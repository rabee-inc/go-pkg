package jsonrpc2

import (
	"context"
	"encoding/json"
)

type Action interface {
	DecodeParams(
		ctx context.Context,
		msg *json.RawMessage,
	) (interface{}, error)

	Exec(
		ctx context.Context,
		method string,
		params interface{},
	) (interface{}, error)
}
