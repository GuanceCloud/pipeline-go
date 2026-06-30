// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2026-present Guance, Inc.

package arbitercheck

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/GuanceCloud/pipeline-go/pkg/arbiter"
	funcs "github.com/GuanceCloud/pipeline-go/pkg/arbiter/builtin-funcs"
	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/dql"
	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/request"
	arbscript "github.com/GuanceCloud/pipeline-go/pkg/arbiter/script"
	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/trigger"
	"github.com/GuanceCloud/platypus/pkg/token"
)

type repeatedString []string

func (v *repeatedString) String() string {
	if v == nil {
		return ""
	}
	return strings.Join(*v, ",")
}

func (v *repeatedString) Set(s string) error {
	*v = append(*v, s)
	return nil
}

type scriptResult struct {
	Name    string `json:"name"`
	Parsed  bool   `json:"parsed"`
	Checked bool   `json:"checked,omitempty"`
	Ran     bool   `json:"ran,omitempty"`
	Content string `json:"content,omitempty"`
}

type checkResult struct {
	OK            bool           `json:"ok"`
	Script        scriptResult   `json:"script"`
	DQLMode       string         `json:"dql_mode,omitempty"`
	Stdout        string         `json:"stdout,omitempty"`
	Triggers      []trigger.Data `json:"triggers,omitempty"`
	DQLQueries    []dqlCall      `json:"dql_queries,omitempty"`
	DQLChecks     []dqlCheck     `json:"dql_checks,omitempty"`
	MockDQLResult map[string]any `json:"mock_dql_result,omitempty"`
	Errors        []string       `json:"errors,omitempty"`
}

type dqlCall struct {
	Query           string   `json:"query"`
	QType           string   `json:"qtype"`
	Limit           int64    `json:"limit"`
	Offset          int64    `json:"offset"`
	SLimit          int64    `json:"slimit"`
	AlignTime       bool     `json:"align_time"`
	DisableSampling bool     `json:"disable_sampling"`
	TimeRange       []any    `json:"time_range,omitempty"`
	WorkspaceUUID   []string `json:"workspace_uuid,omitempty"`
}

type dqlCheck struct {
	Query  string         `json:"query"`
	OK     bool           `json:"ok"`
	Result map[string]any `json:"result,omitempty"`
	Stdout string         `json:"stdout,omitempty"`
	Stderr string         `json:"stderr,omitempty"`
	Error  string         `json:"error,omitempty"`
}

type config struct {
	scriptPath    string
	scriptCmd     string
	parseOnly     bool
	showScript    bool
	dqlResult     string
	dqlResultFile string
	timeRange     string
	duration      string
	liveDQL       bool
	guance        string
	guanceKey     string
	checkDQL      bool
	dqlCheckBin   string

	requireTrigger bool
	expectStatuses repeatedString

	listFunctions   bool
	searchFunctions string
	functionDoc     string
	functionLang    string
	functionLimit   int
}

type mockDQL struct {
	result    map[string]any
	timeRange []int64
	calls     []dqlCall
}

type recordingDQL struct {
	inner dql.DQL
	calls []dqlCall
}

func Run(args []string, stdout, stderr io.Writer) int {
	return run(args, stdout, stderr)
}

