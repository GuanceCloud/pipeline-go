package arbiter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type contextExitSignal struct {
	ctx context.Context
}

func (s contextExitSignal) ExitSignal() bool {
	return s.ctx != nil && s.ctx.Err() != nil
}

func TestRunStopsWhenSignalExits(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	errCh := make(chan error, 1)
	start := time.Now()
	go func() {
		errCh <- Run("test.p", `for ;; {}`, func(c *Config) {
			c.Signal = contextExitSignal{ctx: ctx}
		})
	}()

	select {
	case err := <-errCh:
		require.NoError(t, err)
		require.Less(t, time.Since(start), time.Second)
	case <-time.After(time.Second):
		t.Fatal("arbiter runtime did not stop after signal exit")
	}
}
