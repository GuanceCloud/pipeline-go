package funcs

import (
	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnAppendDesc = runtimev2.FnDesc{
	Name: "append",
	Desc: "Appends a value to a list.",
	Params: []*runtimev2.Param{
		{
			Name: "li",
			Desc: "The list to append to.",
			Typs: []ast.DType{ast.List},
		},
		{
			Name:     "v",
			Desc:     "The value to append.",
			Variable: true,
			Typs:     ast.AllTyp(),
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "The list with the appended value.",
			Typs: []ast.DType{ast.List},
		},
	},
}

func FnAppendCheck(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, expr, FnAppendDesc.Params)
}

func FnAppend(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	li, err := runtimev2.GetParamList(ctx, expr, FnAppendDesc.Params, 0)
	if err != nil {
		return err
	}
	elems, err := runtimev2.GetParamList(ctx, expr, FnAppendDesc.Params, 1)
	if err != nil {
		return err
	}
	li = append(li, elems...)
	ctx.Regs.ReturnAppend(runtimev2.V{V: li, T: ast.List})

	return nil
}
