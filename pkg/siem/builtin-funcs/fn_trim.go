// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"strings"

	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnTrimDesc = runtimev2.FnDesc{
	Name: "trim",
	Desc: "Removes leading and trailing whitespace from a string.",
	Params: []*runtimev2.Param{
		{
			Name: "val",
			Desc: "The string to trim.",
			Typs: []ast.DType{ast.String},
		},
		{
			Name: "cutset",
			Desc: "Characters to remove from the beginning and end of the string. If not specified, whitespace is removed.",
			Typs: []ast.DType{ast.String},
			Val:  func() any { return "" },
		},
		{
			Name: "side",
			Desc: "The side to trim from. If value is 0, trim from both sides. If value is 1, trim from the left side. If value is 2, trim from the right side.",
			Typs: []ast.DType{ast.Int},
			Val:  func() any { return 0 },
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "The trimmed string.",
			Typs: []ast.DType{ast.String},
		},
	},
}

func FnTrimCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnTrimDesc.Params)
}

var spaceCut = func() string {
	s := []rune{'\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0}
	var str string
	for _, v := range s {
		str += string(v)
	}
	return str
}()

func FnTrim(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	val, err := runtimev2.GetParamString(ctx, funcExpr, FnTrimDesc.Params, 0)
	if err != nil {
		return err
	}

	cutset, err := runtimev2.GetParamString(ctx, funcExpr, FnTrimDesc.Params, 1)
	if err != nil {
		return err
	}

	if cutset == "" {
		cutset = spaceCut
	}

	side, err := runtimev2.GetParamInt(ctx, funcExpr, FnTrimDesc.Params, 2)
	if err != nil {
		return err
	}

	switch side {
	case 1:
		val = strings.TrimLeft(val, cutset)
	case 2:
		val = strings.TrimRight(val, cutset)
	default:
		val = strings.Trim(val, cutset)
	}

	ctx.Regs.ReturnAppend(
		runtimev2.V{V: val, T: ast.String},
	)

	return nil
}
