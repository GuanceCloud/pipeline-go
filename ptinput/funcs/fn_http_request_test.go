package funcs

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"github.com/GuanceCloud/cliutils/point"
	"github.com/GuanceCloud/pipeline-go/ptinput"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	cases := []struct {
		url             string
		filterResult    bool
		disableInternal bool
		cidrs, hosts    []string
	}{
		{
			url:             "http://0.0.0.0/",
			filterResult:    true,
			disableInternal: true,
			cidrs:           nil,
		},
		{
			url:             "http://localhost/",
			filterResult:    true,
			disableInternal: true,
			cidrs:           nil,
		},
		{
			url:             "http://127.0.1.1/",
			filterResult:    true,
			disableInternal: true,
			cidrs:           nil,
		},
		{
			url:             "http://1.0.0.1/",
			filterResult:    false,
			disableInternal: true,
			cidrs:           nil,
		},
		{
			url:             "http://1.0.0.1:1234/",
			filterResult:    false,
			disableInternal: true,
			cidrs:           nil,
		},
		{
			url:             "http://[::]:1234/",
			filterResult:    false,
			disableInternal: false,
			cidrs:           nil,
		},
		{
			url:             "http://[::]:1234/",
			filterResult:    true,
			disableInternal: true,
			cidrs:           nil,
		},
		{
			url:             "http://[::]/",
			filterResult:    true,
			disableInternal: true,
			cidrs:           nil,
		},
		{
			url:             "http://1.0.0.1",
			filterResult:    false,
			disableInternal: true,
			cidrs:           []string{"1.0.0.0/16"},
		},
		{
			url:             "http://10.0.0.1",
			filterResult:    true,
			disableInternal: true,
			cidrs:           nil,
		},
		{
			url:             "http://192.168.0.1",
			filterResult:    true,
			disableInternal: true,
			cidrs:           nil,
		},
		{
			url:             "http://10.0.0.1",
			filterResult:    false,
			disableInternal: false,
			cidrs:           nil,
		},
		{
			url:             "http://10.0.0.1",
			filterResult:    false,
			disableInternal: false,
			cidrs:           []string{"10.0.0.1/16"},
		},
		{
			url:             "file://ccc/",
			filterResult:    true,
			disableInternal: true,
			cidrs:           nil,
		},
		{
			url:             "https://guance.com",
			filterResult:    false,
			disableInternal: true,
			cidrs:           nil,
		},
		{
			url:          "https://guance.com",
			hosts:        []string{"guance.com"},
			filterResult: false,
		},
		{
			url:          "https://guance.com",
			hosts:        []string{"guancez.com"},
			filterResult: true,
		},
		{
			url:             "https://127.0.0.1",
			hosts:           []string{"127.0.0.1"},
			filterResult:    false,
			disableInternal: true,
		},
		{
			url:             "https://127.0.0.1",
			cidrs:           []string{"127.0.0.1/32"},
			filterResult:    false,
			disableInternal: true,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			r := filterURL(
				c.url, c.disableInternal, c.cidrs, c.hosts,
			)
			if r != c.filterResult {
				assert.Equal(t, c.filterResult, r)
			}
		})
	}
}

func TestBuildBody(t *testing.T) {
	cases := []struct {
		val    any
		result string
	}{
		{
			val:    float64(123.1),
			result: "123.1",
		},
		{
			val:    int64(123),
			result: "123",
		},
		{
			val:    true,
			result: "true",
		},
		{
			val:    false,
			result: "false",
		},
		{
			val:    "abc",
			result: "abc",
		},
		{
			val:    []any{1, 2, 3},
			result: "[1,2,3]",
		},
		{
			val:    map[string]any{"a": 1, "b": 2},
			result: `{"a":1,"b":2}`,
		},
		{
			val:    nil,
			result: "",
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("index_%d", i), func(t *testing.T) {
			var buf []byte
			if b := buildBody(c.val); b != nil {
				var err error
				buf, err = io.ReadAll(b)
				if err != nil && !errors.Is(err, io.EOF) {
					t.Error(err)
				}
			}
			assert.Equal(t, c.result, string(buf))
		})
	}
}

