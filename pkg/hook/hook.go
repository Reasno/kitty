package hook

import (
	"context"
	"sync"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, event interface{}) error
	Listen(event interface{}, f HandleFunc) error
}

type HandleFunc func(ctx context.Context) error

type defaultDispatcher struct {
	rwLock   sync.RWMutex
	internal map[interface{}][]HandleFunc
}

func NewDefaultDispatcher() *defaultDispatcher {
	return &defaultDispatcher{
		rwLock:   sync.RWMutex{},
		internal: make(map[interface{}][]HandleFunc),
	}
}

func (d *defaultDispatcher) Dispatch(ctx context.Context, event interface{}) error {
	d.rwLock.RLock()
	defer d.rwLock.RUnlock()
	for _, f := range d.internal[event] {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *defaultDispatcher) Listen(event interface{}, handler HandleFunc) error {
	d.rwLock.Lock()
	defer d.rwLock.Unlock()
	d.internal[event] = append(d.internal[event], handler)
	return nil
}

func (d *defaultDispatcher) DispatchWithData(ctx context.Context, event interface{}, data interface{}) error {
	ctx = context.WithValue(ctx, event, data)
	return d.Dispatch(ctx, event)
}

type dataHandler struct {
	event       interface{}
	trueHandler func(ctx context.Context, data interface{}) error
}

func (d *defaultDispatcher) ListenWithData(event interface{}, f func(ctx context.Context, data interface{}) error) error {
	return d.Listen(event, dataHandler{
		event:       event,
		trueHandler: f,
	}.Handle)
}

func (d dataHandler) Handle(ctx context.Context) error {
	data := ctx.Value(d.event)
	return d.trueHandler(ctx, data)
}
