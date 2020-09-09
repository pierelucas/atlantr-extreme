package main

import "github.com/pierelucas/atlantr-extreme/data"

const (
	upload  = false
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

	debug = true
)

var (
	hosterData  map[string]*data.Host
	socksData   []string
	matcherData []string

	lineCount uint64
	lastline  uint64

	conf *data.Config

	machineID string
)
