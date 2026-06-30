package platypus

import (
	"fmt"
	"testing"
	"time"

	"github.com/GuanceCloud/cliutils/point"
	"github.com/GuanceCloud/pipeline-go/constants"
	"github.com/GuanceCloud/pipeline-go/lang"
	"github.com/GuanceCloud/pipeline-go/ptinput"
	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtime"
	"github.com/GuanceCloud/platypus/pkg/errchain"
	"github.com/stretchr/testify/assert"
)

func TestScript(t *testing.T) {
	ret, retErr := NewScripts(map[string]string{
		"abc.p": "if true {}",
	}, lang.WithNS(constants.NSGitRepo),
		lang.WithCat(point.Logging))

	if len(retErr) > 0 {
		t.Fatal(retErr)
	}

	s := ret["abc.p"]

	if ng := s.Engine(); ng == nil {
		t.Fatalf("no engine")
	}
	plpt := ptinput.NewPlPt(point.Logging, "ng", nil, nil, time.Now())
	err := s.Run(plpt, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, plpt.Fields(), map[string]interface{}{"status": constants.DefaultStatus})
	assert.Equal(t, 0, len(plpt.Tags()))
	assert.Equal(t, "abc.p", s.Name())
	assert.Equal(t, point.Logging, s.Category())
	assert.Equal(t, s.NS(), constants.NSGitRepo)

	//nolint:dogsled
	plpt = ptinput.NewPlPt(point.Logging, "ng", nil, nil, time.Now())
	err = s.Run(plpt, nil, &lang.LogOption{DisableAddStatusField: true})
	if err != nil {
		t.Fatal(err)
	}

	if len(plpt.Fields()) != 1 {
		t.Fatal(plpt.Fields())
	} else {
		if _, ok := plpt.Fields()["status"]; !ok {
			t.Fatal("without status")
		}
	}

	//nolint:dogsled
	plpt = ptinput.NewPlPt(point.Logging, "ng", nil, nil, time.Now())
	err = s.Run(plpt, nil, &lang.LogOption{
		DisableAddStatusField: false,
		IgnoreStatus:          []string{constants.DefaultStatus},
	})
	if err != nil {
		t.Fatal(err)
	}
	if plpt.Dropped() != true {
		t.Fatal("!drop")
	}
}

func TestDrop(t *testing.T) {
	ret, retErr := NewScripts(map[string]string{
		"abc.p": "add_key(a, \"a\"); add_key(status, \"debug\"); drop(); add_key(b, \"b\")"},
		lang.WithNS(constants.NSGitRepo),
		lang.WithCat(point.Logging))
	if len(retErr) > 0 {
		t.Fatal(retErr)
	}

	s := ret["abc.p"]

	plpt := ptinput.NewPlPt(point.Logging, "ng", nil, nil, time.Now())
	if err := s.Run(plpt, nil, nil); err != nil {
		t.Fatal(err)
	}

	if plpt.Dropped() != true {
		t.Error("drop != true")
	}
}

func TestScriptRunWithRuntimeOptsPrivate(t *testing.T) {
	const privateKey = "test_result"

	setResult := func(ctx *runtime.Task, funcExpr *ast.CallExpr) *errchain.PlError {
		if len(funcExpr.Param) != 1 {
			return runtime.NewRunError(ctx, "set_result expected 1 arg", funcExpr.NamePos)
		}
		val, _, err := runtime.RunStmt(ctx, funcExpr.Param[0])
		if err != nil {
			return err
		}
		result, ok := ctx.PValue(privateKey)
		if !ok {
			return nil
		}
		p, ok := result.(*string)
		if !ok {
			return runtime.NewRunError(ctx, fmt.Sprintf("%s must be *string", privateKey), funcExpr.NamePos)
		}
		*p = fmt.Sprint(val)
		return nil
	}

	ret, retErr := NewScripts(map[string]string{
		"abc.p": `if service == "api" { set_result("hit") }`,
	}, lang.WithCat(point.Logging), lang.WithFn(
		map[string]runtime.FuncCall{"set_result": setResult},
		map[string]runtime.FuncCheck{"set_result": func(ctx *runtime.Task, funcExpr *ast.CallExpr) *errchain.PlError {
			if len(funcExpr.Param) != 1 {
				return runtime.NewRunError(ctx, "set_result expected 1 arg", funcExpr.NamePos)
			}
			return nil
		}},
	))
	if len(retErr) > 0 {
		t.Fatal(retErr)
	}

	pt := point.NewPoint("abc", point.NewKVs(map[string]any{
		"service": "api",
	}), point.DefaultLoggingOptions()...)
	plpt := ptinput.PtWrap(point.Logging, pt)

	var got string
	err := ret["abc.p"].RunWithRuntimeOpts(plpt, nil, nil, runtime.WithPrivate(map[string]any{
		privateKey: &got,
	}))
	assert.NoError(t, err)
	assert.Equal(t, "hit", got)
	assert.Nil(t, plpt.Point().KVs().Get("test_result"))
}
