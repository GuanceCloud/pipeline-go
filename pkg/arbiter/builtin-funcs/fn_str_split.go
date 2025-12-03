package funcs

import (
	"strings"

	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnStrSplitDesc = runtimev2.FnDesc{
	Name: "str_split",
	Desc: "String split.",
	Params: []*runtimev2.Param{
		{
			Name: "text",
			Desc: "String to be split.",
			Typs: []ast.DType{ast.String},
		},
		{
			Name: "sep",
			Desc: "Separator to be used for splitting.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "List of strings.",
			Typs: []ast.DType{ast.List},
		},
	},
}

func FnStrSplitCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnStrSplitDesc.Params)
}

func FnStrSplit(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	text, err := runtimev2.GetParamString(ctx, funcExpr, FnStrSplitDesc.Params, 0)
	if err != nil {
		return err
	}
	sep, err := runtimev2.GetParamString(ctx, funcExpr, FnStrSplitDesc.Params, 1)
	if err != nil {
		return err
	}

	result := strings.Split(text, sep)
	sLi := make([]any, len(result))
	for i := range result {
		sLi[i] = result[i]
	}

	ctx.Regs.ReturnAppend(runtimev2.V{
		V: sLi,
		T: ast.List,
	})
	return nil
}
