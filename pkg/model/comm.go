package model

import (
	"sync"
	"time"
)

const MaxAge = time.Minute * 20

type comm struct {
	accessed time.Time
	wSecret  string
	m        sync.RWMutex
	*Waiters
}

func newComm() comm {
	return comm{
		accessed: time.Now(),
		Waiters:  NewWaiters(),
	}
}

func (c *comm) WAuth(wSecret string) bool {
	c.m.Lock()
	defer c.m.Unlock()

	if c.wSecret == "" {
		c.wSecret = wSecret
	}
	return c.wSecret == wSecret
}

func (c *comm) Expired() bool {
	c.m.RLock()
	defer c.m.RUnlock()
	return time.Since(c.accessed) > MaxAge && !c.Waiters.HasWaiters()
}

func (c *comm) Accessed() {
	c.m.Lock()
	defer c.m.Unlock()
	c.accessed = time.Now()
}
