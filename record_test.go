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
		{
			name: "APRS-IS packet with Zulu time",
			args: args{
				encoded: "2024-08-03T11:47:50.376903776Z;aprsis-raw;REM2Uk4tOT5BUEJNMUQsREIwQ0osRE1SKixxQVIsREIwQ0o6QDEwNDEwOWg0OTI1LjExTi8wMTE1Mi44NUV2MDE2LzAwME5vcmJlcnQ=",
			},
			want: &Record{
				Timestamp: func() time.Time {
					t, _ := time.Parse(time.RFC3339Nano, "2024-08-03T11:47:50.376903776Z")
					return t
				}(),
				Description: "aprsis-raw",
				Payload:     []byte("DC6RN-9>APBM1D,DB0CJ,DMR*,qAR,DB0CJ:@104109h4925.11N/01152.85Ev016/000Norbert"),
			},
			wantErr: false,
		},
		{
			name: "APRS-IS packet with +0 time",
			args: args{
				encoded: "2024-08-03T11:47:59.657010503+00:00;aprsis-raw;T0U3TUZJLTI+QVBCTTFELE9FN1hVVCxETVIqLHFBUixPRTdYVVQ6PTQ3MjkuMzROLzAxMjM5Ljk2RVswMDAvMDAwL0E9MDA0MDA1Rmxvcmlhbg==",
			},
			want: &Record{
				Timestamp: func() time.Time {
					t, _ := time.Parse(time.RFC3339Nano, "2024-08-03T11:47:59.657010503+00:00")
					return t
				}(),
				Description: "aprsis-raw",
				Payload:     []byte("OE7MFI-2>APBM1D,OE7XUT,DMR*,qAR,OE7XUT:=4729.34N/01239.96E[000/000/A=004005Florian"),
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
