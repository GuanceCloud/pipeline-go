// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2026-present Guance, Inc.

package pipelinecheck

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/GuanceCloud/cliutils/point"
	"github.com/GuanceCloud/pipeline-go/lang"
	"github.com/GuanceCloud/pipeline-go/lang/platypus"
	"github.com/GuanceCloud/pipeline-go/ptinput"
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

type stringFlag struct {
	value string
	set   bool
}

func (v *stringFlag) String() string {
	return v.value
}

func (v *stringFlag) Set(s string) error {
	v.value = s
	v.set = true
	return nil
}

type scriptResult struct {
	Name    string `json:"name"`
	Parsed  bool   `json:"parsed"`
	Ran     bool   `json:"ran,omitempty"`
	Content string `json:"content,omitempty"`
}

type inputResult struct {
	Category        string            `json:"category"`
	Name            string            `json:"name"`
	Tags            map[string]string `json:"tags,omitempty"`
	Fields          map[string]any    `json:"fields"`
	MessageIsJSON   bool              `json:"message_is_json,omitempty"`
	MessageJSONType string            `json:"message_json_type,omitempty"`
	MessageJSONKeys []string          `json:"message_json_keys,omitempty"`
}

type pointResult struct {
	Category string            `json:"category"`
	Name     string            `json:"name"`
	Tags     map[string]string `json:"tags,omitempty"`
	Fields   map[string]any    `json:"fields,omitempty"`
	Dropped  bool              `json:"dropped"`
	Time     string            `json:"time"`
	TimeUnix int64             `json:"time_unix_nano"`
	Sub      []pointResult     `json:"sub_points,omitempty"`
}

type checkResult struct {
	OK       bool           `json:"ok"`
	Script   scriptResult   `json:"script"`
	Input    *inputResult   `json:"input,omitempty"`
	Output   *pointResult   `json:"output,omitempty"`
	Result   *extractResult `json:"result,omitempty"`
	Errors   []string       `json:"errors,omitempty"`
	Required []string       `json:"required_keys,omitempty"`
}

type extractResult struct {
	ExtractedFields map[string]any    `json:"extracted_fields,omitempty"`
	ExtractedTags   map[string]string `json:"extracted_tags,omitempty"`
	Dropped         bool              `json:"dropped"`
	Time            string            `json:"time"`
	TimeUnix        int64             `json:"time_unix_nano"`
	FieldCount      int               `json:"field_count"`
	TagCount        int               `json:"tag_count"`
	SubPointCount   int               `json:"sub_point_count,omitempty"`
}

type config struct {
	scriptPath   string
	scriptCmd    string
	message      stringFlag
	messageFile  string
	category     string
	name         string
	showScript   bool
	addStatus    bool
	tags         repeatedString
	fields       repeatedString
	requireKeys  repeatedString
	expectValues repeatedString

	listFunctions   bool
	searchFunctions string
	functionDoc     string
	functionLang    string
	functionLimit   int
}

func Run(args []string, stdout, stderr io.Writer) int {
	return run(args, stdout, stderr)
}

