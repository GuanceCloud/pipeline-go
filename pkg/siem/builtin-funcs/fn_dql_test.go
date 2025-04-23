package funcs

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GuanceCloud/pipeline-go/pkg/siem/dql"
	"github.com/GuanceCloud/platypus/pkg/engine"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
)

func mockOpenAPIMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"code":200,"content":{"data":[{"async_id":"","column_names":["user","total"],"complete":false,"cost":"7ms","datasource":"raw","group_by":["cpu","guance_site","host","host_ip","project"],"index_name":"","index_names":"","index_store_type":"","interval":0,"is_cross_ws":true,"is_running":false,"next_cursor_time":-1,"points":null,"query_parse":{"fields":{"total":"usage_total","user":"usage_user"},"funcs":{},"namespace":"metric","sources":{"cpu":"exact"}},"query_progress":{"nanos_from_started":0,"nanos_from_submitted":0,"nanos_to_finish":0,"total_percentage":0},"query_status":"","query_type":"guancedb","sample":1,"scan_completed":false,"scan_index":"","series":[{"column_names":["time","user","total"],"columns":["time","user","total"],"name":"cpu","tags":{"cpu":"cpu-total","guance_site":"testing","host":"172.16.241.111","host_ip":"172.16.241.111","name":"cpu","project":"cloudcare-testing"},"units":[null,null,null],"values":[[1744866108991,4.77876106,7.18078381],[1744866103991,7.17009916,10.37376049]]},{"column_names":["time","user","total"],"columns":["time","user","total"],"name":"cpu","tags":{"cpu":"cpu-total","guance_site":"testing","host":"172.16.242.112","host_ip":"172.16.242.112","name":"cpu","project":"cloudcare-testing"},"units":[null,null,null],"values":[[1744866107975,5.69187959,21.75562864],[1744866102975,5.28589581,16.59466328]]}],"window":0}],"declaration":{"business":"test","organization":"636c5b928b83720007027f9a","test":"1,2","托尔斯泰":"23","部门":"产研一部"}},"errorCode":"","message":"","success":true,"traceId":"9344913754735199412"}`))
}

func TestFnDQL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(mockOpenAPIMetric))
	defer server.Close()
	cases := []ProgCase{}
	cases = append(cases, cDQL.Progs...)
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

func mockHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, `{"code":200,"content":{"data":[{"async_id":"","column_names":["usage_user","usage_total","usage_iowait","usage_softirq","usage_irq","usage_system","usage_idle","usage_guest_nice","usage_steal","usage_guest","usage_nice","load5s","host","cluster_name_k8s","cpu","guance_site","host_ip","real_host","region","project"],"complete":false,"cost":"29ms","datasource":"","index_name":"","index_names":"","index_store_type":"","interval":0,"is_cross_ws":true,"is_running":false,"next_cursor_time":-1,"points":null,"query_parse":{"fields":{},"funcs":{},"namespace":"metric","sources":{"cpu":"exact"}},"query_progress":{"nanos_from_started":0,"nanos_from_submitted":0,"nanos_to_finish":0,"total_percentage":0},"query_status":"","query_type":"guancedb","sample":1,"scan_completed":false,"scan_index":"","series":[{"column_names":["time","usage_user","usage_total","usage_iowait","usage_softirq","usage_irq","usage_system","usage_idle","usage_guest_nice","usage_steal","usage_guest","usage_nice","load5s","host","cluster_name_k8s","cpu","guance_site","host_ip","real_host","region","project"],"columns":["time","usage_user","usage_total","usage_iowait","usage_softirq","usage_irq","usage_system","usage_idle","usage_guest_nice","usage_steal","usage_guest","usage_nice","load5s","host","cluster_name_k8s","cpu","guance_site","host_ip","real_host","region","project"],"name":"cpu","units":[null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null,null],"values":[[1744541030005,3.70182556,5.42596349,null,0.35496957,0,1.34381339,94.57403653,0,0,null,0,0,"172.16.242.111",null,"cpu-total","testing","172.16.242.111",null,null,"cloudcare-testing"],[1744541029760,null,null,null,null,0,null,null,null,null,0,null,null,"cluster_a_172.16.231.101","k8s-daily","cpu-total","daily","172.16.231.101","sh-ve-daily-dataflux-k8s-001",null,null],[1744541027723,5.97052845,8.79065041,0.22865854,0.81300813,0,1.77845528,91.20934959,0,0,0,0,0,"172.16.241.112",null,"cpu-total","testing","172.16.241.112",null,null,"cloudcare-testing"]]}],"window":0}],"declaration":{"business":"test","organization":"636c5b928b83720007027f9a","test":"1,2","托尔斯泰":"23","部门":"产研一部"}},"errorCode":"","message":"","success":true,"traceId":"1721846225271127942"}`)
}

func mockKodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, `{"code":200,"content":{"data":[{"async_id":"","column_names":["load5s","usage_user"],"complete":false,"cost":"54ms","datasource":"","group_by":["cpu","guance_site","host","host_ip","project"],"index_name":"","index_names":"","index_store_type":"","interval":0,"is_cross_ws":true,"is_running":false,"next_cursor_time":-1,"points":null,"query_parse":{"fields":{"load5s":"load5s","usage_user":"usage_user"},"funcs":{},"namespace":"metric","sources":{"cpu":"exact"}},"query_progress":{"nanos_from_started":0,"nanos_from_submitted":0,"nanos_to_finish":0,"total_percentage":0},"query_status":"","query_type":"guancedb","sample":1,"scan_completed":false,"scan_index":"","series":[{"column_names":["time","load5s","usage_user"],"columns":["time","load5s","usage_user"],"name":"cpu","tags":{"cpu":"cpu-total","guance_site":"testing","host":"172.16.241.111","host_ip":"172.16.241.111","name":"cpu","project":"cloudcare-testing"},"units":[null,null,null],"values":[[1744535521158,0,8.54830551],[1744535526158,2,9.68314322]]},{"column_names":["time","load5s","usage_user"],"columns":["time","load5s","usage_user"],"name":"cpu","tags":{"cpu":"cpu-total","host":"sh-ve-sit-doris-be-001","name":"cpu","region":"testing"},"units":[null,null,null],"values":[[1744535528182,1,2.49812014],[1744535538182,0,2.85738157]]}],"window":0}],"declaration":{"business":"test","organization":"636c5b928b83720007027f9a","test":"1,2","托尔斯泰":"23","部门":"产研一部"}},"errorCode":"","message":"","success":true,"traceId":"3850506999914266370"}`)
}

func TestFuncDQLLog(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer server.Close()

	svcKodo := httptest.NewServer(http.HandlerFunc(mockKodo))

	cases := []struct {
		name    string
		script  string
		private map[runtimev2.TaskP]any
	}{
		{
			name:   "test",
			script: "v, ok = dql(\"M::cpu limit 3 slimit 3\"); ;printf(\"%v\\n\",v)",
			private: map[runtimev2.TaskP]any{
				PDQLCli: dql.NewDQLOpenAPI(
					server.URL,
					dql.OpenAPIPath,
					"abc", nil,
				)},
		},
		{
			name:   "test",
			script: "v, ok = dql(\"M::cpu limit 3 slimit 3\"); ;printf(\"%v\\n\",v)",
			private: map[runtimev2.TaskP]any{
				PDQLCli: dql.NewDQLOpenAPI(
					svcKodo.URL,
					dql.KodoPath,
					"abc", nil,
				)},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("%d_%s", i, c.name), func(t *testing.T) {
			s, err := engine.ParseV2(c.name, c.script, Funcs)
			if err != nil {
				t.Error(err)
				return
			}
			if c.private == nil {
				c.private = map[runtimev2.TaskP]any{}
			}
			stdout := bytes.NewBuffer([]byte{})
			c.private[PStdout] = stdout
			if err := s.Run(nil, runtimev2.WithPrivate(c.private)); err != nil {
				t.Error(err.Error())
			}
			o := stdout.String()
			t.Log(c.script)
			t.Log(o)
		})
	}
}
