package log

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/rabee-inc/go-pkg/errcode"
	"github.com/rabee-inc/go-pkg/timeutil"
)

// Logger ... ロガー
type Logger struct {
	Writer            Writer
	MinOutSeverity    Severity
	MaxOuttedSeverity Severity
	TraceID           string
	ResponseStatus    int
	ApplicationLogs   []*EntryChild
}

// IsLogging ... レベル毎のログ出力許可
func (l *Logger) IsLogging(severity Severity) bool {
	return l.MinOutSeverity <= severity
}

// SetOuttedSeverity ... 出力された最大のレベルを設定
func (l *Logger) SetOuttedSeverity(severity Severity) {
	if l.MaxOuttedSeverity < severity {
		l.MaxOuttedSeverity = severity
	}
}

// AddApplicationLog ... アプリケーションログ履歴を記録する
func (l *Logger) AddApplicationLog(severity Severity, file string, line int64, function string, msg string, at time.Time) {
	src := &EntryChild{
		Severity: severity.String(),
		Message:  fmt.Sprintf("%s:%d [%s] %s", file, line, function, msg),
		Time:     Time(at),
	}
	l.ApplicationLogs = append(l.ApplicationLogs, src)
}

// WriteRequest ... リクエストログを出力する
func (l *Logger) WriteRequest(r *http.Request, at time.Time, dr time.Duration) {
	if l.IsLogging(l.MaxOuttedSeverity) {
		l.Writer.Request(
			l.MaxOuttedSeverity,
			l.TraceID,
			l.ApplicationLogs,
			r,
			l.ResponseStatus,
			at,
			dr)
	}
}

// WriteJob ... ジョブログを出力する
func (l *Logger) WriteJob(ctx context.Context) {
	l.Writer.Job(l.MaxOuttedSeverity, l.TraceID, l.ApplicationLogs)
}

// NewLogger ... Loggerを作成する
func NewLogger(writer Writer, minSeverity Severity, traceID string) *Logger {
	return &Logger{
		Writer:            writer,
		MinOutSeverity:    minSeverity,
		MaxOuttedSeverity: SeverityDebug,
		TraceID:           traceID,
	}
}

// SetResponseStatus ... レスポンスのステータスコードを設定する
func SetResponseStatus(ctx context.Context, status int) {
	logger := GetLogger(ctx)
	if logger != nil {
		logger.ResponseStatus = status
	}
}

// Debug ... Debugログの定形を出力する
func Debug(ctx context.Context, err error) {
	severity := SeverityDebug
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Debugf ... Debugログを出力する
func Debugf(ctx context.Context, format string, args ...any) {
	severity := SeverityDebug
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := fmt.Sprintf(format, args...)
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Debuge ... Debugログを出力してエラーを生成する
func Debuge(ctx context.Context, format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityDebug
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return err
}

// Debugc ... Debugログを出力してコード付きのエラーを生成する
func Debugc(ctx context.Context, code int, format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityDebug
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return errcode.Set(err, code)
}

// Info ... Infoログの定形を出力する
func Info(ctx context.Context, err error) {
	severity := SeverityInfo
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Infof ... Infoログを出力する
func Infof(ctx context.Context, format string, args ...any) {
	severity := SeverityInfo
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := fmt.Sprintf(format, args...)
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Infoe ... Infoログを出力してエラーを生成する
func Infoe(ctx context.Context, format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityInfo
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return err
}

// Infoc ... Infoログを出力してコード付きのエラーを生成する
func Infoc(ctx context.Context, code int, format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityInfo
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return errcode.Set(err, code)
}

// Warning ... Warningログの定形を出力する
func Warning(ctx context.Context, err error) {
	severity := SeverityWarning
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Warningf ... Warningログを出力する
func Warningf(ctx context.Context, format string, args ...any) {
	severity := SeverityWarning
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := fmt.Sprintf(format, args...)
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Warninge ... Warningログを出力してエラーを生成する
func Warninge(ctx context.Context, format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityWarning
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return err
}

// Warningc ... Warningログを出力してコード付きのエラーを生成する
func Warningc(ctx context.Context, code int, format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityWarning
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return errcode.Set(err, code)
}

// Error ... Errorログの定形を出力する
func Error(ctx context.Context, err error) {
	severity := SeverityError
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Errorf ... Errorログを出力する
func Errorf(ctx context.Context, format string, args ...any) {
	severity := SeverityError
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := fmt.Sprintf(format, args...)
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Errore ... Errorログを出力してエラーを生成する
func Errore(ctx context.Context, format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityError
	logger := GetLogger(ctx)
	msg := err.Error()
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return err
}

// Errorc ... Errorログを出力してコード付きのエラーを生成する
func Errorc(ctx context.Context, code int, format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityError
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return errcode.Set(err, code)
}

// Critical ... Criticalログの定形を出力する
func Critical(ctx context.Context, err error) {
	severity := SeverityCritical
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Criticalf ... Criticalログを出力する
func Criticalf(ctx context.Context, format string, args ...any) {
	severity := SeverityCritical
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := fmt.Sprintf(format, args...)
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
}

// Criticale ... Criticalログを出力してエラーを生成する
func Criticale(ctx context.Context, format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityCritical
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return err
}

// Criticalc ... Criticalログを出力してコード付きのエラーを生成する
func Criticalc(ctx context.Context, code int, format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	severity := SeverityCritical
	logger := GetLogger(ctx)
	if logger != nil && logger.IsLogging(severity) {
		now := timeutil.Now()
		file, line, function := getFileLine()
		msg := err.Error()
		logger.Writer.Application(
			severity,
			logger.TraceID,
			msg,
			file,
			line,
			function,
			now)
		logger.SetOuttedSeverity(severity)
		logger.AddApplicationLog(severity, file, line, function, msg, now)
	}
	return errcode.Set(err, code)
}

// Panic ... Panicをハンドリングする
func Panic(ctx context.Context, rcvr any) string {
	traces := []string{}
	for depth := 0; ; depth++ {
		if depth < 2 {
			continue
		}
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			break
		}
		trace := fmt.Sprintf("%02d: %v:%d", depth-1, file, line)
		traces = append(traces, trace)
	}
	msg := fmt.Sprintf("panic!! %v\n%s", rcvr, strings.Join(traces, "\n"))
	Criticalf(ctx, msg)
	return msg
}

func getFileLine() (string, int64, string) {
	if pt, file, line, ok := runtime.Caller(2); ok {
		parts := strings.Split(file, "/")
		length := len(parts)
		file := fmt.Sprintf("%s/%s", parts[length-2], parts[length-1])

		fParts := strings.Split(runtime.FuncForPC(pt).Name(), ".")
		fLength := len(fParts)
		return file, int64(line), fParts[fLength-1]
	}
	return "", 0, ""
}
