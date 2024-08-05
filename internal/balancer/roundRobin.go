package balancer

import (
	"sync/atomic"
)

type Balancer interface {
	GetHosts(int) []string
}
type roundRobin struct {
	hosts   []string
	current atomic.Uint32
}

func (r *roundRobin) GetHosts(count int) []string {
	index := r.current.Add(uint32(count))
	hosts := make([]string, 0, count)
	for i := uint32(0); i < uint32(count); i++ {
		hosts = append(hosts, r.hosts[(index+i)%uint32(len(r.hosts))])
	}
	return hosts
}

func NewRoundRobinBalancer(hosts []string) Balancer {
	return &roundRobin{hosts: hosts}
}
