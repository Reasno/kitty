//go:generate tr proto/app.proto --svcout app

package main

import (
	"github.com/Reasno/kitty/cmd"
)

func main() {
	cmd.Execute()
}
