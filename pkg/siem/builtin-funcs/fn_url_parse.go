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

var FnURLParseDesc = runtimev2.FnDesc{
	Name: "url_parse",
	Desc: "Parses a URL and returns it as a map.",
	Params: []*runtimev2.Param{
		{
			Name: "url",
			Desc: "The URL to parse.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "Returns the parsed URL as a map.",
			Typs: []ast.DType{ast.Map},
		},
		{
			Desc: "Returns true if the URL is valid.",
			Typs: []ast.DType{ast.Bool},
		},
	},
}

func FnURLParseCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnURLParseDesc.Params)
}

func FnURLParse(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	urlStr, err := runtimev2.GetParamString(ctx, funcExpr, FnURLParseDesc.Params, 0)
	if err != nil {
		return err
	}
	uu, errParse := url.Parse(urlStr)
	if errParse != nil {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: (map[string]any)(nil), T: ast.Map},
			runtimev2.V{V: false, T: ast.Bool},
		)
		return nil
	}

	params := make(map[string]any)
	for k, vs := range uu.Query() {
		vals := make([]any, 0, len(vs))
		for i := range vs {
			vals = append(vals, vs[i])
		}
		params[k] = vals
	}
	res := map[string]any{
		"scheme": uu.Scheme,
		"host":   uu.Host,
		"port":   uu.Port(),
		"path":   uu.Path,
		"params": params,
	}
	ctx.Regs.ReturnAppend(
		runtimev2.V{V: res, T: ast.Map},
		runtimev2.V{V: true, T: ast.Bool},
	)
	return nil
}
