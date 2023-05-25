package util

import "github.com/rabee-inc/go-pkg/maputil"

type eventListener struct {
	f func()
}

type eventEmitter struct {
	listeners maputil.OrderedSet[*eventListener]
}

// NewEventEmitter ... EventEmitter を生成する
func NewEventEmitter() EventEmitter {
	return &eventEmitter{
		listeners: maputil.NewOrderedSet[*eventListener](nil),
	}
}

func (e *eventEmitter) Add(f func()) func() {
	el := &eventListener{f: f}
	e.listeners.Add(el)
	return func() {
		e.listeners.Delete(el)
	}
}

func (e *eventEmitter) AddOnce(f func()) func() {
	var el *eventListener
	once := func() {
		if f != nil {
			f()
		}
		e.listeners.Delete(el)
	}
	el = &eventListener{f: once}

	e.listeners.Add(el)
	return once
}

func (e *eventEmitter) Clear() {
	e.listeners.Clear()
}

func (e *eventEmitter) Emit() {
	for _, el := range e.listeners.Keys() {
		el.f()
	}
}
