// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"time"

	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnTimeNowDesc = runtimev2.FnDesc{
	Name: "time_now",
	Desc: "Get current timestamp with the specified precision.",
	Params: []*runtimev2.Param{
		{
			Name: "precision",
			Desc: "The precision of the timestamp. Supported values: `ns`, `us`, `ms`, `s`.",
			Typs: []ast.DType{ast.String},
			Val:  func() any { return "ns" },
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "Returns the current timestamp.",
			Typs: []ast.DType{ast.Int},
		},
	},
}

func FnTimeNowCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnTimeNowDesc.Params)
}

func FnTimenow(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	precision, err := runtimev2.GetParamString(ctx, funcExpr, FnTimeNowDesc.Params, 0)
	if err != nil {
		return err
	}

	var ts int64
	switch precision {
	case "us":
		ts = time.Now().UnixMicro()
	case "ms":
		ts = time.Now().UnixMilli()
	case "s":
		ts = time.Now().Unix()
	default:
		ts = time.Now().UnixNano()
	}
	ctx.Regs.ReturnAppend(
		runtimev2.V{V: ts, T: ast.Int},
	)
	return nil
}
