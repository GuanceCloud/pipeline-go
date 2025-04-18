// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
	"github.com/mssola/user_agent"
)

var FnUserAgentDesc = runtimev2.FnDesc{
	Name: "user_agent",
	Desc: "Parses a User-Agent header.",
	Params: []*runtimev2.Param{
		{
			Name: "header",
			Desc: "The User-Agent header to parse.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "Returns the parsed User-Agent header as a map.",
			Typs: []ast.DType{ast.Map},
		},
	},
}

func FnUserAgentCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnUserAgentDesc.Params)
}

func FnUserAgent(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {

	header, err := runtimev2.GetParamString(ctx, funcExpr, FnUserAgentDesc.Params, 0)
	if err != nil {
		return err
	}
	dic := UserAgentHandle(header)
	ctx.Regs.ReturnAppend(runtimev2.V{V: dic, T: ast.Map})
	return nil
}

func UserAgentHandle(str string) map[string]interface{} {
	res := make(map[string]interface{})

	ua := user_agent.New(str)

	res["isMobile"] = ua.Mobile()
	res["isBot"] = ua.Bot()
	res["os"] = ua.OS()

	name, version := ua.Browser()
	res["browser"] = name
	res["browserVer"] = version

	en, v := ua.Engine()
	res["engine"] = en
	res["engineVer"] = v

	res["ua"] = ua.Platform()
	return res
}
