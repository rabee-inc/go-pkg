package cloudfirestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

type contextKey string

const (
	ctxTxKey contextKey = "firestore:tx"
	ctxBwKey contextKey = "firestore:bw"
)

func getContextTransaction(ctx context.Context) *firestore.Transaction {
	if tx, ok := ctx.Value(ctxTxKey).(*firestore.Transaction); ok {
		return tx
	}
	return nil
}

func setContextTransaction(ctx context.Context, tx *firestore.Transaction) context.Context {
	return context.WithValue(ctx, ctxTxKey, tx)
}

func getContextBulkWriter(ctx context.Context) *firestore.BulkWriter {
	if bw, ok := ctx.Value(ctxBwKey).(*firestore.BulkWriter); ok {
		return bw
	}
	return nil
}

func setContextBulkWriter(ctx context.Context, bt *firestore.BulkWriter) context.Context {
	return context.WithValue(ctx, ctxBwKey, bt)
}
