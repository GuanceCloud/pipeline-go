// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"strconv"

	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnParseIntDesc = runtimev2.FnDesc{
	Name: "parse_int",
	Desc: "Parses a string into an integer.",
	Params: []*runtimev2.Param{
		{
			Name: "val",
			Desc: "The string to parse.",
			Typs: []ast.DType{ast.String},
		},
		{
			Name: "base",
			Desc: "The base to use for parsing.",
			Typs: []ast.DType{ast.Int},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "The parsed integer.",
			Typs: []ast.DType{ast.Int},
		},
		{
			Desc: "Whether the parsing was successful.",
			Typs: []ast.DType{ast.Bool},
		},
	},
}

var FormatIntDesc = runtimev2.FnDesc{
	Name: "format_int",
	Desc: "Formats an integer into a string.",
	Params: []*runtimev2.Param{
		{
			Name: "val",
			Desc: "The integer to format.",
			Typs: []ast.DType{ast.Int},
		},
		{
			Name: "base",
			Desc: "The base to use for formatting.",
			Typs: []ast.DType{ast.Int},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "The formatted string.",
			Typs: []ast.DType{ast.String},
		},
	},
}

func FnParseIntCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnParseIntDesc.Params)
}

func FnParseInt(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	val, err := runtimev2.GetParamString(ctx, funcExpr, FnParseIntDesc.Params, 0)
	if err != nil {
		return err
	}
	base, err := runtimev2.GetParamInt(ctx, funcExpr, FnParseIntDesc.Params, 1)
	if err != nil {
		return err
	}

	if v, err := strconv.ParseInt(val, int(base), 64); err != nil {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: int64(0), T: ast.Int},
			runtimev2.V{V: false, T: ast.Bool},
		)
	} else {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: v, T: ast.Int},
			runtimev2.V{V: true, T: ast.Bool},
		)
	}
	return nil
}

func FormatIntChecking(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FormatIntDesc.Params)
}

func FormatInt(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	val, err := runtimev2.GetParamInt(ctx, funcExpr, FormatIntDesc.Params, 0)
	if err != nil {
		return err
	}
	base, err := runtimev2.GetParamInt(ctx, funcExpr, FormatIntDesc.Params, 1)
	if err != nil {
		return err
	}

	v := strconv.FormatInt(val, int(base))
	ctx.Regs.ReturnAppend(
		runtimev2.V{V: v, T: ast.String},
	)
	return nil
}
