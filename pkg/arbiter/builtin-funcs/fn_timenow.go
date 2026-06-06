// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"fmt"
	"time"

	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/dql"
	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnTimeNowDesc = runtimev2.FnDesc{
	Name: "time_now",
	Desc: "Get the DQL query start timestamp with the specified precision.",
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
			Desc: "Returns the DQL query start timestamp.",
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

	start, err := dqlQueryStartTime(ctx, funcExpr)
	if err != nil {
		return err
	}

	switch precision {
	case "us":
		start *= int64(time.Millisecond / time.Microsecond)
	case "ms":
	case "s":
		start /= int64(time.Second / time.Millisecond)
	default:
		start *= int64(time.Millisecond / time.Nanosecond)
	}
	ctx.Regs.ReturnAppend(
		runtimev2.V{V: start, T: ast.Int},
	)
	return nil
}

func dqlQueryStartTime(ctx *runtimev2.Task, expr *ast.CallExpr) (int64, *errchain.PlError) {
	v, ok := ctx.PValue(PDQLCli)
	if !ok {
		return 0, runtimev2.NewRunError(ctx, fmt.Sprintf(
			"missing context data named %s", PDQLCli), expr.NamePos)
	}
	dqlCli, ok := v.(dql.DQL)
	if !ok {
		return 0, runtimev2.NewRunError(ctx, fmt.Sprintf(
			"context data %s type is expected", PDQLCli), expr.NamePos)
	}

	r := dqlCli.TimeRange()
	if len(r) == 2 {
		return r[0], nil
	}
	return genTimeRange15min(time.Now().UnixMilli()), nil
}
