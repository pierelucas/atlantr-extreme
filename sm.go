package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/pierelucas/atlantr-extreme/proxy"
)

// Job --
type Job struct {
	lCount int
	User   string
	Pass   string
}

type validProxies struct {
	sync.RWMutex
	proxies    <-chan proxy.Proxy
	validSocks []string
	len        int
}

func (v *validProxies) GetRandomSocks() string {
	rand.Seed(time.Now().UnixNano())
	v.RLock()
	defer v.RUnlock()
	rsocks := v.validSocks[rand.Intn(v.len-1)]
	return rsocks
}

type sm struct {
	jobCH        chan *Job
	resultCH     chan *Job
	uploadCH     chan *Job
	notFoundCH   chan *Job
	validProxies *validProxies
}
