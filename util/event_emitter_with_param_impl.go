package util

import "github.com/rabee-inc/go-pkg/maputil"

type eventListenerWithParam[T any] struct {
	f func(T)
}

type eventEmitterWithParam[T any] struct {
	listeners maputil.OrderedSet[*eventListenerWithParam[T]]
}

// NewEventEmitterWithParam ... パラメータありの EventEmitter を生成する
func NewEventEmitterWithParam[T any]() EventEmitterWithParam[T] {
	return &eventEmitterWithParam[T]{
		listeners: maputil.NewOrderedSet[*eventListenerWithParam[T]](nil),
	}
}

func (e *eventEmitterWithParam[T]) Add(f func(T)) func() {
	el := &eventListenerWithParam[T]{f: f}
	e.listeners.Add(el)
	return func() {
		e.listeners.Delete(el)
	}
}

func (e *eventEmitterWithParam[T]) AddOnce(f func(T)) func() {
	var el *eventListenerWithParam[T]
	el = &eventListenerWithParam[T]{f: func(v T) {
		if f != nil {
			f(v)
		}
		e.listeners.Delete(el)
	}}
	e.listeners.Add(el)
	return func() {
		e.listeners.Delete(el)
	}
}

func (e *eventEmitterWithParam[T]) Clear() {
	e.listeners.Clear()
}

func (e *eventEmitterWithParam[T]) Emit(v T) {
	for _, el := range e.listeners.Keys() {
		el.f(v)
	}
}
