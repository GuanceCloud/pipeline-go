package funcs

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnPrintfDesc = runtimev2.FnDesc{
	Name: "printf",
	Desc: "Output formatted strings to the standard output device.",
	Params: []*runtimev2.Param{
		{
			Name: "format",
			Desc: "String format.",
			Typs: []ast.DType{ast.String},
		},
		{
			Name:     "args",
			Desc:     "Argument list, corresponding to the format specifiers in the format string.",
			Variable: true,
			Typs:     []ast.DType{ast.String, ast.Bool, ast.Int, ast.Float, ast.List, ast.Map},
		},
	},
}

func FnPrintfCheck(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, expr, FnPrintfDesc.Params)
}

func FnPrintf(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	var stdout io.Writer
	if v, ok := ctx.PValue(PStdout); !ok {
		return nil
	} else {
		if v, ok := v.(io.Writer); ok {
			stdout = v
		} else {
			return nil
		}
	}

	format, err := runtimev2.GetParamString(ctx, expr, FnPrintfDesc.Params, 0)
	if err != nil {
		return err
	}

	args, err := runtimev2.GetParamList(ctx, expr, FnPrintfDesc.Params, 1)
	if err != nil {
		return err
	}

	newArgs := convFmtArgs(args)
	s := fmt.Sprintf(format, newArgs...)
	_, _ = stdout.Write([]byte(s))

	return nil
}

func convFmtArgs(args []any) []any {
	newArgs := make([]any, len(args))
	for i := range args {
		r := args[i]
		switch v := args[i].(type) {
		case map[string]any:
			if val, err := json.Marshal(v); err == nil {
				r = string(val)
			}
		case []any:
			if val, err := json.Marshal(v); err == nil {
				r = string(val)
			}
		}
		newArgs[i] = r
	}
	return newArgs
}
