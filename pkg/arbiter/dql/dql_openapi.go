package dql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/errcode"
	"github.com/GuanceCloud/platypus/pkg/token"
)

var (
	_ DQL = (*DQLCliOpenAPI)(nil)
)

type DQLCliOpenAPI struct {
	Endpoint string
	Path     string
	URL      string

	APIKey string
	WSUUID string

	Timerange []int64
}

func NewDQLOpenAPI(endpoint, path, key string, timeRange []int64) *DQLCliOpenAPI {
	u, _ := url.JoinPath(endpoint, path)
	return &DQLCliOpenAPI{
		Endpoint:  endpoint,
		Path:      path,
		URL:       u,
		APIKey:    key,
		Timerange: timeRange,
	}
}

func (cli *DQLCliOpenAPI) Query(pos token.LnColPos, q, qTyp string, limit, offset, slimit int64, timeRange []any, uuids ...string) (map[string]any, error) {
	url := cli.URL
	if url == "" {
		return nil, fmt.Errorf("dql query url is empty")
	}

	query := map[string]any{
		"q":                    q,
		"disable_sampling":     true,
		"limit":                limit,
		"offset":               offset,
		"disableMultipleField": false,
	}

	if slimit > 0 {
		query["slimit"] = slimit
	}

	if len(timeRange) == 2 {
		query["timeRange"] = timeRange
	} else if len(cli.Timerange) == 2 {
		query["timeRange"] = cli.Timerange
	}

	if len(uuids) > 0 {
		query["workspace_uuids"] = uuids
	}

	b, err := json.Marshal(map[string]any{
		"queries": []map[string]any{
			{
				"qtype": qTyp,
				"query": query,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("DF-API-KEY", cli.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := map[string]any{}

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var mp map[string]any
		if err := json.Unmarshal(buf, &mp); err != nil {
			return nil, err
		}

		// result["body"] = mp
		if v, ok := mp["errorCode"]; ok {
			if v, ok := v.(string); ok && v != "" {
				result["error_code"] = v
			}
		} else if v, ok := mp["message"]; ok {
			if v, ok := v.(string); ok && v != "" {
				result["message"] = v
			}
		}
		result["series"] = cli.GetSeries(mp)
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

func (cli *DQLCliOpenAPI) TimeRange() []int64 {
	if len(cli.Timerange) == 2 {
		return append([]int64{}, cli.Timerange...)
	}
	return nil
}

func (cli *DQLCliOpenAPI) GetSeries(resp map[string]any) []any {
	var datas []any
	if v, ok := resp["content"]; ok {
		if v, ok := v.(map[string]any); ok {
			if v, ok := v["data"]; ok {
				if v, ok := v.([]any); ok {
					datas = v
				}
			}
		}
	}

	if len(datas) == 0 {
		return []any{}
	}

	data := datas[0]

	dataMap, ok := data.(map[string]any)
	if !ok {
		return []any{}
	}
	series, ok := dataMap["series"]
	if !ok {
		return []any{}
	}
	seriesLi, ok := series.([]any)
	if !ok {
		return []any{}
	}

	return getSeries(seriesLi)
}

func getSeries(series []any) []any {
	var seriesPts []any
	for _, sElem := range series {
		elem, ok := sElem.(map[string]any)
		if !ok {
			continue
		}
		columns, ok := elem["columns"]
		if !ok {
			continue
		}
		colNames, ok := columns.([]any)
		if !ok {
			continue
		}

		values, ok := elem["values"]
		if !ok {
			continue
		}
		vals, ok := values.([]any)
		if !ok {
			continue
		}

		var tags map[string]any
		if v, ok := elem["tags"]; ok {
			if v, ok := v.(map[string]any); ok {
				tags = v
			}
		}
		if v, ok := elem["name"]; ok {
			if v, ok := v.(string); ok {
				if tags == nil {
					tags = map[string]any{}
				}
				tags["name"] = v
			}
		}

		pts := []any{}

		for _, col := range vals {
			c, ok := col.([]any)
			if !ok {
				continue
			}
			if len(c) != len(colNames) {
				continue
			}
			cols := map[string]any{}
			for i := range c {
				if n, ok := colNames[i].(string); ok {
					cols[n] = c[i]
				}
			}
			pts = append(pts, map[string]any{
				"tags":    tags,
				"columns": cols,
			})

		}
		seriesPts = append(seriesPts, pts)
	}

	return seriesPts
}
