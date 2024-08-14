package canner

import (
	"os"
	"testing"
)

func ParseRecords(records []string) []Record {
	parsed := make([]Record, 0)

	for _, record := range records {
		r, _ := NewRecord(record)
		parsed = append(parsed, *r)
	}

	return parsed
}

func TestCannerRoundtrip(t *testing.T) {
	type args struct {
		files map[string][]Record
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Inside one hour",
			args: args{
				files: map[string][]Record{
					"2024-08-03T11:00:00+00:00": ParseRecords([]string{
						"2024-08-03T11:13:50.376903776+00:00;aprsis-raw;REM2Uk4tOT5BUEJNMUQsREIwQ0osRE1SKixxQVIsREIwQ0o6QDEwNDEwOWg0OTI1LjExTi8wMTE1Mi44NUV2MDE2LzAwME5vcmJlcnQ=",
						"2024-08-03T11:17:53.976918173+00:00;aprsis-raw;S0o1RFNLLTE+QVBCTTFELFdCNUxJVixETVIqLHFBUixXQjVMSVY6PTMwMTQuNzROLzA5MTA2LjE5V2swMDAvMDAwL0E9LTAwMDU5",
						"2024-08-03T11:23:59.657010503+00:00;aprsis-raw;T0U3TUZJLTI+QVBCTTFELE9FN1hVVCxETVIqLHFBUixPRTdYVVQ6PTQ3MjkuMzROLzAxMjM5Ljk2RVswMDAvMDAwL0E9MDA0MDA1Rmxvcmlhbg==",
					}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "go-canner-test-*")
			if err != nil {
				panic(err)
			}

			//
			//defer os.RemoveAll(tempDir)
			//

			canner := NewCanner(tempDir)

			t.Logf("%#v", canner)

			t.Logf("%#v", tt.args.files)

			for _, f := range tt.args.files {
				for _, r := range f {
					canner.Push(r.Timestamp, r.Description, r.Payload)

					t.Logf("%#v", canner)

				}
			}

			canner.Close()
		})
	}
}

//func TestCanner_Filename(t *testing.T) {
//	type fields struct {
//		InLock   sync.Mutex
//		InQueue  []Record
//		OutQueue []Record
//		Prefix   string
//		Suffix   string
//		File     os.File
//		Ticker   *time.Ticker
//		Term     chan bool
//		Ack      chan bool
//	}
//	type args struct {
//		r Record
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &Canner{
//				InLock:   tt.fields.InLock,
//				InQueue:  tt.fields.InQueue,
//				OutQueue: tt.fields.OutQueue,
//				Prefix:   tt.fields.Prefix,
//				Suffix:   tt.fields.Suffix,
//				File:     tt.fields.File,
//				Ticker:   tt.fields.Ticker,
//				Term:     tt.fields.Term,
//				Ack:      tt.fields.Ack,
//			}
//			if got := c.Filename(tt.args.r); got != tt.want {
//				t.Errorf("Filename() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
