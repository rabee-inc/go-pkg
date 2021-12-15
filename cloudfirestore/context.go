package cloudfirestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

type contextKey string

const (
	ctxTxKey contextKey = "firestore:tx"
	ctxBtKey contextKey = "firestore:bt"
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

func getContextWriteBatch(ctx context.Context) *firestore.WriteBatch {
	if bt, ok := ctx.Value(ctxBtKey).(*firestore.WriteBatch); ok {
		return bt
	}
	return nil
}

func setContextWriteBatch(ctx context.Context, bt *firestore.WriteBatch) context.Context {
	return context.WithValue(ctx, ctxBtKey, bt)
}
