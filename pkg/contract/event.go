//go:generate mockery --name=Dispatcher
package contract

import "context"

type Event interface {
	Type() string
	Data() interface{}
	Context() context.Context
}

type Dispatcher interface {
	Dispatch(event Event) error
	Subscribe(listener Listener)
}

type Listener interface {
	Listen() []Event
	Process(event Event) error
}
