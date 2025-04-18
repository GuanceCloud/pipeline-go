// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"github.com/goccy/go-json"

	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnLoadJSONDesc = runtimev2.FnDesc{
	Name: "load_json",
	Desc: "Unmarshal json string",
	Params: []*runtimev2.Param{
		{
			Name: "val",
			Desc: "JSON string.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "Unmarshal result.",
			Typs: ast.AllTyp(),
		},
		{
			Desc: "Unmarshal status.",
			Typs: []ast.DType{ast.Bool},
		},
	},
}

func FnLoadJSONCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnLoadJSONDesc.Params)
}

func FnLoadJSON(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	val, err := runtimev2.GetParamString(ctx, funcExpr, FnLoadJSONDesc.Params, 0)
	if err != nil {
		return err
	}

	var m any
	if err := json.Unmarshal([]byte(val), &m); err != nil {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: nil, T: ast.Nil},
			runtimev2.V{V: false, T: ast.Bool},
		)
		return nil
	} else {
		m, dtype := ast.DectDataType(m)
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: m, T: dtype},
			runtimev2.V{V: true, T: ast.Bool},
		)
	}
	return nil
}
