package cloudfirestore

import "context"

type FuncGetID[D any] func(*D) string
type FuncConvert[S, D any] func(*S) *D
type ConvertibleBatchGetterItem[D any] interface {
	// After ... 取得後に実行されるコールバック関数を追加(コールバック関数の第一引数に nil が入ることはありません)。
	// After のコールバック関数内で Add した場合は、再度 Commit が実行されます。
	After(func(D)) ConvertibleBatchGetterItem[D]
	// RemoveAfter ... After で登録したコールバック関数を削除
	RemoveAfter() ConvertibleBatchGetterItem[D]
	// OnEmpty ... 取得後に nil だった場合に実行されるコールバック関数を追加。
	// OnEmpty のコールバック関数内で Add した場合は、再度 Commit が実行されます。
	OnEmpty(func()) ConvertibleBatchGetterItem[D]
	// RemoveOnEmpty ... OnEmpty で登録したコールバック関数を削除
	RemoveOnEmpty() ConvertibleBatchGetterItem[D]
}

type ConvertibleBatchGetter[S, D any] interface {
	// Add ... 取得対象を追加
	Add(ids ...string) ConvertibleBatchGetterItem[*D]
	// Delete ... 取得対象を削除
	Delete(ids ...string)
	// GetMap ... 取得済みのデータをmapで取得
	// map の key は doc の path になります。
	GetMap() map[string]*D
	// Get ... 取得済みのデータを取得
	Get(ids ...string) *D
	// Commit ... alias: bg の Commit を実行する
	Commit(ctx context.Context) error
	// Set ... 取得後に値を入れる変数を指定して取得対象を追加。
	// 第二引数以降の ids は Add の引数から dst.ID を除いたものを指定する。
	Set(dst *D, ids ...string) ConvertibleBatchGetterItem[*D]
	// SetWithID ... 取得後に値を入れる変数を指定して取得対象を追加。
	// 第二引数以降の ids は Add の引数と同じ。
	SetWithID(dst *D, ids ...string) ConvertibleBatchGetterItem[*D]
}
