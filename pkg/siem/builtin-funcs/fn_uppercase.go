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

var FnUppercaseDesc = runtimev2.FnDesc{
	Name: "uppercase",
	Desc: "Converts a string to uppercase.",
	Params: []*runtimev2.Param{
		{
			Name: "val",
			Desc: "The string to convert.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "Returns the uppercase value.",
			Typs: []ast.DType{ast.String},
		},
	}}

func FnUppercaseCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnUppercaseDesc.Params)
}

func FnUppercase(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	val, err := runtimev2.GetParamString(ctx, funcExpr, FnUppercaseDesc.Params, 0)
	if err != nil {
		return err
	}
	v := strings.ToUpper(val)
	ctx.Regs.ReturnAppend(runtimev2.V{V: v, T: ast.String})
	return nil
}
