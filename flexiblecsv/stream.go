package flexiblecsv

import (
	"encoding/csv"
	"os"
	"reflect"
)

type CSVWriter interface {
	Write(row []string) error
	Flush()
	Error() error
}

type StreamMarshaler interface {
	Marshal(row interface{}) error
	Flush() error
}

type streamMarshaler struct {
	isFirst bool
	writer  CSVWriter
}

func NewStreamMarshaler(w CSVWriter) StreamMarshaler {
	return &streamMarshaler{
		isFirst: true,
		writer:  w,
	}
}

func MarshalFileAsStream(file *os.File) StreamMarshaler {
	return NewStreamMarshaler(csv.NewWriter(file))
}

func (s *streamMarshaler) Marshal(row interface{}) error {
	if s.isFirst {
		// Headerを取得
		header := []string{}
		header = marshalHeader(header, "", reflect.ValueOf(row))
		if err := s.writer.Write(header); err != nil {
			return err
		}
		s.isFirst = false
	}

	record := []string{}
	record = marshalValue(record, reflect.ValueOf(row))
	if err := s.writer.Write(record); err != nil {
		return err
	}
	return nil
}

func (s *streamMarshaler) Flush() error {
	s.writer.Flush()
	return s.writer.Error()
}
