package canner

import (
	"reflect"
	"testing"
	"time"
)

func TestNewRecord(t *testing.T) {
	type args struct {
		encoded string
	}
	tests := []struct {
		name    string
		args    args
		want    *Record
		wantErr bool
	}{
		{
			name: "Empty fields",
			args: args{
				encoded: ";;",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid timestamp",
			args: args{
				encoded: "10000-01-01T00:00:00Z;;",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Empty description",
			args: args{
				encoded: "0001-01-01T00:00:00Z;;",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Empty payload",
			args: args{
				encoded: "0001-01-01T00:00:00Z;foo;",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid payload",
			args: args{
				encoded: "0001-01-01T00:00:00Z;foo;bar",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Minimal valid record",
			args: args{
				encoded: "0001-01-01T00:00:00Z;plain;cGF5bG9hZA==",
			},
			want: &Record{
				Timestamp:   time.Time{},
				Description: "plain",
				Payload:     []byte("payload"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRecord(tt.args.encoded)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRecord() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_Encode(t *testing.T) {
	type fields struct {
		Timestamp   time.Time
		Description string
		Payload     []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Empty fields",
			fields: fields{
				Timestamp:   time.Time{},
				Description: "",
				Payload:     nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Empty description",
			fields: fields{
				Timestamp:   time.Time{},
				Description: "",
				Payload:     []byte("foo"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Empty payload",
			fields: fields{
				Timestamp:   time.Time{},
				Description: "foo",
				Payload:     nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Minimal valid record",
			fields: fields{
				Timestamp:   time.Time{},
				Description: "plain",
				Payload:     []byte("payload"),
			},
			want:    []byte("0001-01-01T00:00:00Z;plain;cGF5bG9hZA=="),
			wantErr: false,
		},
		//{
		//	name: "",
		//	fields: fields{
		//		Timestamp:   time.Time{},
		//		Description: "",
		//		Payload:     nil,
		//	},
		//	want:    nil,
		//	wantErr: false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Record{
				Timestamp:   tt.fields.Timestamp,
				Description: tt.fields.Description,
				Payload:     tt.fields.Payload,
			}
			got, err := r.Encode()
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() got = %s, want %s", got, tt.want)
			}
		})
	}
}
