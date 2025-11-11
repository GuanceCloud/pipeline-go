package funcs

import (
	"fmt"

	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/opt"
	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnSetoptDesc = &runtimev2.FnDesc{
	Name: "setopt",
	Desc: "Set option.",
	Params: []*runtimev2.Param{
		{
			Name: "trigger_keepalive",
			Desc: "Set the trigger keepalive time (seconds).",
			Typs: []ast.DType{ast.Int},
		},
	},
	Returns: []*runtimev2.Param{},
}

func FnSetoptCheck(ctx *runtimev2.Task, funcexpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcexpr, FnSetoptDesc.Params)
}

func FnSetopt(ctx *runtimev2.Task, funcexpr *ast.CallExpr) *errchain.PlError {
	v, ok := ctx.PValue(POptions)
	if !ok {
		return runtimev2.NewRunError(ctx, fmt.Sprintf(
			"missing context data named %s", POptions), funcexpr.NamePos)
	}
	options, ok := v.(*opt.Option)
	if !ok {
		return runtimev2.NewRunError(ctx, fmt.Sprintf(
			"context data %s type is unexpected", POptions), funcexpr.NamePos)
	}

	triggerKeepalive, err := runtimev2.GetParamInt(ctx, funcexpr, FnSetoptDesc.Params, 0)
	if err != nil {
		return err
	}

	options.SetTriggerKeepalive(int(triggerKeepalive))
	return nil
}
