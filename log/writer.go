package log

import (
	"net/http"
	"time"
)

// Writer ... ログの出力
type Writer interface {
	Request(
		severity Severity,
		traceID string,
		applicationLogs []*EntryChild,
		r *http.Request,
		status int,
		at time.Time,
		dr time.Duration,
	)

	Job(
		severity Severity,
		traceID string,
		applicationLogs []*EntryChild,
	)

	Application(
		severity Severity,
		traceID string,
		msg string,
		file string,
		line int64,
		function string,
		at time.Time,
	)
}
