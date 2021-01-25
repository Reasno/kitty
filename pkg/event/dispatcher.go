package event

import (
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"sync"
)

type Dispatcher struct {
	registry map[string][]contract.Listener
	rwLock   sync.RWMutex
}

func (d *Dispatcher) Dispatch(event contract.Event) error {
	d.rwLock.RLock()
	defer d.rwLock.RUnlock()
	listeners, ok := d.registry[event.Type()]

	if !ok {
		return nil
	}
	for _, listener := range listeners {
		if err := listener.Process(event); err != nil {
			return err
		}
	}
	return nil
}

func (d *Dispatcher) Subscribe(listener contract.Listener) {
	d.rwLock.Lock()
	defer d.rwLock.Unlock()

	if d.registry == nil {
		d.registry = make(map[string][]contract.Listener)
	}
	for _, e := range listener.Listen() {
		d.registry[e.Type()] = append(d.registry[e.Type()], listener)
	}
}
