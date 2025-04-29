package utils

import (
	"sync"
	"syscall"
	"time"
)

func NewProcessTracker() *ProcessTracker {
	return &ProcessTracker{
		pids:  map[int]bool{},
		mutex: &sync.Mutex{},
	}
}

type ProcessTracker struct {
	pids  map[int]bool
	mutex *sync.Mutex
}

func (p *ProcessTracker) Add(pid int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.pids[pid] = true
}

func (p *ProcessTracker) kill(pid int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	_ = syscall.Kill(-pid, syscall.SIGKILL)
	delete(p.pids, pid)
}

func (p *ProcessTracker) KillAll() {
	for pid := range p.pids {
		p.kill(pid)
	}

	time.Sleep(time.Millisecond * 250)
}
