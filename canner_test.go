package canner

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
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
		groups map[string][]Record
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Inside one hour",
			args: args{
				groups: map[string][]Record{
					"2024-08-03T00:00:00Z/2024-08-03T11:00:00Z": ParseRecords([]string{
						"2024-08-03T11:00:00+00:00;aprsis-raw;S0MxUkFZLTEwPkFQQk0xRCxLQzFWQVgsRE1SKixxQVIsS0MxVkFYOkAxNTM5MzdoNDE1My4zM04vMDcxMTIuMTJXdjIxNC8wMDBSYXkgQVQtNTc4VVYgTW9iaWxl",
						"2024-08-03T11:13:50.376903776+00:00;aprsis-raw;REM2Uk4tOT5BUEJNMUQsREIwQ0osRE1SKixxQVIsREIwQ0o6QDEwNDEwOWg0OTI1LjExTi8wMTE1Mi44NUV2MDE2LzAwME5vcmJlcnQ=",
						"2024-08-03T11:17:53.976918173+00:00;aprsis-raw;S0o1RFNLLTE+QVBCTTFELFdCNUxJVixETVIqLHFBUixXQjVMSVY6PTMwMTQuNzROLzA5MTA2LjE5V2swMDAvMDAwL0E9LTAwMDU5",
						"2024-08-03T11:23:59.657010503+00:00;aprsis-raw;T0U3TUZJLTI+QVBCTTFELE9FN1hVVCxETVIqLHFBUixPRTdYVVQ6PTQ3MjkuMzROLzAxMjM5Ljk2RVswMDAvMDAwL0E9MDA0MDA1Rmxvcmlhbg==",
						"2024-08-03T11:59:59.999999999+00:00;aprsis-raw;TTdPREEtNz5BUEJNMUQsTTdPREEsRE1SKixxQVIsTTdPREE6QDE1MzkzMmg1NDE1LjUxTi8wMDEyNS40NldbMjA0LzAwME03T0RBIFRlc3Rpbmc=",
					}),
				},
			},
		},
		{
			name: "Span two hours",
			args: args{
				groups: map[string][]Record{
					"2038-01-19T00:00:00Z/2038-01-19T03:00:00Z": ParseRecords([]string{
						"2038-01-19T03:14:07Z;plain;Rm9vIQ==",
						"2038-01-19T03:14:08Z;plain;QmFy",
						"2038-01-19T03:59:59.1Z;plain;YmF6",
					}),
					"2038-01-19T00:00:00Z/2038-01-19T04:00:00Z": ParseRecords([]string{
						"2038-01-19T04:00:00.123456789Z;plain;cXV1eA==",
						"2038-01-19T04:00:00.3000Z;plain;ZnVycmZ1",
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
			defer os.RemoveAll(tempDir)

			canner := NewCanner(tempDir)
			for _, records := range tt.args.groups {
				for _, record := range records {
					canner.Push(record)
				}
			}
			canner.Close()

			for truncatedTimestamp, records := range tt.args.groups {
				filename := filepath.Join(canner.Prefix, fmt.Sprintf("%s%s", truncatedTimestamp, FileExtention))
				t.Logf("Filename %s", filename)

				file, err := os.Open(filename)
				if err != nil {
					panic(err)
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				recordsRemaining := len(records)
				for {
					if scanner.Scan() {
						line := scanner.Bytes()
						t.Logf("line %s", line)

						record, err := NewRecord(string(line))
						if err != nil {
							panic(err)
						}
						t.Logf("record %#v", record)

						// FIXME streamline this
						for _, r := range records {
							a, _ := r.Encode()
							b, _ := record.Encode()
							if bytes.Equal(a, b) {
								recordsRemaining--
								break
							}
						}
					} else {
						break
					}
				}

				if recordsRemaining != 0 {
					t.Errorf("There are %d unaccounted for records %#v", len(records), records)
				}

			}

		})
	}
}

func TestCanner_Filename(t *testing.T) {
	type fields struct {
		Prefix string
	}
	type args struct {
		record Record
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "All zeros",
			fields: fields{
				Prefix: "/path/to",
			},
			args: args{
				record: Record{
					Timestamp: time.Time{},
				},
			},
			want: "/path/to/0001-01-01T00:00:00Z/0001-01-01T00:00:00Z.can",
		},
		{
			name: "All zeros and relative path",
			fields: fields{
				Prefix: "../path/to",
			},
			args: args{
				record: Record{
					Timestamp: time.Time{},
				},
			},
			want: "../path/to/0001-01-01T00:00:00Z/0001-01-01T00:00:00Z.can",
		},
		{
			name: "A timestamp",
			fields: fields{
				Prefix: "/path/to",
			},
			args: args{
				record: Record{
					Timestamp: func() time.Time {
						t, _ := time.Parse(time.RFC3339Nano, "2024-08-17T14:31:42.503657010Z")
						return t
					}(),
				},
			},
			want: "/path/to/2024-08-17T00:00:00Z/2024-08-17T14:00:00Z.can",
		},
		{
			name: "Past y2k38",
			fields: fields{
				Prefix: "/path/to",
			},
			args: args{
				record: Record{
					// 2038-01-19T03:15:11.700Z
					Timestamp: func() time.Time {
						t, _ := time.Parse(time.RFC3339Nano, "2038-01-19T03:15:11.700Z")
						return t
					}(),
				},
			},
			want: "/path/to/2038-01-19T00:00:00Z/2038-01-19T03:00:00Z.can",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Canner{
				Prefix: tt.fields.Prefix,
			}
			if got := c.Filename(tt.args.record); got != tt.want {
				t.Errorf("Filename() = %v, want %v", got, tt.want)
			}
		})
	}
}
