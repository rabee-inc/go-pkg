package cloudfirestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

type BatchGetter interface {
	// Add ... 取得対象を追加
	Add(docRef *firestore.DocumentRef, dst any)
	// Delete ... 取得対象を削除
	Delete(docRef *firestore.DocumentRef)
	// Commit ... 取得処理を実行
	Commit(ctx context.Context) error
	// IsCommittedItem ... 対象データがコミットが済んでいるかどうか
	IsCommittedItem(path string) bool
	// Get ... 対象データを取得
	Get(path string) any
	// OnCommit ... Commit 後に実行されるコールバック関数を追加。
	// このコールバック関数内で、Addをした場合は、再度コミットされます。(OnCommit も再度発火します。)
	OnCommit(func()) func()
	// OnEnd ... Commit 完了後に実行されるコールバック関数を追加。
	// 最後に一度だけ実行されます。
	OnEnd(func()) func()
}
