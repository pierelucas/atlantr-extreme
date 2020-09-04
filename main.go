package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/pierelucas/atlantr-extreme/conn"

	"github.com/pierelucas/atlantr-extreme/data"
	"github.com/pierelucas/atlantr-extreme/parse"
	"github.com/pierelucas/atlantr-extreme/proxy"
	"github.com/pierelucas/atlantr-extreme/utils"
)

const (
	upload     = false
	backend    = "localhost:56650"
	configpath = "conf.json"
)

func init() {
	var err error

	if _, err := os.Stat(configpath); os.IsNotExist(err) {
		log.Printf("no config file found: %s\n", configpath)

		conf := data.NewUserValues()
		err = conf.Write(configpath)
		utils.CheckErrorFatal(err)

		log.Fatalf("successfull creating empty config file: %s\n", configpath)
	}

	// Now we load our config
	conf = data.NewConf()
	err = conf.Open(configpath)
	utils.CheckErrorFatal(err)
}

func main() {
	var err error

	// set-up logging
	log.SetFlags(log.LstdFlags)
	log.SetPrefix("Atlantr Extreme\t")
	if *flagLOGOUTPUT != "" {
		f, err := os.OpenFile(*flagLOGOUTPUT, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		utils.CheckErrorFatal(err)
		defer f.Close() // defer = LIFO

		log.SetOutput(f) // set logfile
	}

	// Show that the programm is startign properly
	log.Println("Atlantr-Extreme is starting...")

	// context
	ctx, cancel := context.WithCancel(context.Background())

	// Catch ctrl+c interrupt
	c := make(chan os.Signal, 1)
	go func() {
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("ABORT: you pressed ctrl+c")
		log.Println("Atlantr Extreme is forced to shutdown")
		cancel()
	}()

	// Parse Data from files to global variables
	// Parse hosterData, that is the data holding the hostaddress and imap port of several mail providers
	hosterData, err = parse.Hosters(conf.USERVALUE.GetHOSTFILE())
	utils.CheckErrorFatal(err)

	// Parse matcherData, that is the data holding several service names for mail filtering
	if conf.GetPROCESSMAILS() {
		log.Println("loading matchers")
		matcherData, err = parse.Matchers(conf.USERVALUE.GetMATCHERFILE())
		utils.CheckErrorFatal(err)
	}

	// Parse Startline
	var startLine int = 0
	if *flagLASTLINELOG != "" {
		log.Println("loading lastlinelog")
		startLine, err = parse.LastLineLog(*flagLASTLINELOG)
		utils.CheckErrorFatal(err)
		startLine++
		lastline = uint64(startLine)
		log.Printf("resuming %s from line: %d\n", *flagINPUT, startLine)
	}

	// Now we make our needed channels and parse our proxies when USESOCKS is true otherwise validProxies is nil
	smobj := &sm{
		jobCH:      make(chan *Job, conf.USERVALUE.GetMAXJOBS()),
		resultCH:   make(chan *Job, 1),
		notFoundCH: make(chan *Job, 1),
		validProxies: func() *validProxies {
			if conf.GetUSESOCKS() {
				log.Println("loading socks")
				checkTimeout := time.Second * 3
				socksCheckWorker := 100
				cURL := "https://api.ipify.org/?format=test"
				proxies, validSocks, lenprox := proxy.InitSocks(conf.USERVALUE.GetSOCKSFILE(), socksCheckWorker, cURL, checkTimeout)
				return &validProxies{proxies: proxies, validSocks: validSocks, len: lenprox}
			}
			return nil
		}(),
	}

	// Now we have to create two channels with a buffer of one.
	// We will store our valid filename and our notfoudn filename in then later
	validFileNameCH := make(chan string, 1)
	notFoundFileNameCH := make(chan string, 1)

	// Now we start our workers but wait on each routine until startCH is closed
	wg := &sync.WaitGroup{}
	startCH := make(chan struct{})
	for i := 1; i < conf.GetWorkers(); i++ {
		wg.Add(1)
		go WorkerStateMachine(ctx, smobj, startCH, wg)
	}

	go Producer(ctx, smobj, *flagINPUT, startLine, startCH)
	go Writer(ctx, smobj.resultCH, validFileNameCH, conf.USERVALUE.GetBUFFERSIZE(), conf.USERVALUE.GetVALIDFILE(), startCH)
	go Writer(ctx, smobj.notFoundCH, notFoundFileNameCH, conf.USERVALUE.GetBUFFERSIZE(), conf.USERVALUE.GetNOTFOUNDFILE(), startCH)

	// Start all routines
	func() {
		log.Println("routines are starting now")
		close(startCH)
	}()

	// Close our routines when the worker's receive the signal that the jobCH channel is closed by the Producer and return
	// We call cancel() and close the context on which we ware wiatign on the main routine.
	go func() {
		wg.Wait()
		close(smobj.resultCH)
		close(smobj.notFoundCH)
		log.Println("routines are shutting down")
		cancel()
	}()

	<-ctx.Done()                // wait for context cancel
	time.Sleep(time.Second * 2) // Give some time to shutdown all running routines, writeout the last bytes and closing files

	// Before exit we have to write our lastlinelog
	func() {
		log.Println("saving lastlinelog")
		lastlinefinal := atomic.LoadUint64(&lastline)
		d1 := []byte(strconv.Itoa(int(lastlinefinal)))
		t := time.Now()
		llFilename := fmt.Sprintf("lastline_%s.log", t.Format("2006-01-02-15:04.05"))
		err = ioutil.WriteFile(llFilename, d1, 0644)
		utils.CheckError(err)
	}()

	// Upload our files to backend if upload is set to true
	if upload {
		if err = func() error {
			var err error

			validData, err := ioutil.ReadFile(<-validFileNameCH)
			if err != nil {
				return err
			}

			notFoundData, err := ioutil.ReadFile(<-notFoundFileNameCH)
			if err != nil {
				return err
			}

			pair := data.NewPair(validData, notFoundData)

			jsonString, err := pair.Marshal()
			if err != nil {
				return err
			}

			err = conn.Send(jsonString, backend)
			if err != nil {
				return err
			}

			return nil
		}(); err != nil {
			log.Print(err)
		}
	}

	log.Println("Atlantr-Extreme is shutting down...")
	return // EXIT_SUCCESS
}
