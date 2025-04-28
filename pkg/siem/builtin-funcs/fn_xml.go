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
	"github.com/antchfx/xmlquery"
)

var FnXMLDesc = runtimev2.FnDesc{
	Name: "xml_query",
	Desc: "Returns the value of an XML field.",
	Params: []*runtimev2.Param{
		{
			Name: "input",
			Desc: "The XML input to get the value of.",
			Typs: []ast.DType{ast.String},
		},
		{
			Name: "xpath",
			Desc: "The XPath expression to get the value of.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "Returns the value of the XML field.",
			Typs: []ast.DType{ast.String},
		},
		{
			Desc: "Returns true if the field exists, false otherwise.",
			Typs: []ast.DType{ast.Bool},
		},
	},
}

func FnXMLChecking(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnXMLDesc.Params)
}

func FnXML(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	input, err := runtimev2.GetParamString(ctx, funcExpr, FnXMLDesc.Params, 0)
	if err != nil {
		return err
	}

	xpathExpr, err := runtimev2.GetParamString(ctx, funcExpr, FnXMLDesc.Params, 1)
	if err != nil {
		return err
	}

	doc, errParse := xmlquery.Parse(strings.NewReader(input))
	if errParse != nil {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: "", T: ast.String},
			runtimev2.V{V: false, T: ast.Bool})
		return nil
	}

	if dest, err := xmlquery.Query(doc, xpathExpr); err != nil {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: "", T: ast.String},
			runtimev2.V{V: false, T: ast.Bool})
		return nil
	} else {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: dest.InnerText(), T: ast.String},
			runtimev2.V{V: true, T: ast.Bool})
	}
	return nil
}
