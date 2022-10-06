package log

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type writerStackdriver struct {
	ProjectID string
}

func NewWriterStackdriver(projectID string) Writer {
	return &writerStackdriver{projectID}
}

func (w *writerStackdriver) Request(
	severity Severity,
	traceID string,
	applicationLogs []*EntryChild,
	r *http.Request,
	status int,
	at time.Time,
	dr time.Duration,
) {
	u := *r.URL
	u.Fragment = ""

	uri, err := url.QueryUnescape(u.RequestURI())
	if err != nil {
		uri = u.RequestURI()
	}

	falseV := false

	e := &Entry{
		Severity: severity.String(),
		Time:     Time(at),
		Trace:    fmt.Sprintf("projects/%s/traces/%s", w.ProjectID, traceID),
		TraceID:  traceID,
		Childs:   applicationLogs,
		Message:  "",
		HTTPRequest: &EntryHTTPRequest{
			RequestMethod:                  r.Method,
			RequestURL:                     uri,
			RequestSize:                    r.ContentLength,
			Status:                         status,
			UserAgent:                      r.UserAgent(),
			Referer:                        r.Referer(),
			Latency:                        Duration(dr),
			CacheLookup:                    &falseV,
			CacheHit:                       &falseV,
			CacheValidatedWithOriginServer: &falseV,
			CacheFillBytes:                 nil,
			Protocol:                       r.Proto,
		},
	}
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stderr, string(b)+"\n")
}

func (w *writerStackdriver) Job(
	severity Severity,
	traceID string,
	applicationLogs []*EntryChild,
) {
	e := &Entry{
		Severity:    severity.String(),
		Time:        Time(time.Now()),
		Trace:       fmt.Sprintf("projects/%s/traces/%s", w.ProjectID, traceID),
		TraceID:     traceID,
		Childs:      applicationLogs,
		Message:     "",
		HTTPRequest: &EntryHTTPRequest{},
	}
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stderr, string(b)+"\n")
}

func (w *writerStackdriver) Application(
	severity Severity,
	traceID string,
	msg string,
	file string,
	line int64,
	function string,
	at time.Time,
) {
	e := &Entry{
		Severity: severity.String(),
		Time:     Time(at),
		Trace:    fmt.Sprintf("projects/%s/traces/%s", w.ProjectID, traceID),
		Message:  fmt.Sprintf("%s:%d [%s] %s", file, line, function, msg),
	}
	b, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stdout, string(b)+"\n")
}
