package util

type EventEmitter interface {
	// Add ... イベントリスナーを追加する。
	// 追加されたリスナーを解除する関数を返します
	Add(f func()) func()
	// AddOnce ... イベントリスナーを一度だけ実行するように追加する。
	// 追加されたリスナーを解除する関数を返します
	AddOnce(f func()) func()
	// Clear ... イベントリスナーを全て削除する
	Clear()
	// Emit ... イベントを発火する
	Emit()
}
