package kerr

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func TestServerError_Error(t *testing.T) {
	err := errors.New("whoops")
	fmt.Printf("%+v\n", err)
}
