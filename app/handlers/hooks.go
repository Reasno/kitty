package handlers

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type destructor []func()

func (d *destructor) Add(close ...func()) {
	*d = append(*d, close...)
}
func (d *destructor) Close() {
	for _, c := range *d {
		c()
	}
}

var destruct destructor

func InterruptHandler(errc chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	terminateError := fmt.Errorf("%s", <-c)
	// Place whatever shutdown handling you want here
	destruct.Close()

	errc <- terminateError
}
