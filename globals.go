package main

import "github.com/pierelucas/atlantr-extreme/data"

var (
	hosterData  map[string]*data.Host
	socksData   []string
	matcherData []string
	lastline    uint64
	conf        *data.Config
)
