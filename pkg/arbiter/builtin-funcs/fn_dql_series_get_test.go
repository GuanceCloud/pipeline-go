package funcs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/dql"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
)

func TestFuncDQLSeries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(mockOpenAPIMetric))
	defer server.Close()
	cases := []ProgCase{}
	cases = append(cases, cDQLSeriesGet.Progs...)
	cases = append(cases, cDQLSeriesFirst.Progs...)
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			runCase(t, tc, map[runtimev2.TaskP]any{
				PDQLCli: dql.NewDQLOpenAPI(
					server.URL,
					dql.OpenAPIPath,
					"abc", nil,
				),
			})
		})
	}
}
