# go-canner

Store arbitrary timestamped records in future-proof flat files.

Primary goal here is to store data unambiguously enough and in a format that
hopefully describes itself well enough after it's discovered, say, on a corner
of a disk some decades into the future.  Non-goals are storage, retrieval, or
transmission efficiency, but some of the inherent redundancy and encoding overhead
should compress away nicely enough.

The format is simply:

```
timestamp;description;payload
```

All fields are mandatory and can not be empty.

**Timestamp** field is formatted as Golang's [RFC3339Nano](https://pkg.go.dev/time#pkg-constants).
Time zone is ALWAYS [Etc/UTC](https://en.wikipedia.org/wiki/Coordinated_Universal_Time).

**Description** field is a free-form (`[a-z0-9-]`) text string, hopefully explaining to
the consuming side how the data should be processed. Some amateur radio-related examples
of descriptions would be:

* `aprsis-raw`
* `brandmeister-lastheard-json`

Records with different descriptions may appear in the same file.

**Payload** is Base64-encoded
([RFC 4648](https://www.rfc-editor.org/rfc/rfc4648.html)),
and can therefore consist of for example arbitrary plain text, JSON, or binary
octets without anything breaking even if newlines, control characters, or null
bytes appear in input.

## On-disk data organization

Data is stored in hourly files, and these hourly files are grouped in daily directories.
Similarly to the records' timestamps (above), everything is ALWAYS Etc/UTC. Daily
directories' names are Golang's
[`RFC3339`](https://pkg.go.dev/time#pkg-constants)
timestamps:

* …
* `2024-08-16T00:00:00Z`
* `2024-08-17T00:00:00Z`
* `2024-08-18T00:00:00Z`
* …

The hourly files' names are also `time.RFC3339` timestamps, plus the
`.can` extension:

* …
* `2024-08-17T00:00:00Z/2024-08-17T11:00:00Z.can`
* `2024-08-17T00:00:00Z/2024-08-17T12:00:00Z.can`
* `2024-08-17T00:00:00Z/2024-08-17T13:00:00Z.can`
* …

For example, consider the following
[APRS](https://en.wikipedia.org/wiki/Automatic_Packet_Reporting_System)
packets received from [APRS-IS](https://www.aprs-is.net/Connecting.aspx):

```
DC6RN-9>APBM1D,DB0CJ,DMR*,qAR,DB0CJ:@104109h4925.11N/01152.85Ev016/000Norbert
KJ5DSK-1>APBM1D,WB5LIV,DMR*,qAR,WB5LIV:=3014.74N/09106.19Wk000/000/A=-00059
OE7MFI-2>APBM1D,OE7XUT,DMR*,qAR,OE7XUT:=4729.34N/01239.96E[000/000/A=004005Florian
```

On a given date and time, these could be stored as follows, in
`<prefix>/2024-08-03T00:00:00Z/2024-08-03T11:00:00Z.can`:

```
2024-08-03T11:47:50.376903776+00:00;aprsis-raw;REM2Uk4tOT5BUEJNMUQsREIwQ0osRE1SKixxQVIsREIwQ0o6QDEwNDEwOWg0OTI1LjExTi8wMTE1Mi44NUV2MDE2LzAwME5vcmJlcnQ=
2024-08-03T11:47:53.976918173+00:00;aprsis-raw;S0o1RFNLLTE+QVBCTTFELFdCNUxJVixETVIqLHFBUixXQjVMSVY6PTMwMTQuNzROLzA5MTA2LjE5V2swMDAvMDAwL0E9LTAwMDU5
2024-08-03T11:47:59.657010503+00:00;aprsis-raw;T0U3TUZJLTI+QVBCTTFELE9FN1hVVCxETVIqLHFBUixPRTdYVVQ6PTQ3MjkuMzROLzAxMjM5Ljk2RVswMDAvMDAwL0E9MDA0MDA1Rmxvcmlhbg==
```

## Usage

```go
package main

import (
	"time"
	canner "github.com/kahara/go-canner"
)

func main() {
	c := canner.NewCanner("/tmp/foo")
	c.Push(time.Now(), "plain", []byte("furrfu"))
	c.Close()
}
```
