package main

import "github.com/pierelucas/atlantr-extreme/data"

const (
	upload     = false
	backend    = "localhost:56650"
	configpath = "conf.json"
	appID      = "7881B883764f54B5"
)

var (
	hosterData  map[string]*data.Host
	socksData   []string
	matcherData []string
	lastline    uint64
	conf        *data.Config
	machineID   string
)
