package canner

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	RecordSeparator = ";"
	TimestampFormat = time.RFC3339Nano
)

type Record struct {
	Timestamp   time.Time
	Description string
	Payload     []byte
}

func NewRecord(encoded string) (*Record, error) {
	var record Record

	parts := strings.Split(encoded, RecordSeparator)

	if len(parts) != 3 {
		return nil, errors.New(fmt.Sprintf("invalid record '%s'", encoded))
	}

	// Timestamp
	if timestamp, err := time.Parse(TimestampFormat, parts[0]); err != nil {
		return nil, err
	} else {
		record.Timestamp = timestamp.UTC()
	}

	// Description
	if parts[1] == "" {
		return nil, errors.New(fmt.Sprintf("empty description '%s'", parts[1]))
	}
	record.Description = parts[1]

	// Payload
	if parts[2] == "" {
		return nil, errors.New(fmt.Sprintf("empty payload '%s'", parts[2]))
	}
	if payload, err := base64.StdEncoding.DecodeString(parts[2]); err != nil {
		return nil, err
	} else {
		record.Payload = payload
	}

	return &record, nil
}

func (r *Record) Encode() ([]byte, error) {
	var encoded []byte

	if r.Description == "" {
		return nil, errors.New(fmt.Sprintf("empty description '%s'", r.Description))
	}
	if r.Payload == nil {
		return nil, errors.New(fmt.Sprintf("empty payload '%s'", r.Payload))
	}

	encoded = append(encoded, r.Timestamp.UTC().Format(TimestampFormat)...)
	encoded = append(encoded, RecordSeparator...)
	encoded = append(encoded, r.Description...)
	encoded = append(encoded, RecordSeparator...)
	encoded = base64.StdEncoding.AppendEncode(encoded, r.Payload)

	return encoded, nil
}