func run(args []string, stdout, stderr io.Writer) int {
	cfg := config{
		guance:    envDefault("GUANCE_OPENAPI_ENDPOINT", "https://openapi.guance.com"),
		guanceKey: envDefault("GUANCE_API_KEY", envDefault("DF_API_KEY", "")),
		duration:  "15m",
	}

	fs := flag.NewFlagSet("arbiter-check", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&cfg.scriptPath, "script", "", "arbiter script file")
	fs.StringVar(&cfg.scriptPath, "P", "", "arbiter script file")
	fs.StringVar(&cfg.scriptCmd, "cmd", "", "arbiter script passed as a string")
	fs.StringVar(&cfg.scriptCmd, "c", "", "arbiter script passed as a string")
	fs.BoolVar(&cfg.parseOnly, "parse-only", false, "only parse and type-check the script")
	fs.BoolVar(&cfg.showScript, "show-script", false, "include script content in JSON output")
	fs.StringVar(&cfg.dqlResult, "dql-result", "", "mock DQL JSON result used by dql()")
	fs.StringVar(&cfg.dqlResultFile, "dql-result-file", "", "file containing mock DQL JSON result used by dql()")
	fs.StringVar(&cfg.timeRange, "time-range", "", "DQL time range in millisecond epoch format: start,end")
	fs.StringVar(&cfg.duration, "duration", cfg.duration, "live DQL query time range duration, such as 15m, 1h, 24h; ignored when --time-range is set")
	fs.BoolVar(&cfg.liveDQL, "live-dql", false, "run dql() against GuanceCloud OpenAPI instead of the mock DQL result")
	fs.StringVar(&cfg.guance, "guance", cfg.guance, "GuanceCloud OpenAPI endpoint for --live-dql")
	fs.StringVar(&cfg.guanceKey, "guance-key", cfg.guanceKey, "GuanceCloud OpenAPI key for --live-dql; defaults to GUANCE_API_KEY or DF_API_KEY")
	fs.BoolVar(&cfg.checkDQL, "check-dql", false, "validate each dql() query with dqlcheck")
	fs.StringVar(&cfg.dqlCheckBin, "dqlcheck-bin", "", "path to dqlcheck; defaults to bundled skills/dql/bin/dqlcheck or PATH")
	fs.BoolVar(&cfg.requireTrigger, "require-trigger", false, "fail if the script does not call trigger()")
	fs.Var(&cfg.expectStatuses, "expect-status", "require at least one trigger with this status; repeatable")
	fs.BoolVar(&cfg.listFunctions, "list-functions", false, "list embedded arbiter function docs")
	fs.BoolVar(&cfg.listFunctions, "fn-list", false, "list embedded arbiter function docs")
	fs.StringVar(&cfg.searchFunctions, "search-functions", "", "search embedded arbiter function docs")
	fs.StringVar(&cfg.searchFunctions, "fn-search", "", "search embedded arbiter function docs")
	fs.StringVar(&cfg.functionDoc, "function-doc", "", "show one embedded arbiter function doc")
	fs.StringVar(&cfg.functionDoc, "fn-doc", "", "show one embedded arbiter function doc")
	fs.StringVar(&cfg.functionLang, "function-lang", "zh", "function doc language: zh, en, or all")
	fs.IntVar(&cfg.functionLimit, "function-limit", 0, "maximum function docs returned by list/search; 0 means no limit")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	positional := fs.Args()
	if cfg.scriptPath == "" && cfg.scriptCmd == "" && len(positional) > 0 {
		cfg.scriptPath = positional[0]
	}

	if cfg.hasFunctionDocMode() {
		res, exitCode := executeFunctionDocs(cfg)
		writeJSON(stdout, stderr, res)
		return exitCode
	}

	res := checkResult{}
	exitCode := execute(cfg, &res)
	writeJSON(stdout, stderr, res)
	return exitCode
}

func writeJSON(stdout, stderr io.Writer, v any) {
	enc := json.NewEncoder(stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintf(stderr, "encode result: %v\n", err)
	}
}

func (cfg config) hasFunctionDocMode() bool {
	return cfg.listFunctions || cfg.searchFunctions != "" || cfg.functionDoc != ""
}

func execute(cfg config, res *checkResult) int {
	name, script, err := loadScript(cfg.scriptPath, cfg.scriptCmd)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return 1
	}
	res.Script.Name = name
	if cfg.showScript {
		res.Script.Content = script
	}

	if _, err := arbscript.Parse(name, script, funcs.Funcs); err != nil {
		res.Errors = append(res.Errors, err.Error())
		return 1
	}
	res.Script.Parsed = true
	res.Script.Checked = true

	if cfg.parseOnly {
		res.OK = true
		return 0
	}

	timeRange, err := effectiveTimeRange(cfg)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return 1
	}

	stdout := bytes.NewBuffer(nil)
	tr := trigger.NewTr()

	dqlClient, err := newDQLClient(cfg, timeRange, res)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return 1
	}

	if err := arbiter.Run(name, script,
		arbiter.WithFuncs(funcs.Funcs),
		arbiter.WithStdout(stdout),
		arbiter.WithTrigger(tr),
		arbiter.WithDQLClient(dqlClient),
		arbiter.WithHTTPClient(request.NewHTTPClient(0)),
	); err != nil {
		res.Stdout = stdout.String()
		res.Triggers = tr.Result()
		res.DQLQueries = dqlClient.Calls()
		appendDQLChecks(cfg, res)
		res.Errors = append(res.Errors, err.Error())
		return 1
	}

	res.Script.Ran = true
	res.Stdout = stdout.String()
	res.Triggers = tr.Result()
	res.DQLQueries = dqlClient.Calls()
	appendDQLChecks(cfg, res)

	if cfg.requireTrigger && len(res.Triggers) == 0 {
		res.Errors = append(res.Errors, "required trigger output not found")
	}
	for _, status := range cfg.expectStatuses {
		if !hasTriggerStatus(res.Triggers, status) {
			res.Errors = append(res.Errors, fmt.Sprintf("expected trigger status %q not found", status))
		}
	}

	res.OK = len(res.Errors) == 0
	if !res.OK {
		return 1
	}
	return 0
}

