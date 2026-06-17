package dql

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GuanceCloud/platypus/pkg/token"
	"github.com/stretchr/testify/require"
)

func TestDQLCliKodoAlignTime(t *testing.T) {
	var bodies []map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		var got map[string]any
		require.NoError(t, json.Unmarshal(body, &got))
		bodies = append(bodies, got)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"","error_code":"","content":[{"Series":[],"query_progress":{"scanned_compressed_bytes":0,"scanned_rows":0}}]}`))
	}))
	defer srv.Close()

	cli := NewDQLKodo(srv.URL, "wksp", nil)
	_, err := cli.Query(token.LnColPos{}, "M::cpu", "dql", 10, 0, 0, nil, true, true, false)
	require.NoError(t, err)
	_, err = cli.Query(token.LnColPos{}, "M::cpu", "dql", 10, 0, 0, nil, false, false, true)
	require.NoError(t, err)

	require.Len(t, bodies, 2)
	require.Equal(t, true, kodoQuery(t, bodies[0])["align_time"])
	require.Equal(t, true, kodoQuery(t, bodies[0])["disable_sampling"])
	require.Equal(t, false, kodoQuery(t, bodies[0])["disable_streaming_aggregation"])
	require.Equal(t, false, kodoQuery(t, bodies[1])["align_time"])
	require.Equal(t, false, kodoQuery(t, bodies[1])["disable_sampling"])
	require.Equal(t, true, kodoQuery(t, bodies[1])["disable_streaming_aggregation"])
}

func TestDQLCliOpenAPIAlignTime(t *testing.T) {
	var bodies []map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		var got map[string]any
		require.NoError(t, json.Unmarshal(body, &got))
		bodies = append(bodies, got)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"content":{"data":[{"series":[]}]},"errorCode":"","message":"","success":true}`))
	}))
	defer srv.Close()

	cli := NewDQLOpenAPI(srv.URL, OpenAPIPath, "key", nil)
	_, err := cli.Query(token.LnColPos{}, "M::cpu", "dql", 10, 0, 0, nil, true, true, false)
	require.NoError(t, err)
	_, err = cli.Query(token.LnColPos{}, "M::cpu", "dql", 10, 0, 0, nil, false, false, true)
	require.NoError(t, err)

	require.Len(t, bodies, 2)
	require.Equal(t, true, openAPIQuery(t, bodies[0])["alignTime"])
	require.Equal(t, true, openAPIQuery(t, bodies[0])["disable_sampling"])
	require.Nil(t, openAPIQuery(t, bodies[0])["disableStreamingAggregation"])
	require.Equal(t, false, openAPIQuery(t, bodies[1])["alignTime"])
	require.Equal(t, false, openAPIQuery(t, bodies[1])["disable_sampling"])
	require.Equal(t, true, openAPIQuery(t, bodies[1])["disableStreamingAggregation"])
}

func kodoQuery(t *testing.T, body map[string]any) map[string]any {
	t.Helper()
	queries, ok := body["queries"].([]any)
	require.True(t, ok)
	require.Len(t, queries, 1)
	query, ok := queries[0].(map[string]any)
	require.True(t, ok)
	return query
}

func openAPIQuery(t *testing.T, body map[string]any) map[string]any {
	t.Helper()
	queries, ok := body["queries"].([]any)
	require.True(t, ok)
	require.Len(t, queries, 1)
	queryWrap, ok := queries[0].(map[string]any)
	require.True(t, ok)
	query, ok := queryWrap["query"].(map[string]any)
	require.True(t, ok)
	return query
}
