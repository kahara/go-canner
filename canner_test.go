package canner

import (
	"os"
	"sync"
	"testing"
	"time"
)

//func TestCannerRoundtrip(t *testing.T) {
//	type args struct {
//		files map[string][]Record
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "",
//			args: args{
//				files: map[string][]Record{
//					"foo": []Record{},
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewCanner(tt.args.prefix, tt.args.suffix); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewCanner() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestCanner_Filename(t *testing.T) {
	type fields struct {
		InLock   sync.Mutex
		InQueue  []Record
		OutQueue []Record
		Prefix   string
		Suffix   string
		File     os.File
		Ticker   *time.Ticker
		Term     chan bool
		Ack      chan bool
	}
	type args struct {
		r Record
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Canner{
				InLock:   tt.fields.InLock,
				InQueue:  tt.fields.InQueue,
				OutQueue: tt.fields.OutQueue,
				Prefix:   tt.fields.Prefix,
				Suffix:   tt.fields.Suffix,
				File:     tt.fields.File,
				Ticker:   tt.fields.Ticker,
				Term:     tt.fields.Term,
				Ack:      tt.fields.Ack,
			}
			if got := c.Filename(tt.args.r); got != tt.want {
				t.Errorf("Filename() = %v, want %v", got, tt.want)
			}
		})
	}
}
