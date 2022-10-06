package flexiblecsv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"reflect"
)

func MarshalString(srcs interface{}) (string, error) {
	records, err := marshal(srcs)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	for _, record := range records {
		if err := w.Write(record); err != nil {
			return "", err
		}
	}
	w.Flush()

	return buf.String(), nil
}

func MarshalFile(srcs interface{}, file *os.File) error {
	records, err := marshal(srcs)
	if err != nil {
		return err
	}

	w := csv.NewWriter(file)
	for _, record := range records {
		if err := w.Write(record); err != nil {
			return err
		}
	}
	w.Flush()

	return nil
}

func marshal(srcs interface{}) ([][]string, error) {
	// Pointerの場合は参照先を取得
	rvs := reflect.Indirect(reflect.ValueOf(srcs))

	// Sliceかどうか判定
	if rvs.Type().Kind() != reflect.Slice {
		return nil, errors.New("srcs is not slice")
	}

	// 空か判定
	if rvs.Len() == 0 {
		return [][]string{}, nil
	}

	dsts := [][]string{}

	// Headerを取得
	dst := []string{}
	dst = marshalHeader(dst, "", rvs.Index(0))
	dsts = append(dsts, dst)

	// Valueを取得
	for i := 0; i < rvs.Len(); i++ {
		dst := []string{}
		dst = marshalValue(dst, rvs.Index(i))
		dsts = append(dsts, dst)
	}
	return dsts, nil
}

func marshalHeader(dsts []string, prevTag string, rp reflect.Value) []string {
	rv := reflect.Indirect(rp)
	for i := 0; i < rv.Type().NumField(); i++ {
		tag := rv.Type().Field(i).Tag.Get("csv")
		if rv.Field(i).Type().Kind() == reflect.Ptr {
			if rv.Field(i).IsZero() {
				continue
			}
			if tag != "" {
				tag = fmt.Sprintf("%s.", tag)
			}
			dsts = marshalHeader(dsts, tag, rv.Field(i))
		} else {
			dsts = append(dsts, fmt.Sprintf("%s%s", prevTag, tag))
		}
	}
	return dsts
}

func marshalValue(dsts []string, rp reflect.Value) []string {
	rv := reflect.Indirect(rp)
	for i := 0; i < rv.Type().NumField(); i++ {
		if rv.Field(i).Type().Kind() == reflect.Ptr {
			if rv.Field(i).IsZero() {
				continue
			}
			dsts = marshalValue(dsts, rv.Field(i))
		} else {
			dsts = append(dsts, fmt.Sprintf("%v", rv.Field(i).Interface()))
		}
	}
	return dsts
}