func run(args []string, stdout, stderr io.Writer) int {
	cfg := config{
		category: "logging",
		name:     "pipeline_check",
	}

	fs := flag.NewFlagSet("pipeline-check", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&cfg.scriptPath, "script", "", "pipeline script file")
	fs.StringVar(&cfg.scriptPath, "P", "", "pipeline script file")
	fs.StringVar(&cfg.scriptCmd, "cmd", "", "pipeline script passed as a string")
	fs.StringVar(&cfg.scriptCmd, "c", "", "pipeline script passed as a string")
	fs.Var(&cfg.message, "message", "message field value used as pipeline input")
	fs.Var(&cfg.message, "T", "message field value used as pipeline input")
	fs.StringVar(&cfg.messageFile, "message-file", "", "file containing the message field value")
	fs.StringVar(&cfg.category, "category", cfg.category, "point category, such as logging, metric, tracing, rum, security")
	fs.StringVar(&cfg.name, "name", cfg.name, "point name")
	fs.BoolVar(&cfg.showScript, "show-script", false, "include script content in JSON output")
	fs.BoolVar(&cfg.addStatus, "add-status", false, "enable logging status normalization while running")
	fs.Var(&cfg.tags, "tag", "extra input tag in key=value form; repeatable")
	fs.Var(&cfg.fields, "field", "extra input field in key=value form; JSON values are decoded when possible; repeatable")
	fs.Var(&cfg.requireKeys, "require-key", "require output tag or field to exist after the script runs; repeatable")
	fs.Var(&cfg.expectValues, "expect", "require output tag or field to equal key=value; JSON values are decoded when possible; repeatable")
	fs.BoolVar(&cfg.listFunctions, "list-functions", false, "list embedded pipeline function docs")
	fs.BoolVar(&cfg.listFunctions, "fn-list", false, "list embedded pipeline function docs")
	fs.StringVar(&cfg.searchFunctions, "search-functions", "", "search embedded pipeline function docs")
	fs.StringVar(&cfg.searchFunctions, "fn-search", "", "search embedded pipeline function docs")
	fs.StringVar(&cfg.functionDoc, "function-doc", "", "show one embedded pipeline function doc")
	fs.StringVar(&cfg.functionDoc, "fn-doc", "", "show one embedded pipeline function doc")
	fs.StringVar(&cfg.functionLang, "function-lang", "zh", "function doc language: zh, en, or all")
	fs.IntVar(&cfg.functionLimit, "function-limit", 0, "maximum function docs returned by list/search; 0 means no limit")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	positional := fs.Args()
	if cfg.scriptPath == "" && cfg.scriptCmd == "" && len(positional) > 0 {
		cfg.scriptPath = positional[0]
		positional = positional[1:]
	}
	if !cfg.message.set && cfg.messageFile == "" && len(positional) > 0 {
		cfg.message.value = strings.Join(positional, " ")
		cfg.message.set = true
	}

	if cfg.hasFunctionDocMode() {
		res, exitCode := executeFunctionDocs(cfg)
		enc := json.NewEncoder(stdout)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		if err := enc.Encode(res); err != nil {
			fmt.Fprintf(stderr, "encode result: %v\n", err)
			return 1
		}
		return exitCode
	}

	res := checkResult{}
	exitCode := execute(cfg, &res)

	enc := json.NewEncoder(stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(res); err != nil {
		fmt.Fprintf(stderr, "encode result: %v\n", err)
		return 1
	}

	return exitCode
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

	cat, err := parseCategory(cfg.category)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return 1
	}

	scripts, parseErrs := platypus.NewScripts(map[string]string{name: script}, lang.WithCat(cat))
	if err := parseErrs[name]; err != nil {
		res.Errors = append(res.Errors, err.Error())
		return 1
	}

	pipelineScript := scripts[name]
	if pipelineScript == nil {
		res.Errors = append(res.Errors, "script parser did not return a runnable script")
		return 1
	}
	defer pipelineScript.Cleanup()
	res.Script.Parsed = true

	message, hasMessage, err := loadMessage(cfg.message, cfg.messageFile)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return 1
	}

	tags, err := parseTags(cfg.tags)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return 1
	}
	fields, err := parseFields(cfg.fields)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return 1
	}

	if !hasMessage {
		res.OK = true
		return 0
	}

	fields[ptinput.Originkey] = message
	res.Input = newInputResult(cat, cfg.name, tags, fields, message)

	plpt := ptinput.NewPlPoint(cat, cfg.name, tags, fields, time.Now())
	if err := pipelineScript.Run(plpt, nil, &lang.LogOption{
		DisableAddStatusField: !cfg.addStatus,
	}); err != nil {
		res.Output = snapshotPoint(plpt)
		res.Result = summarizeResult(res.Input, res.Output)
		res.Errors = append(res.Errors, err.Error())
		return 1
	}
	res.Script.Ran = true
	res.Output = snapshotPoint(plpt)
	res.Result = summarizeResult(res.Input, res.Output)

	for _, key := range cfg.requireKeys {
		if _, ok := outputValue(res.Output, key); !ok {
			res.Errors = append(res.Errors, fmt.Sprintf("required key %q not found in output tags or fields", key))
			res.Required = append(res.Required, key)
		}
	}

	for _, expr := range cfg.expectValues {
		key, expectedRaw, err := parseKeyValue(expr)
		if err != nil {
			res.Errors = append(res.Errors, fmt.Sprintf("invalid expectation %q: %v", expr, err))
			continue
		}
		got, ok := outputValue(res.Output, key)
		if !ok {
			res.Errors = append(res.Errors, fmt.Sprintf("expected key %q not found in output tags or fields", key))
			continue
		}
		expected := parseLooseValue(expectedRaw)
		if !valuesEqual(got, expected) {
			res.Errors = append(res.Errors, fmt.Sprintf("expected %s=%v, got %v", key, expected, got))
		}
	}

	res.OK = len(res.Errors) == 0
	if !res.OK {
		return 1
	}
	return 0
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
		return "", "", fmt.Errorf("missing pipeline script: pass --script FILE or --cmd TEXT")
	}
}

func loadMessage(message stringFlag, file string) (string, bool, error) {
	if message.set && file != "" {
		return "", false, fmt.Errorf("use only one of --message or --message-file")
	}
	if file != "" {
		b, err := os.ReadFile(file)
		if err != nil {
			return "", false, err
		}
		return string(b), true, nil
	}
	return message.value, message.set, nil
}

func parseCategory(s string) (point.Category, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "log", "logs":
		return point.Logging, nil
	}

	if cat := point.CatString(strings.ToLower(s)); cat != point.UnknownCategory {
		return cat, nil
	}
	if cat := point.CatAlias(strings.ToUpper(s)); cat != point.UnknownCategory {
		return cat, nil
	}

	return point.UnknownCategory, fmt.Errorf("unknown category %q", s)
}

