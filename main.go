/*
	Atlantr-Extreme Mailgrabber
	AUTHOR: github.com/pierelucas
*/

package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	pbar "github.com/schollz/progressbar/v3"

	"github.com/pierelucas/atlantr-extreme/conn"
	"github.com/pierelucas/atlantr-extreme/license"

	"github.com/pierelucas/atlantr-extreme/data"
	"github.com/pierelucas/atlantr-extreme/parse"
	"github.com/pierelucas/atlantr-extreme/proxy"
	"github.com/pierelucas/atlantr-extreme/utils"
)

// let's fool around and use init() for the license and config setup, instead of variable initialisation (:
// i know, maybe we define a own function for this later and use init for the real init things
func init() {
	var err error

	// Generate unique computer identifier
	machineID, err = utils.GenerateID(appID)
	utils.CheckErrorPrintFatal(err)

	// check for a valid license if licenseSystem is set to true
	// We do this before any user input or any input validation can happen
	if licenseSystem {
		// Check if there is a license.dat file in our working folder
		var licenseKey string
		if _, err := os.Stat(licensepath); os.IsNotExist(err) {
			fmt.Printf("no license file found: %s\n", licensepath)

			// read license key from commandline
			read := bufio.NewReader(os.Stdin)

			// reade license key from commandline
			fmt.Printf("Please enter license key\n\n-> ")
			licenseKey, _ = read.ReadString('\n')
			licenseKey = strings.Replace(licenseKey, "\n", "", -1)
		} else {
			// read license
			data, err := ioutil.ReadFile(licensepath)
			utils.CheckErrorPrintFatal(err)

			licenseKey = strings.Replace(strings.TrimSpace(string(data)), "\n", "", -1)
		}

		// validate user input and check for some possible exploit cases
		err = license.ValidateOrKill(licenseKey)
		utils.CheckErrorPrintFatal(err)

		// Now we make a json string with our machineID and license key
		pair, err := license.NewPair(licenseKey, machineID)
		utils.CheckErrorPrintFatal(err)

		jsonString, err := pair.Marshal()
		utils.CheckErrorPrintFatal(err)

		// Now we send our key to the backend server, if err != nil the key is not valid, already used or expired
		err = conn.Send(jsonString, licenseSystemBackend, debug)
		if err != nil {
			fmt.Println("error: your license is not valid, already in use or expired. Please contact your vendor for support\nAlso please make sure you have a working internet connection, when not, fix that and try again")
			os.Exit(1)
		}

		// write license
		ioutil.WriteFile(licensepath, []byte(licenseKey), 0644)
	}

	// check if the config file exist
	if _, err := os.Stat(configpath); os.IsNotExist(err) {
		fmt.Printf("no config file found: %s\n", configpath)

		conf := data.NewUserValues()

		// Set default values
		conf.SetVALIDFILE(defaultVALIDFILE)
		conf.SetNOTFOUNDFILE(defaultNOTFOUNDFILE)
		conf.SetHOSTFILE(defaultHOSTERFILE)
		conf.SetMAXJOBS(defaultMAXJOBS)
		conf.SetBUFFERSIZE(defaultBUFFERSIZE)
		conf.SetSAVELASTLINELOG(defaultSAVELASTLINELOG)

		// write config
		err = conf.Write(configpath)
		utils.CheckErrorPrintFatal(err)

		fmt.Printf("successfull creating default config file: %s\n", configpath)
		os.Exit(1)
	}

	// Now we load our config
	conf = data.NewConf()
	err = conf.Open(configpath)
	utils.CheckErrorPrintFatal(err)

	// parse and validate commandline args
	parseFlags()

	// Got lineCount from mail:pass input file
	lineCount, err = utils.GotLineCount(*flagINPUT)
	utils.CheckErrorPrintFatal(err)
}

