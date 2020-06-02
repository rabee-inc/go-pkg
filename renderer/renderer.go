package renderer

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"

	"github.com/unrolled/render"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/rabee-inc/go-pkg/errcode"
	"github.com/rabee-inc/go-pkg/log"
)

// HandleError ... 一番典型的なエラーハンドリング
func HandleError(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	code, ok := errcode.Get(err)
	if !ok {
		Error(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}
	texts := []string{}
	if msg != "" {
		text := fmt.Sprintf("%d %s", code, msg)
		texts = append(texts, text)
	}
	if err != nil {
		text := err.Error()
		texts = append(texts, text)
	}
	text := strings.Join(texts, "\n")

	switch code {
	case http.StatusBadRequest:
		log.Warningf(ctx, text)
		Error(ctx, w, code, text)
	case http.StatusUnauthorized:
		log.Warningf(ctx, text)
		Error(ctx, w, code, text)
	case http.StatusForbidden:
		log.Warningf(ctx, text)
		Error(ctx, w, code, text)
	case http.StatusNotFound:
		log.Warningf(ctx, text)
		Error(ctx, w, code, text)
	default:
		log.Errorf(ctx, text)
		Error(ctx, w, code, text)
	}
}

// Success ... 成功レスポンスをレンダリングする
func Success(ctx context.Context, w http.ResponseWriter) {
	status := http.StatusOK
	r := render.New()
	r.JSON(w, http.StatusOK, NewResponseOK(http.StatusOK))
	log.SetResponseStatus(ctx, status)
}

// Error ... エラーレスポンスをレンダリングする
func Error(ctx context.Context, w http.ResponseWriter, status int, msg string) {
	r := render.New()
	r.JSON(w, status, NewResponseError(status, msg))
	log.SetResponseStatus(ctx, status)
}

// JSON ... JSONをレンダリングする
func JSON(ctx context.Context, w http.ResponseWriter, status int, v interface{}) {
	r := render.New()
	r.JSON(w, status, v)
	log.SetResponseStatus(ctx, status)
}

// HTML ... HTMLをレンダリングする
func HTML(ctx context.Context, w http.ResponseWriter, status int, name string, values interface{}) {
	r := render.New()
	r.HTML(w, status, name, values)
	log.SetResponseStatus(ctx, status)
}

// Text ... テキストをレンダリングする
func Text(ctx context.Context, w http.ResponseWriter, status int, body string) {
	r := render.New()
	r.Text(w, status, body)
	log.SetResponseStatus(ctx, status)
}

// CSV ... CSV(UTF8)をレンダリングする
func CSV(ctx context.Context, w http.ResponseWriter, name string, data [][]string) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s.csv", name))

	writer := csv.NewWriter(w)
	for _, datum := range data {
		writer.Write(datum)
	}
	writer.Flush()
	log.SetResponseStatus(ctx, http.StatusOK)
}

// CSVByShiftJIS ... CSV(ShiftJIS)をレンダリングする
func CSVByShiftJIS(ctx context.Context, w http.ResponseWriter, name string, data [][]string) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s.csv", name))

	writer := csv.NewWriter(transform.NewWriter(w, japanese.ShiftJIS.NewEncoder()))
	for _, datum := range data {
		writer.Write(datum)
	}
	writer.Flush()
	log.SetResponseStatus(ctx, http.StatusOK)
}
