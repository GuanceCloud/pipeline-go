package arbiter

import (
	"context"
	"testing"
	"time"

	funcs "github.com/GuanceCloud/pipeline-go/pkg/arbiter/builtin-funcs"
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

type logWriter struct {
	index string
	data  any
}

func (w *logWriter) Push(index string, data any) (map[string]any, error) {
	w.index = index
	w.data = data
	return map[string]any{"ok": true, "accepted": int64(1), "error": ""}, nil
}

func TestRunWithLogWriter(t *testing.T) {
	writer := &logWriter{}
	err := Run("push_log.p", `push_log({"fields": {"message": "risk detected"}})`,
		WithFuncs(funcs.Funcs),
		WithLogWriter(writer),
	)
	require.NoError(t, err)
	require.Empty(t, writer.index)
	require.Equal(t, map[string]any{
		"fields": map[string]any{"message": "risk detected"},
	}, writer.data)
}
