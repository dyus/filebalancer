package balancer

import (
	"slices"
	"sync"
	"sync/atomic"
)

type roundRobin struct {
	hosts   []string
	current atomic.Uint32
	mutex   sync.Mutex
}

func (r *roundRobin) GetHosts(count int) []string {
	index := r.current.Add(uint32(count))
	hosts := make([]string, 0, count)

	for i := range uint32(count) {
		hosts = append(hosts, r.hosts[(index+i)%uint32(len(r.hosts))])
	}

	return hosts
}

func (r *roundRobin) AddHost(host string) {
	r.mutex.Lock()

	defer r.mutex.Unlock()

	if slices.Contains(r.hosts, host) {
		return
	}

	r.hosts = append(r.hosts, host)
}

func NewRoundRobinBalancer(hosts []string) Balancer {
	return &roundRobin{hosts: hosts}
}
