// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"fmt"

	"github.com/DataDog/datadog-agent/pkg/obfuscate"
	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnSQLCoverDesc = runtimev2.FnDesc{
	Name: "sql_cover",
	Desc: "Obfuscate SQL query string.",
	Params: []*runtimev2.Param{
		{
			Name: "val",
			Desc: "The sql to obfuscate.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "The obfuscated sql.",
			Typs: []ast.DType{ast.String},
		},
		{
			Desc: "The obfuscate status.",
			Typs: []ast.DType{ast.Bool},
		},
	},
}

func FnSQLCoverCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnSQLCoverDesc.Params)
}

func FnSQLCover(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	o := obfuscate.NewObfuscator(obfuscate.Config{})

	val, err := runtimev2.GetParamString(ctx, funcExpr, FnSQLCoverDesc.Params, 0)
	if err != nil {
		return err
	}

	if v, errObfuscate := obfuscatedResource(o, "sql", val); errObfuscate != nil {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: "", T: ast.String},
			runtimev2.V{V: false, T: ast.Bool},
		)
		return nil
	} else {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: v, T: ast.String},
			runtimev2.V{V: true, T: ast.Bool},
		)
	}
	return nil
}

func obfuscatedResource(o *obfuscate.Obfuscator, typ, resource string) (string, error) {
	if typ != "sql" {
		return resource, nil
	}
	oq, err := o.ObfuscateSQLString(resource)
	if err != nil {
		err = fmt.Errorf("error obfuscating stats group resource %q: %w", resource, err)
		return "", err
	}
	return oq.Query, nil
}
