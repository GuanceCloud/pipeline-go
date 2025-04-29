// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnLenDesc = runtimev2.FnDesc{
	Name: "len",
	Desc: "Get the length of the value. If the value is a string, returns the length of the string. If the value is a list or map, returns the length of the list or map.",
	Params: []*runtimev2.Param{
		{
			Name: "val",
			Desc: "The value to get the length of.",
			Typs: []ast.DType{ast.Map, ast.List, ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "The length of the value.",
			Typs: []ast.DType{ast.Int},
		},
	},
}

func FnLenCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnLenDesc.Params)
}

func FnLen(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	val, err := runtimev2.GetParam(ctx, funcExpr, FnLenDesc.Params, 0)
	if err != nil {
		return err
	}
	switch val := val.(type) { //nolint:exhaustive
	case string:
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: int64(len(val)), T: ast.Int},
		)
	case []any:
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: int64(len(val)), T: ast.Int},
		)
	case map[string]any:
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: int64(len(val)), T: ast.Int},
		)
	default:
		return runtimev2.NewRunError(ctx, "unsupported type for len", funcExpr.NamePos)
	}
	return nil
}
