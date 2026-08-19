package funcs

import (
	"fmt"

	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/logging"
	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnPushLogDesc = runtimev2.FnDesc{
	Name: "push_log",
	Desc: "Push one or more log records to the current workspace.",
	Params: []*runtimev2.Param{
		{
			Name: "data",
			Desc: "A structured log record or a list of records. Each record contains a non-empty fields map; source, time and tags are optional, and message is not required. If time is omitted, the push_log call time is used; time accepts a Unix timestamp in nanoseconds.",
			Typs: []ast.DType{ast.Map, ast.List},
		},
		{
			Name: "index",
			Desc: "Optional destination log index. An empty value uses the default index.",
			Typs: []ast.DType{ast.String},
			Val:  func() any { return "" },
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "Push result containing ok, accepted and error.",
			Typs: []ast.DType{ast.Map},
		},
	},
}

func FnPushLogCheck(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, expr, FnPushLogDesc.Params)
}

func FnPushLog(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	private, ok := ctx.PValue(PLogWriter)
	if !ok {
		return runtimev2.NewRunError(ctx,
			fmt.Sprintf("missing context data named %s", PLogWriter), expr.NamePos)
	}

	writer, ok := private.(logging.Writer)
	if !ok || writer == nil {
		return runtimev2.NewRunError(ctx,
			fmt.Sprintf("context data %s type is unexpected", PLogWriter), expr.NamePos)
	}

	data, err := runtimev2.GetParam(ctx, expr, FnPushLogDesc.Params, 0)
	if err != nil {
		return err
	}
	index, err := runtimev2.GetParamString(ctx, expr, FnPushLogDesc.Params, 1)
	if err != nil {
		return err
	}

	result, pushErr := writer.Push(index, data)
	if pushErr != nil {
		result = map[string]any{
			"ok":       false,
			"accepted": int64(0),
			"error":    pushErr.Error(),
		}
	}
	ctx.Regs.ReturnAppend(runtimev2.V{V: result, T: ast.Map})
	return nil
}
