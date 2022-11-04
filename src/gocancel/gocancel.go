package gocancel

import (
	"context"
)

// ctx cancel stops this go routine.
func GoRoutineWithContextCancel(ctx context.Context, f func()) {
	f()
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}
