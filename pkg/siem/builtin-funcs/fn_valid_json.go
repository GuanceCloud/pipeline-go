// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"encoding/json"

	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnValidJSONDesc = runtimev2.FnDesc{
	Name: "valid_json",
	Desc: "Returns true if the value is a valid JSON.",
	Params: []*runtimev2.Param{
		{
			Name: "val",
			Desc: "The value to check.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "Returns true if the value is a valid JSON.",
			Typs: []ast.DType{ast.Bool},
		},
	},
}

func FnValidJSONCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnValidJSONDesc.Params)
}

func FnValidJSON(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	val, err := runtimev2.GetParamString(ctx, funcExpr, FnValidJSONDesc.Params, 0)
	if err != nil {
		return err
	}

	valid := json.Valid([]byte(val))
	ctx.Regs.ReturnAppend(
		runtimev2.V{V: valid, T: ast.Bool},
	)
	return nil
}