func parseTags(raw []string) (map[string]string, error) {
	tags := map[string]string{}
	for _, item := range raw {
		k, v, err := parseKeyValue(item)
		if err != nil {
			return nil, err
		}
		tags[k] = v
	}
	return tags, nil
}

func parseFields(raw []string) (map[string]any, error) {
	fields := map[string]any{}
	for _, item := range raw {
		k, v, err := parseKeyValue(item)
		if err != nil {
			return nil, err
		}
		fields[k] = parseLooseValue(v)
	}
	return fields, nil
}

func parseKeyValue(s string) (string, string, error) {
	k, v, ok := strings.Cut(s, "=")
	if !ok {
		return "", "", fmt.Errorf("expected key=value")
	}
	k = strings.TrimSpace(k)
	if k == "" {
		return "", "", fmt.Errorf("key is empty")
	}
	return k, v, nil
}

func parseLooseValue(s string) any {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return s
	}

	dec := json.NewDecoder(strings.NewReader(trimmed))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err == nil && dec.Decode(&struct{}{}) == io.EOF {
		return normalizeJSONNumber(v)
	}
	return s
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

func newInputResult(cat point.Category, name string, tags map[string]string, fields map[string]any, message string) *inputResult {
	valid, kind, keys := inspectJSON(message)
	return &inputResult{
		Category:        cat.String(),
		Name:            name,
		Tags:            cloneStringMap(tags),
		Fields:          cloneAnyMap(fields),
		MessageIsJSON:   valid,
		MessageJSONType: kind,
		MessageJSONKeys: keys,
	}
}

func inspectJSON(s string) (bool, string, []string) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return false, "", nil
	}

	dec := json.NewDecoder(strings.NewReader(trimmed))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return false, "", nil
	}
	if dec.Decode(&struct{}{}) != io.EOF {
		return false, "", nil
	}

	switch x := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return true, "object", keys
	case []any:
		return true, "array", nil
	case string:
		return true, "string", nil
	case bool:
		return true, "bool", nil
	case nil:
		return true, "null", nil
	default:
		return true, "number", nil
	}
}

func snapshotPoint(plpt ptinput.PlInputPt) *pointResult {
	if plpt == nil {
		return nil
	}

	tm := plpt.PtTime()
	res := &pointResult{
		Category: plpt.Category().String(),
		Name:     plpt.GetPtName(),
		Tags:     cloneStringMap(plpt.Tags()),
		Fields:   cloneAnyMap(plpt.Fields()),
		Dropped:  plpt.Dropped(),
		Time:     tm.Format(time.RFC3339Nano),
		TimeUnix: tm.UnixNano(),
	}

	for _, sub := range plpt.GetSubPoint() {
		if p := snapshotPoint(sub); p != nil {
			res.Sub = append(res.Sub, *p)
		}
	}

	return res
}

func summarizeResult(input *inputResult, output *pointResult) *extractResult {
	if output == nil {
		return nil
	}

	inputFields := map[string]any{}
	message := ""
	if input != nil {
		for k, v := range input.Fields {
			inputFields[k] = v
		}
		if v, ok := input.Fields[ptinput.Originkey]; ok {
			message = fmt.Sprint(v)
		}
	}

	fields := map[string]any{}
	if message != "" {
		fields[ptinput.Originkey] = message
	}
	for k, v := range output.Fields {
		if old, existed := inputFields[k]; existed && valuesEqual(old, v) {
			continue
		}
		fields[k] = v
	}

	inputTags := map[string]string{}
	if input != nil {
		for k, v := range input.Tags {
			inputTags[k] = v
		}
	}

	tags := map[string]string{}
	for k, v := range output.Tags {
		if old, existed := inputTags[k]; existed && old == v {
			continue
		}
		tags[k] = v
	}

	return &extractResult{
		ExtractedFields: cloneAnyMap(fields),
		ExtractedTags:   cloneStringMap(tags),
		Dropped:         output.Dropped,
		Time:            output.Time,
		TimeUnix:        output.TimeUnix,
		FieldCount:      len(fields),
		TagCount:        len(tags),
		SubPointCount:   len(output.Sub),
	}
}

func outputValue(out *pointResult, key string) (any, bool) {
	if out == nil {
		return nil, false
	}
	if v, ok := out.Fields[key]; ok {
		return v, true
	}
	if v, ok := out.Tags[key]; ok {
		return v, true
	}
	return nil, false
}

func valuesEqual(got, expected any) bool {
	if reflect.DeepEqual(got, expected) {
		return true
	}
	return fmt.Sprint(got) == fmt.Sprint(expected)
}

func cloneStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
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

func jsonSafeValue(v any) any {
	switch x := v.(type) {
	case map[string]any:
		return cloneAnyMap(x)
	case []any:
		out := make([]any, len(x))
		for i := range x {
			out[i] = jsonSafeValue(x[i])
		}
		return out
	default:
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(x); err == nil {
			return x
		}
		return fmt.Sprint(x)
	}
}
