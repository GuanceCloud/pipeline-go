// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"net/url"

	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnURLDecodeDesc = runtimev2.FnDesc{
	Name: "url_decode",
	Desc: "Decodes a URL-encoded string.",
	Params: []*runtimev2.Param{
		{
			Name: "val",
			Desc: "The URL-encoded string to decode.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "The decoded string.",
			Typs: []ast.DType{ast.String},
		},
		{
			Desc: "The decoding status.",
			Typs: []ast.DType{ast.Bool},
		},
	},
}

func FnURLDecodeCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnURLDecodeDesc.Params)
}

func FnURLDecode(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	val, err := runtimev2.GetParamString(ctx, funcExpr, FnURLDecodeDesc.Params, 0)
	if err != nil {
		return err
	}

	if v, err := url.QueryUnescape(val); err != nil {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: "", T: ast.String},
			runtimev2.V{V: false, T: ast.Bool},
		)
	} else {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: v, T: ast.String},
			runtimev2.V{V: true, T: ast.Bool},
		)
	}

	return nil
}
