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
	"github.com/spf13/cast"
)

var FnCastDesc = runtimev2.FnDesc{
	Name: "cast",
	Desc: "Convert the value to the target type.",
	Params: []*runtimev2.Param{
		{
			Name: "val",
			Desc: "The value of the type to be converted.",
			Typs: []ast.DType{ast.Bool, ast.Int, ast.Float, ast.String},
		},
		{
			Name: "typ",
			Desc: "Target type. One of (`bool`, `int`, `float`, `str`).",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "The value after the conversion.",
			Typs: []ast.DType{ast.Bool, ast.Int, ast.Float, ast.String},
		},
	},
}

func FnCastCheck(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, expr, FnCastDesc.Params)
}

func FnCast(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	val, err := runtimev2.GetParam(ctx, expr, FnCastDesc.Params, 0)
	if err != nil {
		return err
	}

	targetTyp, err := runtimev2.GetParamString(ctx, expr, FnCastDesc.Params, 1)
	if err != nil {
		return err
	}

	val, typ := doCast(val, targetTyp)
	ctx.Regs.ReturnAppend(runtimev2.V{V: val, T: typ})
	return nil
}

func doCast(result any, tInfo string) (any, ast.DType) {
	switch strings.ToLower(tInfo) {
	case "bool":
		return cast.ToBool(result), ast.Bool

	case "int":
		return cast.ToInt64(cast.ToFloat64(result)), ast.Int

	case "float":
		return cast.ToFloat64(result), ast.Float

	case "str":
		return cast.ToString(result), ast.String
	}

	return nil, ast.Nil
}
