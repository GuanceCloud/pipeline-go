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

var FnStrJoinDesc = runtimev2.FnDesc{
	Name: "str_join",
	Desc: "String join.",
	Params: []*runtimev2.Param{
		{
			Name: "li",
			Desc: "List to be joined with separator. The elements type need to be string, if not, they will be ignored.",
			Typs: []ast.DType{ast.List},
		},
		{
			Name: "sep",
			Desc: "Separator to be used between elements.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "Joined string.",
			Typs: []ast.DType{ast.String},
		},
	},
}

func FnStrJoinCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnStrJoinDesc.Params)
}

func FnStrJoin(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	val, err := runtimev2.GetParamList(ctx, funcExpr, FnDeleteDesc.Params, 0)
	if err != nil {
		return err
	}
	sep, err := runtimev2.GetParamString(ctx, funcExpr, FnDeleteDesc.Params, 1)
	if err != nil {
		return err
	}

	var li []string
	for _, v := range val {
		if v, ok := v.(string); ok {
			li = append(li, v)
		}
	}

	r := strings.Join(li, sep)
	ctx.Regs.ReturnAppend(runtimev2.V{
		V: r,
		T: ast.String,
	})

	return nil
}
