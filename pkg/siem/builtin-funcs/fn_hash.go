// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"

	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnHashDesc = runtimev2.FnDesc{
	Name: "hash",
	Params: []*runtimev2.Param{
		{
			Name: "text",
			Desc: "The string used to calculate the hash.",
			Typs: []ast.DType{ast.String},
		},
		{
			Name: "method",
			Desc: "Hash Algorithms, allowing values including `md5`, `sha1`, `sha256`, `sha512`.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "The hash value.",
			Typs: []ast.DType{ast.String},
		},
	},
}

func FnHashCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnHashDesc.Params)
}

func FnHash(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	text, err := runtimev2.GetParamString(ctx, funcExpr, FnHashDesc.Params, 0)
	if err != nil {
		return err
	}

	method, err := runtimev2.GetParamString(ctx, funcExpr, FnHashDesc.Params, 1)
	if err != nil {
		return err
	}

	var sum string
	switch method {
	case "md5":
		b := md5.Sum([]byte(text))
		sum = hex.EncodeToString(b[:])
	case "sha1":
		b := sha1.Sum([]byte(text))
		sum = hex.EncodeToString(b[:])
	case "sha256":
		b := sha256.Sum256([]byte(text))
		sum = hex.EncodeToString(b[:])
	case "sha512":
		b := sha512.Sum512([]byte(text))
		sum = hex.EncodeToString(b[:])
	default:
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: "", T: ast.String},
		)
		return nil
	}

	ctx.Regs.ReturnAppend(
		runtimev2.V{V: sum, T: ast.String},
	)
	return nil

}
