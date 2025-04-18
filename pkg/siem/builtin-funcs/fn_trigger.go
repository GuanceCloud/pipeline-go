package funcs

import (
	"github.com/GuanceCloud/pipeline-go/pkg/siem/trigger"
	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnTriggerDesc = runtimev2.FnDesc{
	Name: "trigger",
	Desc: "Trigger a security event.",
	Params: []*runtimev2.Param{
		{
			Name: "result",
			Desc: "Event check result.",
			Typs: []ast.DType{ast.Int, ast.Float, ast.Bool, ast.String},
		},
		{
			Name: "level",
			Desc: "Event Level. One of: (`critical`, `high`, `medium`, `low`, `info`).",
			Typs: []ast.DType{ast.String},
			Val: func() any {
				return ""
			},
		},
		{
			Name: "dim_tags",
			Desc: "Dimension tags.",
			Typs: []ast.DType{ast.Map},
			Val: func() any {
				return map[string]any{}
			},
		},
		{
			Name: "extra_data",
			Desc: "Extra data.",
			Typs: []ast.DType{ast.Map},
			Val: func() any {
				return map[string]any{}
			},
		},
	},
}

func FnTriggerCheck(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, expr, FnTriggerDesc.Params)
}

func FnTrigger(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	var tr *trigger.Trigger
	if v, ok := ctx.PValue(PTrigger); !ok {
		return nil
	} else {
		if v, ok := v.(*trigger.Trigger); ok {
			tr = v
		}
	}
	if tr == nil {
		return nil
	}

	result, err := runtimev2.GetParam(ctx, expr, FnDQLDesc.Params, 0)
	if err != nil {
		return err
	}

	level, err := runtimev2.GetParamString(ctx, expr, FnTriggerDesc.Params, 1)
	if err != nil {
		return err
	}
	dimTags, err := runtimev2.GetParamMap(ctx, expr, FnTriggerDesc.Params, 2)
	if err != nil {
		return err
	}
	extraData, err := runtimev2.GetParamMap(ctx, expr, FnTriggerDesc.Params, 3)
	if err != nil {
		return err
	}
	tr.Trigger(result, level, dimTags, extraData)
	return nil
}
