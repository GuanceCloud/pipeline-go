package funcs

import (
	"fmt"

	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnStrFmtDesc = runtimev2.FnDesc{
	Name: "strfmt",
	Params: []*runtimev2.Param{
		{
			Name: "format",
			Desc: "String format.",
			Typs: []ast.DType{ast.String},
		},
		{
			Name:     "args",
			Desc:     "Parameters to replace placeholders.",
			Typs:     ast.AllTyp(),
			Variable: true,
		},
	},
	Returns: []*runtimev2.Param{},
}

func FnStrFmtCheck(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, expr, FnStrFmtDesc.Params)
}

func FnStrFmt(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	format, err := runtimev2.GetParamString(ctx, expr, FnStrFmtDesc.Params, 0)
	if err != nil {
		return err
	}
	args, err := runtimev2.GetParamList(ctx, expr, FnStrFmtDesc.Params, 1)
	if err != nil {
		return err
	}
	newArgs := convFmtArgs(args)
	s := fmt.Sprintf(format, newArgs...)
	ctx.Regs.ReturnAppend(runtimev2.V{V: s, T: ast.String})
	return nil
}
