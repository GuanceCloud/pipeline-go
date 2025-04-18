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

var FnValueTypeDesc = runtimev2.FnDesc{
	Name: "value_type",
	Desc: "Returns the type of the value.",
	Params: []*runtimev2.Param{
		{
			Name: "val",
			Desc: "The value to get the type of.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "Returns the type of the value. One of (`bool`, `int`, `float`, `str`, `list`, `map`, `nil`). If the value and the type is nil, returns `nil`.",
			Typs: []ast.DType{ast.String},
		},
	},
}

func FnValueTypeCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnValueTypeDesc.Params)
}

func FnValueType(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	val, err := runtimev2.GetParam(ctx, funcExpr, FnValueTypeDesc.Params, 0)
	if err != nil {
		return err
	}

	_, dtyp := ast.DectDataType(val)

	var v string
	switch dtyp { //nolint:exhaustive
	case ast.Bool:
		v = "bool"
	case ast.Int:
		v = "int"
	case ast.Float:
		v = "float"
	case ast.String:
		v = "str"
	case ast.List:
		v = "list"
	case ast.Map:
		v = "map"
	case ast.Nil:
		v = "nil"
	}

	ctx.Regs.ReturnAppend(
		runtimev2.V{V: v, T: ast.String},
	)

	return nil
}
