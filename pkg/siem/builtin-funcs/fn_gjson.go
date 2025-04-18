// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
	"github.com/tidwall/gjson"
)

var FnGJSONDesc = runtimev2.FnDesc{
	Name: "gjson",
	Desc: "GJSON provides a fast and easy way to get values from a JSON document.",
	Params: []*runtimev2.Param{
		{
			Name: "input",
			Desc: "JSON format string to parse.",
			Typs: []ast.DType{ast.String},
		},
		{
			Name: "json_path",
			Desc: "JSON path.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "Parsed result.",
			Typs: ast.AllTyp(),
		},
		{
			Desc: "Parsed status.",
			Typs: []ast.DType{ast.Bool},
		},
	},
}

func FnGJSONCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnGJSONDesc.Params)
}

func FnGJSON(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	jsonStr, err := runtimev2.GetParamString(ctx, funcExpr, FnGJSONDesc.Params, 0)
	if err != nil {
		return err
	}

	jpath, err := runtimev2.GetParamString(ctx, funcExpr, FnGJSONDesc.Params, 1)
	if err != nil {
		return err
	}

	res := gjson.Get(jsonStr, jpath)

	rType := res.Type
	var v runtimev2.V
	switch rType {
	case gjson.Number:
		v.V = res.Float()
		v.T = ast.Float
	case gjson.True, gjson.False:
		v.V = res.Bool()
		v.T = ast.Bool
	case gjson.Null:
		v.V = nil
		v.T = ast.Nil
	case gjson.String:
		v.V = res.String()
		v.T = ast.String
	case gjson.JSON:
		if res.IsObject() {
			val := res.Value()
			v.V = val.(map[string]any)
			v.T = ast.Map
		} else {
			val := res.Value()
			v.V = val.([]any)
			v.T = ast.List
		}
	default:
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: map[string]any(nil), T: ast.Map},
			runtimev2.V{V: false, T: ast.Bool},
		)
		return nil
	}

	ctx.Regs.ReturnAppend(
		v,
		runtimev2.V{V: true, T: ast.Bool},
	)
	return nil
}
