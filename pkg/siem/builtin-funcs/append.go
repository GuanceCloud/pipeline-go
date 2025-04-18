package funcs

import (
	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var AppendDesc = runtimev2.FnDesc{
	Name: "append",
	Desc: "Adding element to a list",
	Params: []*runtimev2.Param{
		{
			Name: "li",
			Desc: "a list",
			Typs: []ast.DType{ast.List},
		},
		{
			Name:     "elems",
			Desc:     "elements",
			Variable: true,
			Typs:     ast.AllTyp(),
		},
	},
}

func AppendChecking(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, expr, AppendDesc.Params)
}

func Append(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	li, err := runtimev2.GetParamList(ctx, expr, AppendDesc.Params, 0)
	if err != nil {
		return err
	}
	elems, err := runtimev2.GetParamList(ctx, expr, AppendDesc.Params, 1)
	if err != nil {
		return err
	}
	li = append(li, elems...)
	ctx.Regs.ReturnAppend(runtimev2.V{V: li, T: ast.List})
	return nil
}
