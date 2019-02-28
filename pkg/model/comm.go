package model

import (
	"time"
)

const MaxAge = time.Minute * 20

type comm struct {
	accessed time.Time
	wSecret  string
	*Waiters
}

func newComm() comm {
	return comm{
		accessed: time.Now(),
		Waiters:  NewWaiters(),
	}
}

func (c *comm) WAuth(wSecret string) bool {
	c.Lock()
	defer c.Unlock()

	if c.wSecret == "" {
		c.wSecret = wSecret
	}
	return c.wSecret == wSecret
}

func (c *comm) Expired() bool {
	c.RLock()
	defer c.RUnlock()
	return time.Since(c.accessed) > MaxAge && !c.Waiters.HasWaiters()
}

func (c *comm) Accessed() {
	c.Lock()
	c.accessed = time.Now()
	c.Unlock()
}
