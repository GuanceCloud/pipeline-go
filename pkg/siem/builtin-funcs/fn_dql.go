package funcs

import (
	"fmt"

	"github.com/GuanceCloud/pipeline-go/pkg/siem/dql"
	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

// triger
// refer_url

var FnDQLDesc = runtimev2.FnDesc{
	Name: "dql",
	Desc: "Query data from the GuanceCloud using dql or promql.",
	Params: []*runtimev2.Param{
		{
			Name: "query",
			Desc: "DQL or PromQL query statements.",
			Typs: []ast.DType{ast.String},
		},
		{
			Name: "qtype",
			Desc: "Query language, One of `dql` or `promql`, default is `dql`.",
			Val:  func() any { return "dql" },
			Typs: []ast.DType{ast.String},
		},
		{
			Name: "limit",
			Desc: "Query limit.",
			Val: func() any {
				return int64(2000)
			},
			Typs: []ast.DType{ast.Int},
		},
		{
			Name: "offset",
			Desc: "Query offset.",
			Val: func() any {
				return int64(0)
			},
			Typs: []ast.DType{ast.Int},
		},
		{
			Name: "slimit",
			Desc: "Query slimit.",
			Val: func() any {
				return int64(2000)
			},
			Typs: []ast.DType{ast.Int},
		},
		{
			Name: "time_range",
			Desc: "Query timestamp range, " +
				"the default value can be modified externally by the script caller.",
			Val: func() any {
				return []any{}
			},
			Typs: []ast.DType{ast.List},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "Query response.",
			Typs: []ast.DType{ast.Map},
		},
		{
			Desc: "Query execution status",
			Typs: []ast.DType{ast.Bool},
		},
	},
}

func FnDQLCheck(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, expr, FnDQLDesc.Params)
}

func FnDQL(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	v, ok := ctx.PValue(PDQLCli)
	if !ok {
		return runtimev2.NewRunError(ctx, fmt.Sprintf(
			"missing context data named %s", PDQLCli), expr.NamePos)
	}
	dqlCli, ok := v.(dql.DQL)
	if !ok {
		return runtimev2.NewRunError(ctx, fmt.Sprintf("context data %s type is expected", PDQLCli), expr.NamePos)
	}
	if dqlCli == nil {
		return runtimev2.NewRunError(ctx, fmt.Sprintf(
			"context data %s value is nil", PDQLCli), expr.NamePos)
	}

	var r map[string]any

	p := make([]any, 6)

	for i := range p {
		var err *errchain.PlError
		p[i], err = runtimev2.GetParam(ctx, expr, FnDQLDesc.Params, i)
		if err != nil {
			return err
		}
	}

	query, err := runtimev2.GetParamString(ctx, expr, FnDQLDesc.Params, 0)
	if err != nil {
		return err
	}

	qtype, err := runtimev2.GetParamString(ctx, expr, FnDQLDesc.Params, 1)
	if err != nil {
		return err
	}

	limit, err := runtimev2.GetParamInt(ctx, expr, FnDQLDesc.Params, 2)
	if err != nil {
		return err
	}

	offset, err := runtimev2.GetParamInt(ctx, expr, FnDQLDesc.Params, 3)
	if err != nil {
		return err
	}

	slimit, err := runtimev2.GetParamInt(ctx, expr, FnDQLDesc.Params, 4)
	if err != nil {
		return err
	}

	timeRange, err := runtimev2.GetParamList(ctx, expr, FnDQLDesc.Params, 5)
	if err != nil {
		return err
	}

	r, errQ := dqlCli.Query(query, qtype, limit, offset, slimit, timeRange)
	if errQ != nil {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: (map[string]any)(nil), T: ast.Map},
			runtimev2.V{V: false, T: ast.Bool})
	} else {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: r, T: ast.Map},
			runtimev2.V{V: true, T: ast.Bool})
	}
	return nil
}
