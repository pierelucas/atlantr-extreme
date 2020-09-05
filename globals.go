package main

import "github.com/pierelucas/atlantr-extreme/data"

const (
	upload  = false
	backend = "localhost:56650"

	appID = "7881B883764f54B5"

	configpath          = "conf.json"
	defaultVALIDFILE    = "valid"
	defaultNOTFOUNDFILE = "notfound"
	defaultHOSTERFILE   = "hosters.txt"
	defaultMAXJOBS      = 100
	defaultBUFFERSIZE   = 1
)

var (
	hosterData  map[string]*data.Host
	socksData   []string
	matcherData []string
	lastline    uint64
	conf        *data.Config
	machineID   string
)
