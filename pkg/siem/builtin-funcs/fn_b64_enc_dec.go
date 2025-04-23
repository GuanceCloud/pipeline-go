package funcs

import (
	"encoding/base64"

	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnB64DecDesc = runtimev2.FnDesc{
	Name: "b64dec",
	Desc: "Base64 decoding.",
	Params: []*runtimev2.Param{
		{
			Name: "data",
			Desc: "Data that needs to be base64 decoded.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "The decoded string.",
			Typs: []ast.DType{ast.String},
		},
		{
			Desc: "Whether decoding is successful.",
			Typs: []ast.DType{ast.Bool},
		},
	},
}

var FnB64EncDesc = runtimev2.FnDesc{
	Name: "b64enc",
	Desc: "Base64 encoding.",
	Params: []*runtimev2.Param{
		{
			Name: "data",
			Desc: "Data that needs to be base64 encoded.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "The encoded string.",
			Typs: []ast.DType{ast.String},
		},
		{
			Desc: "Whether encoding is successful.",
			Typs: []ast.DType{ast.Bool},
		},
	},
}

func FnB64DecCheck(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, expr, FnB64DecDesc.Params)
}

func FnB64Dec(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	data, err := runtimev2.GetParamString(ctx, funcExpr, FnB64DecDesc.Params, 0)
	if err != nil {
		return err
	}
	if res, err := base64.StdEncoding.DecodeString(data); err != nil {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: "", T: ast.String},
			runtimev2.V{V: false, T: ast.Bool})
	} else {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: res, T: ast.String},
			runtimev2.V{V: false, T: ast.Bool})
	}
	return nil
}

func FnB64EncCheck(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, expr, FnB64EncDesc.Params)

}

func FnB64Enc(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	data, err := runtimev2.GetParamString(ctx, funcExpr, FnB64DecDesc.Params, 0)
	if err != nil {
		return err
	}
	res := base64.StdEncoding.EncodeToString([]byte(data))
	ctx.Regs.ReturnAppend(
		runtimev2.V{V: res, T: ast.String})
	return nil
}
