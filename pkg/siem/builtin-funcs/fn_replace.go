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

var FnReplaceDesc = runtimev2.FnDesc{
	Name: "replace",
	Desc: "Replaces text in a string.",
	Params: []*runtimev2.Param{
		{
			Name: "input",
			Desc: "The string to replace text in.",
			Typs: []ast.DType{ast.String},
		},
		{
			Name: "pattern",
			Desc: "Regular expression pattern.",
			Typs: []ast.DType{ast.String},
		},
		{
			Name: "replacement",
			Desc: "Replacement text to use.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "The string with text replaced.",
			Typs: []ast.DType{ast.String},
		},
	},
}

func FnReplaceCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnReplaceDesc.Params)
}

func FnReplace(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	val, err := runtimev2.GetParamString(ctx, funcExpr, FnReplaceDesc.Params, 0)
	if err != nil {
		return err
	}
	pattern, err := runtimev2.GetParamString(ctx, funcExpr, FnReplaceDesc.Params, 1)
	if err != nil {
		return err
	}
	replacement, err := runtimev2.GetParamString(ctx, funcExpr, FnReplaceDesc.Params, 2)
	if err != nil {
		return err
	}

	if re, err := regexpCache.Get(pattern); err != nil {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: val, T: ast.String},
			runtimev2.V{V: false, T: ast.Bool},
		)
	} else {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: re.ReplaceAllString(val, replacement), T: ast.String},
			runtimev2.V{V: true, T: ast.Bool},
		)
	}
	return nil
}
