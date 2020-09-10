package main

import (
	"sync"

	"github.com/pierelucas/atlantr-extreme/data"
)

const (
	upload  = true
	backend = "localhost:56650"

	licenseSystem        = true
	licenseSystemBackend = "localhost:56560"
	licensepath          = "license.dat"

	appID = "7881B883764f54B5"

	configpath             = "conf.json"
	defaultVALIDFILE       = "valid"
	defaultNOTFOUNDFILE    = "notfound"
	defaultHOSTERFILE      = "hosters.txt"
	defaultMAXJOBS         = 100
	defaultBUFFERSIZE      = 1
	defaultSAVELASTLINELOG = true

	defaultLOGFILENAME = "log.txt"

	debug = false
)

var (
	hosterData  map[string]*data.Host
	socksData   []string
	matcherData []string

	lineCount int32
	llcounter *lineCounter

	conf *data.Config

	machineID string
)

// The lastlinelog have to be thread-safe. Of course we can use a aytomix operation for this action, but its better to
// use synchronisation at all.
type lineCounter struct {
	sync.Mutex
	lastline int32
}

func newLineCounter() *lineCounter {
	return &lineCounter{
		lastline: 0,
	}
}

func (lc *lineCounter) add(n int32) {
	lc.Lock()
	defer lc.Unlock()
	lc.lastline += n
}

func (lc *lineCounter) value() int32 {
	lc.Lock()
	defer lc.Unlock()
	return lc.lastline
}