func appendDQLChecks(cfg config, res *checkResult) {
	if !cfg.checkDQL {
		return
	}
	checks, errs := runDQLChecks(cfg, res.DQLQueries)
	res.DQLChecks = checks
	res.Errors = append(res.Errors, errs...)
}

func runDQLChecks(cfg config, queries []dqlCall) ([]dqlCheck, []string) {
	bin := resolveDQLCheckBin(cfg.dqlCheckBin)
	var checks []dqlCheck
	var errs []string
	for i, q := range queries {
		check := runDQLCheck(bin, q.Query)
		checks = append(checks, check)
		if !check.OK {
			errs = append(errs, fmt.Sprintf("dql query %d failed dqlcheck", i+1))
		}
	}
	return checks, errs
}

func resolveDQLCheckBin(configured string) string {
	if configured != "" {
		return configured
	}

	candidates := []string{
		filepath.Join("skills", "dql", "bin", dqlCheckName()),
		filepath.Join("..", "dql", "bin", dqlCheckName()),
	}

	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(exeDir, "..", "..", "dql", "bin", dqlCheckName()),
		)
	}

	for _, local := range candidates {
		if st, err := os.Stat(local); err == nil && !st.IsDir() {
			return local
		}
	}
	return dqlCheckName()
}

func dqlCheckName() string {
	if runtime.GOOS == "windows" {
		return "dqlcheck.exe"
	}
	return "dqlcheck"
}

func runDQLCheck(bin, query string) dqlCheck {
	check := dqlCheck{
		Query: query,
	}
	cmd := exec.Command(bin, "--stdin", "--format=json", "--pretty")
	cmd.Stdin = strings.NewReader(query)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	check.Stdout = strings.TrimSpace(stdout.String())
	check.Stderr = strings.TrimSpace(stderr.String())
	if check.Stdout != "" {
		var result map[string]any
		dec := json.NewDecoder(strings.NewReader(check.Stdout))
		dec.UseNumber()
		if err := dec.Decode(&result); err == nil {
			check.Result = cloneAnyMap(normalizeJSONNumber(result).(map[string]any))
			if ok, _ := check.Result["ok"].(bool); ok {
				check.OK = true
			}
		}
	}
	if err != nil {
		check.Error = err.Error()
	}
	return check
}

func envDefault(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

type recordedDQL interface {
	dql.DQL
	Calls() []dqlCall
}

func newDQLClient(cfg config, timeRange []int64, res *checkResult) (recordedDQL, error) {
	if cfg.liveDQL {
		if cfg.guanceKey == "" {
			return nil, fmt.Errorf("--live-dql requires --guance-key or GUANCE_API_KEY/DF_API_KEY")
		}
		if cfg.dqlResult != "" || cfg.dqlResultFile != "" {
			return nil, fmt.Errorf("do not use --dql-result or --dql-result-file with --live-dql")
		}
		res.DQLMode = "live"
		return &recordingDQL{
			inner: dql.NewDQLOpenAPI(cfg.guance, dql.OpenAPIPath, cfg.guanceKey, timeRange),
		}, nil
	}

	mockResult, err := loadDQLResult(cfg.dqlResult, cfg.dqlResultFile)
	if err != nil {
		return nil, err
	}
	res.DQLMode = "mock"
	res.MockDQLResult = cloneAnyMap(mockResult)
	return &mockDQL{
		result:    mockResult,
		timeRange: timeRange,
	}, nil
}

func loadScript(path, inline string) (string, string, error) {
	switch {
	case path != "" && inline != "":
		return "", "", fmt.Errorf("use only one of --script or --cmd")
	case inline != "":
		return "inline.p", inline, nil
	case path != "":
		b, err := os.ReadFile(path)
		if err != nil {
			return "", "", err
		}
		return filepath.Base(path), string(b), nil
	default:
		return "", "", fmt.Errorf("missing arbiter script: pass --script FILE or --cmd TEXT")
	}
}

func loadDQLResult(raw, file string) (map[string]any, error) {
	if raw != "" && file != "" {
		return nil, fmt.Errorf("use only one of --dql-result or --dql-result-file")
	}
	if file != "" {
		b, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		raw = string(b)
	}
	if strings.TrimSpace(raw) == "" {
		return map[string]any{
			"series":      []any{},
			"status_code": int64(200),
		}, nil
	}

	var v any
	dec := json.NewDecoder(strings.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	if dec.Decode(&struct{}{}) != io.EOF {
		return nil, fmt.Errorf("mock DQL result must contain a single JSON value")
	}
	m, ok := normalizeJSONNumber(v).(map[string]any)
	if !ok {
		return nil, fmt.Errorf("mock DQL result must be a JSON object")
	}
	return m, nil
}

func parseTimeRange(s string) ([]int64, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("--time-range expects start,end")
	}
	start, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid time range start: %w", err)
	}
	end, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid time range end: %w", err)
	}
	return []int64{start, end}, nil
}

