package util

type EventEmitterWithParam[T any] interface {
	// Add ... イベントリスナーを追加する。
	// 追加されたリスナーを解除する関数を返します
	Add(func(T)) func()
	// AddOnce ... イベントリスナーを一度だけ実行するように追加する。
	// 追加されたリスナーを解除する関数を返します
	AddOnce(func(T)) func()
	// Clear ... イベントリスナーを全て削除する
	Clear()
	// Emit ... イベントを発火する
	Emit(T)
}
