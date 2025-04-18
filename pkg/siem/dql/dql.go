package dql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/GuanceCloud/pipeline-go/pkg/siem/errcode"
)

const (
	OpenAPIPath = "/api/v1/df/query_data_v1"
	KodoPath    = "/v1/query"
)

func dialCtx(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}

var client = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: dialCtx(&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}),

		ForceAttemptHTTP2: true,
		MaxIdleConns:      100,
		MaxConnsPerHost:   19200,

		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

type DQL interface {
	Query(q, qTyp string, limit, offset, slimit int64, timeRange []any) (map[string]any, error)
	GetSeries(resp map[string]any) []any // []points
}

var (
	_ DQL = (*DQLCliKodo)(nil)
	_ DQL = (*DQLCliOpenAPI)(nil)
)

type DQLCliOpenAPI struct {
	Endpoint string
	Path     string
	URL      string

	APIKey string

	WSUUID string
}

type DQLCliKodo struct {
	Endpoint string
	Path     string
	URL      string

	WSToken string

	WSUUID string
}

func NewDQLKodo(endpoint, path, token, uuid string) *DQLCliKodo {
	u, _ := url.JoinPath(endpoint, path)
	return &DQLCliKodo{
		Endpoint: endpoint,
		Path:     path,
		URL:      u,
		WSToken:  token,
		WSUUID:   uuid,
	}
}

func NewDQLOpenAPI(endpoint, path, key string) *DQLCliOpenAPI {
	u, _ := url.JoinPath(endpoint, path)
	return &DQLCliOpenAPI{
		Endpoint: endpoint,
		Path:     path,
		URL:      u,
		APIKey:   key,
	}
}

func (cli *DQLCliKodo) Query(q, qTyp string, limit, offset, slimit int64, timeRange []any) (map[string]any, error) {
	url := cli.URL
	if url == "" {
		return nil, fmt.Errorf("dql query url is empty")
	}

	query := map[string]any{
		"query":                  q,
		"qtype":                  qTyp,
		"disable_sampling":       true,
		"limit":                  limit,
		"offset":                 offset,
		"slimit":                 slimit,
		"disable_multiple_field": false,
	}

	if len(timeRange) == 2 {
		query["timeRange"] = timeRange
	}

	b, err := json.Marshal(map[string]any{
		"workspace_uuid": cli.WSUUID,
		"token":          cli.WSToken,
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
		if v, ok := mp["error_code"]; ok {
			result["error_code"] = v
		}
		if v, ok := mp["message"]; ok {
			result["message"] = v
		}
		result["series"] = cli.GetSeries(mp)
	} else {
		result["error_code"] = errcode.SIEMFnErr
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

func (cli *DQLCliOpenAPI) Query(q, qTyp string, limit, offset, slimit int64, timeRange []any) (map[string]any, error) {
	url := cli.URL
	if url == "" {
		return nil, fmt.Errorf("dql query url is empty")
	}

	query := map[string]any{
		"q":                    q,
		"disable_sampling":     true,
		"limit":                limit,
		"offset":               offset,
		"slimit":               slimit,
		"disableMultipleField": false,
	}

	if len(timeRange) == 2 {
		query["timeRange"] = timeRange
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
		result["error_code"] = errcode.SIEMFnErr
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

func (cli *DQLCliKodo) GetSeries(resp map[string]any) []any {
	var datas []any
	if v, ok := resp["content"]; ok {
		if v, ok := v.([]any); ok {
			datas = v
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
