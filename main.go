//go:generate tr proto/app.proto
package main

import (
	"github.com/Reasno/kitty/cmd"
)

func main() {
	cmd.Execute()
}
