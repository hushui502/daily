package runner

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

// SignalError is returned by the signal handler's execute function
// when it terminates due to a received signal.
type SignalError struct {
	Signal os.Signal
}

func SignalHandler(ctx context.Context, signals ...os.Signal) (excute func() error, interrupt func(error)) {
	ctx, cancel := context.WithCancel(ctx)

	return func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, signals...)
			defer signal.Stop(c)

			select {
			case sig := <-c:
				return SignalError{Signal: sig}
			case <-ctx.Done():
				return ctx.Err()
			}
		}, func(error) {
			cancel()
		}
}

func (s SignalError) Error() string {
	return fmt.Sprintf("received signal %s", s.Signal)
}
