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

var FnDeleteDesc = runtimev2.FnDesc{
	Name: "delete",
	Desc: "Delete key from the map.",
	Params: []*runtimev2.Param{
		{
			Name: "m",
			Desc: "The map for deleting key",
			Typs: []ast.DType{ast.Map},
		},
		{
			Name: "key",
			Desc: "Key need delete from map.",
			Typs: []ast.DType{ast.String},
		},
	},
}

func FnDeleteCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnDeleteDesc.Params)
}

func FnDelete(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	val, err := runtimev2.GetParamMap(ctx, funcExpr, FnDeleteDesc.Params, 0)
	if err != nil {
		return err
	}
	key, err := runtimev2.GetParamString(ctx, funcExpr, FnDeleteDesc.Params, 1)
	if err != nil {
		return err
	}
	delete(val, key)

	return nil
}
