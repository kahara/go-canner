package canner

import (
	"os"
	"sync"
	"time"
)

type Canner struct {
	InLock   sync.Mutex
	InQueue  []Record
	OutQueue []Record
	Prefix   string
	Suffix   string
	File     os.File // One file at a time, assume timestamps of arriving records are in order
	Ticker   *time.Ticker
	Done     chan bool
}

func NewCanner(prefix string, suffix string) *Canner {
	c := Canner{
		InQueue:  make([]Record, 0),
		OutQueue: make([]Record, 0),
		Prefix:   prefix,
		Suffix:   suffix,
		Ticker:   time.NewTicker(time.Second),
		Done:     make(chan bool),
	}

	// Flush periodically
	go func() {
		for {
			select {
			case <-c.Done:
				c.Flush()
				return
			case <-c.Ticker.C:
				c.Flush()
			}
		}
	}()

	return &c
}

func (c *Canner) Push(t time.Time, d string, p []byte) {
	c.InLock.Lock()
	c.InQueue = append(c.InQueue, Record{
		Timestamp:   t,
		Description: d,
		Payload:     p,
	})
	c.InLock.Unlock()
}

func (c *Canner) Flush() {
	// Prepare to consume incoming records
	c.InLock.Lock()
	c.OutQueue = append(c.OutQueue, c.InQueue...)
	c.InQueue = nil
	c.InLock.Unlock()

	if len(c.OutQueue) < 1 {
		return
	}

	// Write records
	for _, record := range c.OutQueue {
		c.Write(record)
	}
	c.OutQueue = nil
}

func (c *Canner) Write(r Record) {

}

func (c *Canner) Close() {
	c.Ticker.Stop()
	c.Done <- true
}
