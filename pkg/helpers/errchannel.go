// Package helpers some simple helpers
package helpers

import (
	"context"

	"github.com/chapsuk/wait"
)

// RegisterErrorChannel create goroutine with error channel return
func RegisterErrorChannel(fn func() error) <-chan error {
	ch := make(chan error)
	wg := wait.Group{}
	wg.Add(func() {
		ch <- fn()
	})

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func RegisterErrorChannelContext(ctx context.Context, fn func(_ context.Context) error) <-chan error {
	ch := make(chan error)
	wg := wait.Group{}
	wg.Add(func() {
		ch <- fn(ctx)
	})

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}
