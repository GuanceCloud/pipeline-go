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

var FnParseDurationDesc = runtimev2.FnDesc{
	Name: "parse_duration",
	Desc: "Parses a golang duration string into a duration. " +
		"A duration string is a sequence of possibly signed decimal numbers with optional fraction and unit suffixes for each number, such as `300ms`, `-1.5h` or `2h45m`. " +
		"Valid units are `ns`, `us` (or `μs`), `ms`, `s`, `m`, `h`. ",
	Params: []*runtimev2.Param{
		{
			Name: "s",
			Desc: "The string to parse.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "The duration in nanoseconds.",
			Typs: []ast.DType{ast.Int},
		},
		{
			Desc: "The duration in nanoseconds.",
			Typs: []ast.DType{ast.Int},
		},
	},
}

func FnParseDurationCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnParseDurationDesc.Params)
}

func FnParseDuration(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	duStr, err := runtimev2.GetParamString(ctx, funcExpr, FnParseDurationDesc.Params, 0)
	if err != nil {
		return err
	}
	if du, err := time.ParseDuration(duStr); err != nil {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: int64(0), T: ast.Int},
			runtimev2.V{V: false, T: ast.Bool},
		)
		return nil
	} else {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: du.Nanoseconds(), T: ast.Int},
			runtimev2.V{V: true, T: ast.Bool},
		)
	}
	return nil
}
