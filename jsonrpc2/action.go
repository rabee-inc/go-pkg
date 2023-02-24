package jsonrpc2

import (
	"context"
	"encoding/json"
)

type Action interface {
	DecodeParams(
		ctx context.Context,
		msg *json.RawMessage,
	) (any, error)

	Exec(
		ctx context.Context,
		method string,
		params any,
	) (any, error)
}