func effectiveTimeRange(cfg config) ([]int64, error) {
	if strings.TrimSpace(cfg.timeRange) != "" {
		return parseTimeRange(cfg.timeRange)
	}
	if !cfg.liveDQL {
		return nil, nil
	}
	d, err := time.ParseDuration(cfg.duration)
	if err != nil {
		return nil, err
	}
	if d <= 0 {
		return nil, fmt.Errorf("--duration must be greater than zero")
	}
	end := time.Now().UnixMilli()
	return []int64{end - d.Milliseconds(), end}, nil
}

func (m *mockDQL) Query(pos token.LnColPos, q, qTyp string, limit, offset, slimit int64, timeRange []any, alignTime, disableSampling bool, uuids ...string) (map[string]any, error) {
	m.calls = append(m.calls, dqlCall{
		Query:           q,
		QType:           qTyp,
		Limit:           limit,
		Offset:          offset,
		SLimit:          slimit,
		AlignTime:       alignTime,
		DisableSampling: disableSampling,
		TimeRange:       cloneAnySlice(timeRange),
		WorkspaceUUID:   append([]string{}, uuids...),
	})
	return cloneAnyMap(m.result), nil
}

func (m *mockDQL) TimeRange() []int64 {
	if len(m.timeRange) != 2 {
		return nil
	}
	return append([]int64{}, m.timeRange...)
}

func (m *mockDQL) Calls() []dqlCall {
	return append([]dqlCall{}, m.calls...)
}

var _ dql.DQL = (*mockDQL)(nil)
var _ recordedDQL = (*mockDQL)(nil)

func (r *recordingDQL) Query(pos token.LnColPos, q, qTyp string, limit, offset, slimit int64, timeRange []any, alignTime, disableSampling bool, uuids ...string) (map[string]any, error) {
	r.calls = append(r.calls, dqlCall{
		Query:           q,
		QType:           qTyp,
		Limit:           limit,
		Offset:          offset,
		SLimit:          slimit,
		AlignTime:       alignTime,
		DisableSampling: disableSampling,
		TimeRange:       cloneAnySlice(timeRange),
		WorkspaceUUID:   append([]string{}, uuids...),
	})
	return r.inner.Query(pos, q, qTyp, limit, offset, slimit, timeRange, alignTime, disableSampling, uuids...)
}

func (r *recordingDQL) TimeRange() []int64 {
	return r.inner.TimeRange()
}

func (r *recordingDQL) Calls() []dqlCall {
	return append([]dqlCall{}, r.calls...)
}

var _ recordedDQL = (*recordingDQL)(nil)

func hasTriggerStatus(triggers []trigger.Data, status string) bool {
	for _, tr := range triggers {
		if tr.Status == status {
			return true
		}
	}
	return false
}

func normalizeJSONNumber(v any) any {
	switch x := v.(type) {
	case json.Number:
		if i, err := x.Int64(); err == nil {
			return i
		}
		if f, err := x.Float64(); err == nil {
			return f
		}
		return x.String()
	case []any:
		for i := range x {
			x[i] = normalizeJSONNumber(x[i])
		}
		return x
	case map[string]any:
		for k := range x {
			x[k] = normalizeJSONNumber(x[k])
		}
		return x
	default:
		return v
	}
}

func cloneAnyMap(in map[string]any) map[string]any {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = jsonSafeValue(v)
	}
	return out
}

func cloneAnySlice(in []any) []any {
	if len(in) == 0 {
		return []any{}
	}
	out := make([]any, len(in))
	for i, v := range in {
		out[i] = jsonSafeValue(v)
	}
	return out
}

func jsonSafeValue(v any) any {
	switch x := v.(type) {
	case map[string]any:
		return cloneAnyMap(x)
	case []any:
		return cloneAnySlice(x)
	default:
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(x); err == nil {
			return x
		}
		return fmt.Sprint(x)
	}
}

func valuesEqual(got, expected any) bool {
	if reflect.DeepEqual(got, expected) {
		return true
	}
	return fmt.Sprint(got) == fmt.Sprint(expected)
}
