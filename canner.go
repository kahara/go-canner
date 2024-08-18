package canner

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	LineSeparator = '\n'
	FileExtention = ".can"
)

type Canner struct {
	Lock   sync.Mutex
	Queue  []Record
	Prefix string
	File   *os.File // One file at a time, assume timestamps of arriving records are in order
	Ticker *time.Ticker
	Term   chan bool
	Ack    chan bool
}

func NewCanner(prefix string) *Canner {
	c := Canner{
		Queue:  make([]Record, 0),
		Prefix: prefix,
		Ticker: time.NewTicker(time.Second),
		Term:   make(chan bool),
		Ack:    make(chan bool),
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

func (c *Canner) Push(record Record) {
	c.Lock.Lock()
	c.Queue = append(c.Queue, record)
	c.Lock.Unlock()
}

func (c *Canner) Flush() {
	if len(c.Queue) == 0 {
		return
	}

	// Prepare to consume incoming records
	outQueue := make([]Record, 0)

	c.Lock.Lock()
	outQueue = append(outQueue, c.Queue...)
	c.Queue = nil
	c.Lock.Unlock()

	// Write records
	for _, record := range outQueue {
		c.Write(record)
	}
}

func (c *Canner) Write(record Record) {
	filename := c.Filename(record)

	if c.File != nil && filename != c.File.Name() {
		if err := c.File.Close(); err != nil {
			panic(err)
		}
		c.File = nil
	}

	if c.File == nil {
		if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
			panic(err)
		}

		if file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644); err != nil {
			panic(err)
		} else {
			c.File = file
		}
	}

	if buf, err := record.Encode(); err != nil {
		panic(err)
	} else {
		buf = append(buf, LineSeparator)
		if _, err := c.File.Write(buf); err != nil {
			panic(err)
		}
	}
}

func (c *Canner) Filename(record Record) string {
	return filepath.Join(c.Prefix,
		record.Timestamp.UTC().Truncate(24*time.Hour).Format(time.RFC3339),
		fmt.Sprintf("%s%s", record.Timestamp.UTC().Truncate(time.Hour).Format(time.RFC3339), FileExtention))
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