func saveLastLineLog() error {
	utils.MultiLogf("saving lastlinelog\n")

	lastlinefinal := atomic.LoadUint64(&lastline)
	d1 := []byte(strconv.Itoa(int(lastlinefinal)))

	t := time.Now()
	llFilename := fmt.Sprintf("lastline_%s.log", t.Format("2006-01-02-15:04.05"))

	err := ioutil.WriteFile(llFilename, d1, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var err error

	// set-up logging
	f := func() *os.File {
		log.SetFlags(log.LstdFlags)
		log.SetPrefix("Atlantr Extreme\t")

		var logOut string
		switch *flagLOGOUTPUT {
		case "":
			logOut = defaultLOGFILENAME
		default:
			logOut = *flagLOGOUTPUT
		}

		f, err := os.OpenFile(logOut, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		utils.CheckErrorFatal(err)

		log.SetOutput(f) // set logfile

		return f
	}()
	defer f.Close() // defer = LIFO

	// Show that the programm is startign properly
	utils.MultiLogf("Atlantr-Extreme is starting...\n")

	// context
	ctx, cancel := context.WithCancel(context.Background())

	// Catch ctrl+c interrupt
	c := make(chan os.Signal, 1)
	go func() {
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		utils.MultiLogf("ABORT: you pressed ctrl+c\n")
		cancel()
	}()

	// Parse Data from files to global variables
	// Parse hosterData, that is the data holding the hostaddress and imap port of several mail providers
	hosterData, err = parse.Hosters(conf.USERVALUE.GetHOSTFILE())
	utils.CheckErrorFatal(err)

	// Parse matcherData, that is the data holding several service names for mail filtering
	if conf.GetPROCESSMAILS() {
		utils.MultiLogf("loading matchers\n")
		matcherData, err = parse.Matchers(conf.USERVALUE.GetMATCHERFILE())
		utils.CheckErrorFatal(err)
	}

	// Parse Startline
	var startLine int = 0
	if *flagLASTLINELOG != "" {
		utils.MultiLogf("loading lastlinelog\n")
		startLine, err = parse.LastLineLog(*flagLASTLINELOG)
		utils.CheckErrorFatal(err)
		startLine++
		lastline = uint64(startLine)
		utils.MultiLogf("resuming %s from line: %d\n", *flagINPUT, startLine)
	}

	// Now we make our needed channels and parse our proxies when USESOCKS is true otherwise validProxies is nil
	smobj := &sm{
		jobCH:      make(chan *Job, conf.USERVALUE.GetMAXJOBS()),
		resultCH:   make(chan *Job, 1),
		uploadCH:   make(chan *Job, 1),
		notFoundCH: make(chan *Job, 1),
		validProxies: func() *validProxies {
			if conf.GetUSESOCKS() {
				utils.MultiLogf("loading socks\n")
				checkTimeout := time.Second * 3
				socksCheckWorker := 100
				cURL := "https://api.ipify.org/?format=test"
				proxies, validSocks, lenprox := proxy.InitSocks(conf.USERVALUE.GetSOCKSFILE(), socksCheckWorker, cURL, checkTimeout)
				return &validProxies{proxies: proxies, validSocks: validSocks, len: lenprox}
			}
			return nil
		}(),
	}

	// make progressbar
	bar := pbar.NewOptions(int(lineCount),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan][reset]Checking Mails..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	// Now we start our workers but wait on each routine until startCH is closed
	handlerWG := &sync.WaitGroup{}
	startCH := make(chan struct{})
	for i := 1; i < conf.GetWorkers(); i++ {
		handlerWG.Add(1)
		go WorkerStateMachine(ctx, smobj, startCH, handlerWG, bar)
	}

	// We need a second waitgroup for our writer and uploader
	workerWG := &sync.WaitGroup{}
	workerWG.Add(3)

	go Producer(ctx, smobj, *flagINPUT, startLine, startCH) // we don't need a waitgroup on Producer() cause WorkerStateMachine() will not return until the producer is done
	go Writer(ctx, smobj.resultCH, conf.USERVALUE.GetBUFFERSIZE(), conf.USERVALUE.GetVALIDFILE(), startCH, workerWG)
	go Writer(ctx, smobj.notFoundCH, conf.USERVALUE.GetBUFFERSIZE(), conf.USERVALUE.GetNOTFOUNDFILE(), startCH, workerWG)
	go Uploader(ctx, smobj, backend, startCH, workerWG)

	// Start all routines
	go func() {
		utils.MultiLogf("routines are starting now\n")
		close(startCH)
	}()

	// Close our routines when the worker's receive the signal that the jobCH channel is closed by the Producer.
	// We call cancel() and close the context on which we ware waiting on the main routine.
	go func() {
		handlerWG.Wait()
		close(smobj.resultCH) // close the writer and upload channels and let the Writer() and Uploader() routines begin to shutdown
		close(smobj.uploadCH)

		utils.MultiLogf("routines finish and shutting down now, clean-up is starting and files will be written\n")
		workerWG.Wait() // Wait till the last bytes are written and uploads are finished
		cancel()
	}()

	<-ctx.Done()            // wait for context cancel
	time.Sleep(time.Second) // chill down

	// Before exit we have to write our lastlinelog, if SAVELASTLINELOG is true
	if conf.USERVALUE.IsSAVELASTLINELOG() {
		saveLastLineLog()
	}

	utils.MultiLogf("Atlantr-Extreme is shutting down...\n")
	time.Sleep(time.Second) // chill down

	return // EXIT_SUCCESS
}