func TestHTTPRequest(t *testing.T) {
	server := HTTPServer()
	defer server.Close()

	url := `"` + server.URL + "/testResp" + `"`
	fmt.Println(url)

	cases := []struct {
		name, pl, in string
		expected     interface{}
		absent       []string
		fail         bool
		outkey       string
	}{
		{
			name: "test_post",
			pl: fmt.Sprintf(`
			resp = http_request("POST", %s, {"extraHeader": "1",
			"extraHeader": "1"}, {"a": "1"})
			add_key(abc, resp["body"])
			`, url),
			in:       `[]`,
			outkey:   "abc",
			expected: `{"a":"1"}`,
		},
		{
			name: "test_file",
			pl: `
			resp = http_request("POST", "file:///etc/", {"extraHeader": "1", 
			"extraHeader": "1"}, {"a": "1"})
			add_key(abc, resp)
			`,
			in:       `[]`,
			outkey:   "abc",
			expected: nil,
		},
		{
			name: "test_put",
			pl: fmt.Sprintf(`
			resp = http_request("put", %s, {"extraHeader": "1",
			"extraHeader": "1"}, {"a": "1"})
			add_key(abc, resp["body"])
			`, url),
			in:       `[]`,
			outkey:   "abc",
			expected: `{"a":"1"}`,
		},
		{
			name: "legacy_required_args_only",
			pl: fmt.Sprintf(`
			resp = http_request("GET", %s)
			add_key(status_code, resp["status_code"])
			`, url),
			in:       `[]`,
			outkey:   "status_code",
			expected: int64(200),
		},
		{
			name: "test_prefix",
			pl: fmt.Sprintf(`
			resp = http_request("POST", %s, {"extraHeader": "1",
			"extraHeader": "1"}, {"a": "1"}, "hr_")
			add_key(abc, resp["hr_body"])
			if resp["body"] != nil {
				add_key(unexpected_unprefixed_key, true)
			}
			`, url),
			in:       `[]`,
			outkey:   "abc",
			expected: `{"a":"1"}`,
			absent:   []string{"unexpected_unprefixed_key"},
		},
		{
			name: "test_prefix_status_code",
			pl: fmt.Sprintf(`
			resp = http_request("GET", %s, {"extraHeader": "1",
			"extraHeader": "1"}, nil, "hr_")
			add_key(abc, resp["hr_status_code"])
			`, url),
			in:       `[]`,
			outkey:   "abc",
			expected: int64(200),
		},
		{
			name: "test_prefix_variable",
			pl: fmt.Sprintf(`
			p = "hr_"
			resp = http_request("POST", %s, {"extraHeader": "1",
			"extraHeader": "1"}, {"a": "1"}, p)
			add_key(abc, resp["hr_body"])
			`, url),
			in:       `[]`,
			outkey:   "abc",
			expected: `{"a":"1"}`,
		},
		{
			name: "test_named_prefix_skips_optional_args",
			pl: fmt.Sprintf(`
			resp = http_request("GET", %s, prefix="hr_")
			add_key(abc, resp["hr_body"])
			if resp["body"] != nil {
				add_key(unexpected_unprefixed_key, true)
			}
			`, url),
			in:       `[]`,
			outkey:   "abc",
			expected: `{"a":"hello"}`,
			absent:   []string{"unexpected_unprefixed_key"},
		},
		{
			name: "invalid prefix type",
			pl: fmt.Sprintf(`
			resp = http_request("POST", %s, {"extraHeader": "1",
			"extraHeader": "1"}, {"a": "1"}, 123)
			`, url),
			in:   `[]`,
			fail: true,
		},
		{
			name: "invalid prefix type variable",
			pl: fmt.Sprintf(`
			p = 123
			resp = http_request("POST", %s, {"extraHeader": "1",
			"extraHeader": "1"}, {"a": "1"}, p)
			`, url),
			in:   `[]`,
			fail: true,
		},
	}

	for idx, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runner, err := NewTestingRunner(tc.pl)
			if err != nil {
				if !tc.fail {
					t.Fatalf("[%d] unexpected compile error: %s", idx, err)
				}
				return
			}
			pt := ptinput.NewPlPt(
				point.Logging, "test", nil, map[string]any{"message": tc.in}, time.Now())
			errR := runScript(runner, pt)

			if errR != nil {
				if !tc.fail {
					t.Fatalf("[%d] unexpected runtime error: %s", idx, errR)
				}
				return
			}
			if tc.fail {
				t.Fatal("expected an error, got nil")
			}

			v, _, _ := pt.Get(tc.outkey)
			// tu.Equals(t, nil, err)
			assert.Equal(t, tc.expected, v)
			for _, k := range tc.absent {
				_, _, err := pt.Get(k)
				assert.Error(t, err, "key %q indicates an unprefixed result entry was generated", k)
			}

			t.Logf("[%d] PASS", idx)
		})
	}
}

func HTTPServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			headers := r.Header

			var respData []byte
			var err error
			if headers.Get("extraHeader1") != "" && headers.Get("extraHeader2") != "" {
				responseData := map[string]string{"a": "hello world"}
				respData, err = json.Marshal(responseData)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			} else {
				switch r.Method {
				case http.MethodGet:
					responseData := map[string]string{"a": "hello"}
					respData, err = json.Marshal(responseData)
					if err != nil {
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					}
				default:
					d, _ := io.ReadAll(r.Body)
					respData = d
				}
			}

			_, _ = w.Write(respData)
			w.WriteHeader(http.StatusOK)
		},
	))
	return server
}
