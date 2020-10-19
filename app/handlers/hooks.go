package handlers

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

var jaegerCloser io.Closer

func InterruptHandler(errc chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	terminateError := fmt.Errorf("%s", <-c)
	// Place whatever shutdown handling you want here
	if jaegerCloser != nil {
		_ = jaegerCloser.Close()
	}

	errc <- terminateError
}
