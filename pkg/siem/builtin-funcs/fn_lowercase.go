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

var FnLowercaseDesc = runtimev2.FnDesc{
	Name: "lowercase",
	Desc: "Converts a string to lowercase.",
	Params: []*runtimev2.Param{
		{
			Name: "val",
			Desc: "The string to convert.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "Returns the lowercase value.",
			Typs: []ast.DType{ast.String},
		},
	},
}

func FnLowercaseCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnLowercaseDesc.Params)
}

func FnLowercase(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	val, err := runtimev2.GetParamString(ctx, funcExpr, FnLowercaseDesc.Params, 0)
	if err != nil {
		return err
	}
	v := strings.ToLower(val)
	ctx.Regs.ReturnAppend(runtimev2.V{V: v, T: ast.String})
	return nil
}
