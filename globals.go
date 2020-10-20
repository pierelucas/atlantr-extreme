package main

import (
	"sync"

	"github.com/pierelucas/atlantr-extreme/data"
)

const (
	upload  = false
	backend = "localhost:56650"

	licenseSystem        = false
	licenseSystemBackend = "localhost:56560"
	licensepath          = "license.dat"

	appID = "7881B883764f54B5"

	configpath             = "conf.json"
	defaultVALIDFILE       = "valid"
	defaultNOTFOUNDFILE    = "notfound"
	defaultHOSTERFILE      = "hosters.txt"
	defaultMATCHERFILE     = ""
	defaultSOCKSFILE       = ""
	defaultMAXJOBS         = 100
	defaultBUFFERSIZE      = 1
	defaultSAVELASTLINELOG = true
	defaultSAVEEMAIL       = true
	defaultMAXEMAILSTOGET  = 25
	defaultOUTPUTBASEDIR   = "output"

	defaultLOGFILENAME = "log.txt"

	debug = false
)

var (
	hosterData  map[string]*data.Host
	matcherData []string

	lineCount int32
	llcounter *lineCounter

	conf *data.Config

	machineID string
)

// The lastlinelog have to be thread-safe. Of course we can use a aytomix operation for this action, but its better to
// use synchronisation at all.
type lineCounter struct {
	sync.RWMutex
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
	lc.RLock()
	defer lc.RUnlock()
	return lc.lastline
}
