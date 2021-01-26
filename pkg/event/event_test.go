package event

import (
	"context"
	"fmt"
	"testing"
)

type TestE struct {
}

func TestEvent(t *testing.T) {
	testE := NewEvent(context.Background(), TestE{})
	fmt.Println(testE.Type())
}
