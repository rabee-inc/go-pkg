package validation

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/go-playground/validator.v9"

	"github.com/rabee-inc/go-pkg/errcode"
)

func ConvertErrorMessage(err error, prefix string, fn func(tag, field, value string) string) error {
	var vErrs validator.ValidationErrors
	var ok bool
	if vErrs, ok = err.(validator.ValidationErrors); !ok {
		return err
	}

	msgs := []string{prefix}
	for _, vErr := range vErrs {
		tag := vErr.Tag()
		field := vErr.Field()
		value := vErr.Param()
		msg := fn(tag, field, value)
		msgs = append(msgs, msg)
	}

	dst := errors.New(strings.Join(msgs, "\n"))
	dst = errcode.Set(dst, http.StatusBadRequest)
	return dst
}

func ConvertErrorMessageByDefault(err error, fn func(field string) string) error {
	var vErrs validator.ValidationErrors
	var ok bool
	if vErrs, ok = err.(validator.ValidationErrors); !ok {
		return err
	}

	msgs := []string{"以下の入力項目を確認してください"}
	for _, vErr := range vErrs {
		var msg string
		tag := vErr.Tag()
		field := vErr.Field()
		if fn != nil {
			field = fn(field)
		}
		value := vErr.Param()
		switch tag {
		case "required":
			msg = fmt.Sprintf("・%s は必須", field)
		case "email":
			msg = fmt.Sprintf("・%s はメールアドレス形式", field)
		case "min":
			msg = fmt.Sprintf("・%s は%s文字以上", field, value)
		case "max":
			msg = fmt.Sprintf("・%s は%s文字以下", field, value)
		case "len":
			msg = fmt.Sprintf("・%s は%s文字固定", field, value)
		case "gt":
			msg = fmt.Sprintf("・%s は%sを超える数", field, value)
		case "gte":
			msg = fmt.Sprintf("・%s は%s以上の数", field, value)
		case "lt":
			msg = fmt.Sprintf("・%s は%s未満の数", field, value)
		case "lte":
			msg = fmt.Sprintf("・%s は%s以下の数", field, value)
		case "numeric":
			msg = fmt.Sprintf("・%s は数字", field)
		case "url":
			msg = fmt.Sprintf("・%s はURL形式", field)
		case "hexcolor":
			msg = fmt.Sprintf("・%s はカラーコード形式", field)
		case "contains":
			msg = fmt.Sprintf("・%s は%sを含む", field, value)
		case "startswith":
			msg = fmt.Sprintf("・%s は %s で始める", field, value)
		case "endswith":
			msg = fmt.Sprintf("・%s は %s で終わる", field, value)
		default:
			// 他に追加したい場合は下記を参照
			// https://godoc.org/gopkg.in/go-playground/validator.v9
		}
		msgs = append(msgs, msg)
	}

	dst := errors.New(strings.Join(msgs, "\n"))
	dst = errcode.Set(dst, http.StatusBadRequest)
	return dst
}
