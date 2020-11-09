package hook

import (
	"context"
	"testing"
)

func TestHookAddValue(t *testing.T) {
	var FooEvent struct{}
	d := NewDefaultDispatcher()
	d.Listen(FooEvent, func(ctx context.Context) error {
		return nil
	})
	err := d.Dispatch(context.Background(), FooEvent)
	if err != nil {
		t.Fail()
	}
}

func TestHookChangeValue(t *testing.T) {
	var FooEvent struct{}
	d := NewDefaultDispatcher()
	_ = d.ListenWithData(FooEvent, func(ctx context.Context, data interface{}) error {
		ptr := data.(*int)
		*ptr = 2
		return nil
	})
	var value = 1
	err := d.DispatchWithData(context.Background(), FooEvent, &value)
	if err != nil {
		t.Fail()
	}
	if value != 2 {
		t.Fatalf("expect 2, got %d", value)
	}
}
