package main

import "github.com/pierelucas/atlantr-extreme/proxy"

// Job --
type Job struct {
	lCount int
	User   string
	Pass   string
}

type validProxies struct {
	proxies    <-chan proxy.Proxy
	validSocks []string
	len        int
}

type sm struct {
	jobCH        chan *Job
	resultCH     chan *Job
	notFoundCH   chan *Job
	validProxies *validProxies
}
