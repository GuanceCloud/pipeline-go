// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"bytes"

	"github.com/goccy/go-json"

	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnDumpJSONDesc = runtimev2.FnDesc{
	Name: "dump_json",
	Desc: "Returns the JSON encoding of v.",
	Params: []*runtimev2.Param{
		{
			Name: "v",
			Desc: "Object to encode.",
			Typs: []ast.DType{ast.String},
		},
		{
			Name: "indent",
			Desc: "Indentation prefix.",
			Val: func() any {
				return ""
			},
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "JSON encoding of v.",
			Typs: []ast.DType{ast.String},
		},
		{
			Desc: "Whether decoding is successful.",
			Typs: []ast.DType{ast.Bool},
		},
	},
}

func FnDumpJSONCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnDumpJSONDesc.Params)
}

func FnDumpJSON(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	val, err := runtimev2.GetParam(ctx, funcExpr, FnDumpJSONDesc.Params, 0)
	if err != nil {
		return err
	}

	indent, err := runtimev2.GetParamString(ctx, funcExpr, FnDumpJSONDesc.Params, 1)
	if err != nil {
		return err
	}

	b := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(b)
	if indent != "" {
		enc.SetIndent("", indent)
	}
	if err := enc.Encode(val); err != nil {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: "", T: ast.String},
			runtimev2.V{V: false, T: ast.Bool},
		)
	} else {
		cxx := b.String()
		_ = cxx
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: b.String(), T: ast.String},
			runtimev2.V{V: true, T: ast.Bool},
		)
	}
	return nil
}
