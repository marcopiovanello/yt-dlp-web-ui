package logging

import (
	"context"
)

/*
Logger implementation using the observable pattern.
Implements io.Writer interface.

The observable is an event source which drops everythigng unless there's
a subscriber connected.

The observer implementatios are a http ServerSentEvents handler and a
websocket one in handler.go
*/
type ObservableLogger struct {
	logsChan chan []byte
}

func NewObservableLogger() *ObservableLogger {
	return &ObservableLogger{
		logsChan: make(chan []byte, 100),
	}
}

func (o *ObservableLogger) Write(p []byte) (n int, err error) {
	select {
	case o.logsChan <- p:
		n = len(p)
		err = nil
		return
	default:
		return
	}
}

func (o *ObservableLogger) Observe(ctx context.Context) <-chan string {
	logs := make(chan string)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case logLine := <-o.logsChan:
				logs <- string(logLine)
			}
		}
	}()

	return logs
}
