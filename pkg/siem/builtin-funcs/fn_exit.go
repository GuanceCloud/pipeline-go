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

var FnExitDesc = runtimev2.FnDesc{
	Name:   "exit",
	Desc:   "Exit the program",
	Params: []*runtimev2.Param{},
}

func FnExitCheck(ctx *runtimev2.Task, node *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, node, FnExitDesc.Params)
}

func FnExit(ctx *runtimev2.Task, node *ast.CallExpr) *errchain.PlError {
	ctx.SetExit()
	return nil
}
