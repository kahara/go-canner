package canner

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const FileExtention = ".can"

type Canner struct {
	InLock   sync.Mutex
	InQueue  []Record
	OutQueue []Record
	Prefix   string
	File     *os.File // One file at a time, assume timestamps of arriving records are in order
	Ticker   *time.Ticker
	Term     chan bool
	Ack      chan bool
}

func NewCanner(prefix string) *Canner {
	c := Canner{
		InQueue:  make([]Record, 0),
		OutQueue: make([]Record, 0),
		Prefix:   prefix,
		Ticker:   time.NewTicker(time.Second),
		Term:     make(chan bool),
		Ack:      make(chan bool),
	}

	// Flush periodically
	go func() {
		for {
			select {
			case <-c.Term:
				c.Flush()
				c.Ack <- true // Notify Close() we're done flushing
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
	if len(c.InQueue) == 0 {
		return
	}

	// Prepare to consume incoming records
	c.InLock.Lock()
	c.OutQueue = append(c.OutQueue, c.InQueue...)
	c.InQueue = nil
	c.InLock.Unlock()

	// Write records
	for _, record := range c.OutQueue {
		c.Write(record)
	}
	c.OutQueue = nil
}

func (c *Canner) Write(r Record) {
	filename := c.Filename(r)

	if c.File != nil && filename != c.File.Name() {
		if err := c.File.Close(); err != nil {
			panic(err)
		}
	}

	if file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644); err != nil {
		panic(err)
	} else {
		c.File = file
	}

	if buf, err := r.Encode(); err != nil {
		panic(err)
	} else {
		buf = append(buf, '\n')
		if _, err := c.File.Write(buf); err != nil {
			panic(err)
		}
	}
}

func (c *Canner) Filename(r Record) string {
	timestamp := r.Timestamp.Truncate(time.Hour)

	return filepath.Join(c.Prefix, fmt.Sprintf("%s%s", timestamp.UTC().Format(time.RFC3339), FileExtention))
}

func (c *Canner) Close() {
	c.Ticker.Stop()
	c.Term <- true
	<-c.Ack

	if c.File == nil {
		return
	}
	if err := c.File.Close(); err != nil {
		panic(err)
	}
}
