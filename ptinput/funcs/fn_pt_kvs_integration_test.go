package funcs

import (
	"testing"
	"time"

	"github.com/GuanceCloud/cliutils/point"
	"github.com/GuanceCloud/pipeline-go/ptinput"
	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/stretchr/testify/assert"
)

func TestPtKvsGetInOperator(t *testing.T) {
	pl := `
	arr = pt_kvs_get("nums", true)
	if 2 in arr {
		pt_kvs_set("hit", true)
	}
	`

	runner, err := NewTestingRunner(pl)
	if err != nil {
		t.Fatal(err)
	}

	raw := point.NewPoint("test",
		point.KVs{
			point.NewKV("message", "test"),
			point.NewKV("nums", []int{1, 2, 3}),
		},
		point.DefaultLoggingOptions()...)
	pt := ptinput.PtWrap(point.Logging, raw)

	errR := runScript(runner, pt)
	if errR != nil {
		t.Fatal(errR.Error())
	}

	v, _, err := pt.Get("hit")
	assert.NoError(t, err)
	assert.Equal(t, true, v)
}

func TestPtKvsGetAppendLenChain(t *testing.T) {
	pl := `
	arr = pt_kvs_get("nums", true)
	arr = append(arr, 4)
	pt_kvs_set("size", len(arr))
	pt_kvs_set("arr", arr, false, true)
	`

	runner, err := NewTestingRunner(pl)
	if err != nil {
		t.Fatal(err)
	}

	raw := point.NewPoint("test",
		point.KVs{
			point.NewKV("message", "test"),
			point.NewKV("nums", []int{1, 2, 3}),
		},
		append(point.DefaultLoggingOptions(), point.WithTime(time.Now()))...,
	)
	pt := ptinput.PtWrap(point.Logging, raw)

	errR := runScript(runner, pt)
	if errR != nil {
		t.Fatal(errR.Error())
	}

	v, _, err := pt.Get("size")
	assert.NoError(t, err)
	assert.Equal(t, int64(4), v)

	v, _, err = pt.Get("arr")
	assert.NoError(t, err)
	assert.Equal(t, "[1,2,3,4]", v)

	v, dt, err := pt.GetRaw("arr")
	assert.NoError(t, err)
	assert.Equal(t, ast.List, dt)
	assert.Equal(t, []any{int64(1), int64(2), int64(3), int64(4)}, v)
}

func TestPtKvsSetMap(t *testing.T) {
	pl := `
	obj = {"a": 1, "b": "x"}
	pt_kvs_set("obj", obj)
	`

	runner, err := NewTestingRunner(pl)
	if err != nil {
		t.Fatal(err)
	}

	pt := ptinput.NewPlPt(point.Logging, "test", nil, map[string]any{"message": "test"}, time.Now())
	errR := runScript(runner, pt)
	if errR != nil {
		t.Fatal(errR.Error())
	}

	v, dt, err := pt.GetRaw("obj")
	assert.NoError(t, err)
	assert.Equal(t, ast.String, dt)
	assert.Equal(t, `{"a":1,"b":"x"}`, v)

	v, dt, err = pt.Get("obj")
	assert.NoError(t, err)
	assert.Equal(t, ast.String, dt)
	assert.Equal(t, `{"a":1,"b":"x"}`, v)
}

func TestPtKvsSetMessageMapCompatibility(t *testing.T) {
	pl := `
	result_msg = {}
	all_keys = pt_kvs_keys()
	for k in all_keys {
		if k != "message" {
			result_msg[k] = pt_kvs_get(k)
		}
	}
	pt_kvs_set("message", result_msg)
	`

	runner, err := NewTestingRunner(pl)
	if err != nil {
		t.Fatal(err)
	}

	pt := ptinput.NewPlPt(point.Logging, "test", nil, map[string]any{
		"message": "test",
		"a":       "x",
		"b":       int64(1),
	}, time.Now())
	errR := runScript(runner, pt)
	if errR != nil {
		t.Fatal(errR.Error())
	}

	v, _, err := pt.Get("message")
	assert.NoError(t, err)
	assert.Equal(t, `{"a":"x","b":1,"status":"info"}`, v)

	assert.Equal(t, `{"a":"x","b":1,"status":"info"}`, pt.Fields()["message"])
	assert.Equal(t, `{"a":"x","b":1,"status":"info"}`, pt.Point().KVs().InfluxFields()["message"])
}

func TestPtKvsGetMapCompatibility(t *testing.T) {
	pl := `
	obj = {"a": 1, "b": "x"}
	pt_kvs_set("obj", obj)
	add_key("obj_type", value_type(pt_kvs_get("obj")))
	`

	runner, err := NewTestingRunner(pl)
	if err != nil {
		t.Fatal(err)
	}

	pt := ptinput.NewPlPt(point.Logging, "test", nil, map[string]any{"message": "test"}, time.Now())
	errR := runScript(runner, pt)
	if errR != nil {
		t.Fatal(errR.Error())
	}

	v, _, err := pt.Get("obj_type")
	assert.NoError(t, err)
	assert.Equal(t, "str", v)
}

func TestPtKvsGetRawOption(t *testing.T) {
	pl := `
	pt_kvs_set("obj", {"a": 1}, false, true)
	pt_kvs_set("raw_type", value_type(pt_kvs_get("obj", true)))
	pt_kvs_set("plain_type", value_type(pt_kvs_get("obj", false)))
	pt_kvs_set("plain_value", pt_kvs_get("obj", false))
	`

	runner, err := NewTestingRunner(pl)
	if err != nil {
		t.Fatal(err)
	}

	pt := ptinput.NewPlPt(point.Logging, "test", nil, map[string]any{"message": "test"}, time.Now())
	errR := runScript(runner, pt)
	if errR != nil {
		t.Fatal(errR.Error())
	}

	v, _, err := pt.Get("raw_type")
	assert.NoError(t, err)
	assert.Equal(t, "map", v)

	v, _, err = pt.Get("plain_type")
	assert.NoError(t, err)
	assert.Equal(t, "str", v)

	v, _, err = pt.Get("plain_value")
	assert.NoError(t, err)
	assert.Equal(t, `{"a":1}`, v)
}

func TestPtKvsSetRawOption(t *testing.T) {
	pl := `
	obj = {"a": 1}
	pt_kvs_set("obj_raw", obj, false, true)
	pt_kvs_set("obj_str", obj, false, false)
	`

	runner, err := NewTestingRunner(pl)
	if err != nil {
		t.Fatal(err)
	}

	pt := ptinput.NewPlPt(point.Logging, "test", nil, map[string]any{"message": "test"}, time.Now())
	errR := runScript(runner, pt)
	if errR != nil {
		t.Fatal(errR.Error())
	}

	v, dt, err := pt.GetRaw("obj_raw")
	assert.NoError(t, err)
	assert.Equal(t, ast.Map, dt)
	assert.Equal(t, map[string]any{"a": int64(1)}, v)

	v, dt, err = pt.GetRaw("obj_str")
	assert.NoError(t, err)
	assert.Equal(t, ast.String, dt)
	assert.Equal(t, `{"a":1}`, v)
}
