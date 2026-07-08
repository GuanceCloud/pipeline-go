package funcs

import (
	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnDQLSeriesGetDesc = runtimev2.FnDesc{
	Name: "dql_series_get",
	Desc: "get series data",
	Params: []*runtimev2.Param{
		{
			Name: "series",
			Desc: "dql query result",
			Typs: []ast.DType{ast.Map},
		},
		{
			Name: "name",
			Desc: "column or tag name",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "specified column or tag value for the series",
			Typs: []ast.DType{ast.List},
		},
	},
}

var FnDQLSeriesFirstDesc = runtimev2.FnDesc{
	Name: "dql_series_first",
	Desc: "get the first column or tag value from series data",
	Params: []*runtimev2.Param{
		{
			Name: "series",
			Desc: "dql query result",
			Typs: []ast.DType{ast.Map},
		},
		{
			Name: "name",
			Desc: "column or tag name",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "first specified column or tag value for the series",
			Typs: []ast.DType{ast.Bool, ast.Int, ast.Float, ast.String, ast.List, ast.Map, ast.Nil},
		},
	},
}

func FnDQLSeriesGetCheck(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, expr, FnDQLSeriesGetDesc.Params)
}

func FnDQLSeriesFirstCheck(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, expr, FnDQLSeriesFirstDesc.Params)
}

func FnDQLSeriesGet(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	dqlResult, err := runtimev2.GetParamMap(ctx, expr, FnDQLSeriesGetDesc.Params, 0)
	if err != nil {
		return err
	}
	name, err := runtimev2.GetParamString(ctx, expr, FnDQLSeriesGetDesc.Params, 1)
	if err != nil {
		return err
	}

	ctx.Regs.ReturnAppend(runtimev2.V{T: ast.List, V: dqlSeriesValues(dqlResult, name)})

	return nil
}

func FnDQLSeriesFirst(ctx *runtimev2.Task, expr *ast.CallExpr) *errchain.PlError {
	dqlResult, err := runtimev2.GetParamMap(ctx, expr, FnDQLSeriesFirstDesc.Params, 0)
	if err != nil {
		return err
	}
	name, err := runtimev2.GetParamString(ctx, expr, FnDQLSeriesFirstDesc.Params, 1)
	if err != nil {
		return err
	}

	v := dqlSeriesFirstValue(dqlResult, name)
	v, typ := ast.DectDataType(v)
	if typ == ast.Invalid {
		v, typ = nil, ast.Nil
	}
	ctx.Regs.ReturnAppend(runtimev2.V{T: typ, V: v})

	return nil
}

func dqlSeriesValues(dqlResult map[string]any, name string) []any {
	vecVec := []any{}

	tagOrCol := -1

	series, ok := dqlResult["series"]
	_ = ok

	for _, v := range getList(series) {
		s0 := getList(v)

		var vec []any
		for _, elem := range s0 {
			elemMap := getMap(elem)
			switch tagOrCol {
			case -1:
				if v, ok := elemMap["columns"]; ok {
					colMap := getMap(v)
					if v, ok := colMap[name]; ok {
						vec = append(vec, v)
						tagOrCol = 0
						continue
					}
				}
				if v, ok := elemMap["tags"]; ok {
					tagMap := getMap(v)
					if v, ok := tagMap[name]; ok {
						vec = append(vec, v)
						tagOrCol = 1
						continue
					}
				}
			case 0:
				if v, ok := elemMap["columns"]; ok {
					colMap := getMap(v)
					if v, ok := colMap[name]; ok {
						vec = append(vec, v)
						continue
					}
				}

			case 1:
				if v, ok := elemMap["tags"]; ok {
					tagMap := getMap(v)
					if v, ok := tagMap[name]; ok {
						vec = append(vec, v)
						continue
					}
				}
			}
			vec = append(vec, nil)
		}
		vecVec = append(vecVec, vec)
	}

	return vecVec
}

func dqlSeriesFirstValue(dqlResult map[string]any, name string) any {
	series := dqlResult["series"]

	for _, v := range getList(series) {
		s0 := getList(v)
		if len(s0) == 0 {
			continue
		}

		elemMap := getMap(s0[0])
		if v, ok := elemMap["columns"]; ok {
			colMap := getMap(v)
			if v, ok := colMap[name]; ok {
				return v
			}
		}
		if v, ok := elemMap["tags"]; ok {
			tagMap := getMap(v)
			if v, ok := tagMap[name]; ok {
				return v
			}
		}
		return nil
	}
	return nil
}

func getMap(v any) map[string]any {
	if v, ok := v.(map[string]any); ok {
		return v
	}
	return nil
}

func getList(v any) []any {
	if v, ok := v.([]any); ok {
		return v
	}
	return nil
}
