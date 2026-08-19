package funcs

import (
	"errors"
	"testing"

	"github.com/GuanceCloud/platypus/pkg/engine"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type logWriterMock struct {
	index string
	data  any
	err   error
}

func (w *logWriterMock) Push(index string, data any) (map[string]any, error) {
	w.index = index
	w.data = data
	if w.err != nil {
		return nil, w.err
	}
	return map[string]any{
		"ok":       true,
		"accepted": int64(1),
		"error":    "",
	}, nil
}

func TestPushLog(t *testing.T) {
	for _, tc := range cPushLog.Progs {
		t.Run(tc.Name, func(t *testing.T) {
			writer := &logWriterMock{}
			runCase(t, tc, map[runtimev2.TaskP]any{PLogWriter: writer})

			data, ok := writer.data.(map[string]any)
			require.True(t, ok)
			assert.Equal(t, "risk_check", data["source"])
			fields, ok := data["fields"].(map[string]any)
			require.True(t, ok)
			assert.Equal(t, "risk detected", fields["message"])
			if tc.Name == "push_log_default_index" {
				assert.Empty(t, writer.index)
			} else {
				assert.Equal(t, "security", writer.index)
			}
		})
	}
}

func TestPushLogReturnsWriterError(t *testing.T) {
	writer := &logWriterMock{err: errors.New("kodo unavailable")}
	prog := `result = push_log({"fields": {"message": "risk detected"}})
printf("%v", result)`
	script, err := engine.ParseV2("push_log_error", prog, Funcs)
	require.NoError(t, err)

	stdout := &testWriter{}
	if runErr := script.Run(nil, runtimev2.WithPrivate(map[runtimev2.TaskP]any{
		PLogWriter: writer,
		PStdout:    stdout,
	})); runErr != nil {
		t.Fatal(runErr)
	}
	assert.JSONEq(t, `{"accepted":0,"error":"kodo unavailable","ok":false}`, stdout.String())
}

func TestPushLogRequiresWriter(t *testing.T) {
	script, err := engine.ParseV2("push_log_missing_writer",
		`push_log({"fields": {"message": "risk detected"}})`, Funcs)
	require.NoError(t, err)

	err = script.Run(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), string(PLogWriter))
}

type testWriter struct {
	data []byte
}

func (w *testWriter) Write(p []byte) (int, error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

func (w *testWriter) String() string {
	return string(w.data)
}
