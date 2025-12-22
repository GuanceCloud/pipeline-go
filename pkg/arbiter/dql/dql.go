package dql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/errcode"
	"github.com/GuanceCloud/platypus/pkg/token"
)

const (
	OpenAPIPath = "/api/v1/df/query_data_v1"
	KodoPath    = "/v1/query"
)

var (
	client = &http.Client{
		Transport: &http.Transport{
			DialContext: ((&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext),
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: time.Minute * 5,
	}
)

var (
	_ DQL = (*DQLCliKodo)(nil)
)

type DQL interface {
	Query(pos token.LnColPos, q, qTyp string, limit, offset, slimit int64, timeRange []any, uuids ...string) (map[string]any, error)
	TimeRange() []int64
}

type DQLCliKodo struct {
	URL       string
	WSUUID    string
	Timerange []int64

	MetricTotal *ProgressMetric
}

func NewDQLKodo(url, uuid string, timeRange []int64) *DQLCliKodo {
	return &DQLCliKodo{
		URL:         url,
		WSUUID:      uuid,
		Timerange:   timeRange,
		MetricTotal: &ProgressMetric{},
	}
}

func (cli *DQLCliKodo) TimeRange() []int64 {
	if len(cli.Timerange) == 2 {
		return append([]int64{}, cli.Timerange...)
	}
	return nil
}

func (cli *DQLCliKodo) Query(pos token.LnColPos, q, qTyp string, limit, offset, slimit int64, timeRange []any, uuids ...string) (map[string]any, error) {
	url := cli.URL
	if url == "" {
		return nil, fmt.Errorf("dql query url is empty")
	}
	var curMetric ProgressMetric

	query := map[string]any{
		"query":                         q,
		"qtype":                         qTyp,
		"limit":                         limit,
		"offset":                        offset,
		"disable_sampling":              true,
		"disable_multiple_field":        false,
		"disable_streaming_aggregation": true,
		"disable_slimit":                true,
		"mask_visible":                  true,
	}

	if slimit > 0 {
		query["slimit"] = slimit
	}

	if len(timeRange) == 2 {
		query["time_range"] = timeRange
	} else if len(cli.Timerange) == 2 {
		query["time_range"] = cli.Timerange
	}

	if len(uuids) > 0 {
		query["workspace_uuids"] = uuids
	}

	b, err := json.Marshal(map[string]any{
		"workspace_uuid": cli.WSUUID,
		// "query_source":   "arbiter",
		"queries": []map[string]any{
			query,
		},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "arbiter")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := map[string]any{}

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		qResp := &QueryResultKodo{}
		if err := json.NewDecoder(resp.Body).Decode(qResp); err != nil {
			return nil, err
		}
		if qResp.ErrorCode != nil {
			result["error_code"] = qResp.ErrorCode
		}
		if qResp.Message != nil {
			result["message"] = qResp.Message
		}

		result["series"] = cli.GetSeries(qResp, &curMetric)

		cli.MetricTotal.ScannedCompressedBytes += curMetric.ScannedCompressedBytes
		cli.MetricTotal.ScannedRows += curMetric.ScannedRows
		cli.MetricTotal.Count++
	} else {
		result["error_code"] = errcode.ArbiterFnErr
		result["message"] = fmt.Sprintf(
			"expected Content-Type of the request response was `%s`, but it was actually `%s`",
			"application/json",
			contentType,
		)
		result["series"] = []any{}
	}

	result["status_code"] = resp.StatusCode

	return result, nil
}

// GetSeries returns series, scanned_compressed_bytes and scanned_rows
func (cli *DQLCliKodo) GetSeries(qResult *QueryResultKodo, mertic *ProgressMetric) []any {
	if len(qResult.Content) == 0 || qResult.Content[0] == nil {
		return []any{}
	}

	data := qResult.Content[0]

	mertic.ScannedCompressedBytes = data.QueryProgress.ScannedCompressedBytes
	mertic.ScannedRows = data.QueryProgress.ScannedRows
	mertic.Count++

	return getSeries(data.Series)
}

type ProgressMetric struct {
	ScannedCompressedBytes int64
	ScannedRows            int64

	Count int64
}

type QueryResultKodo struct {
	Message   any             `json:"message"`
	ErrorCode any             `json:"error_code"`
	Content   []*QueryContent `json:"content"`
}

type QueryContent struct {
	Series        []any          `json:"Series"`
	QueryProgress *QueryProgress `json:"query_progress"`
}

type Series struct {
	Columns []string `json:"columns"`
	Values  []any    `json:"values"`
}

type QueryProgress struct {
	ScannedCompressedBytes int64 `json:"scanned_compressed_bytes"`
	ScannedRows            int64 `json:"scanned_rows"`
}
