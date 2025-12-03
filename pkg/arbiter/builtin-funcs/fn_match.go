// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"regexp"

	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
	"github.com/dgraph-io/ristretto"
)

var FnMatchDesc = runtimev2.FnDesc{
	Name: "match",
	Desc: "Regular expression matching.",
	Params: []*runtimev2.Param{
		{
			Name: "val",
			Desc: "The string to match.",
			Typs: []ast.DType{ast.String},
		},
		{
			Name: "pattern",
			Desc: "Regular expression pattern.",
			Typs: []ast.DType{ast.String},
		},
		{
			Name: "n",
			Desc: "The number of matches to return. Defaults to 1, -1 for all matches.",
			Typs: []ast.DType{ast.Int},
			Val:  func() any { return 1 },
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "Returns the matched value.",
			Typs: []ast.DType{ast.List},
		},
		{
			Desc: "Returns true if the regular expression matches.",
			Typs: []ast.DType{ast.Bool},
		},
	},
}

func FnMatchCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnMatchDesc.Params)
}

type RegexpCache struct {
	c *ristretto.Cache
}

func (rc *RegexpCache) Get(pattern string) (*regexp.Regexp, error) {
	var re *regexp.Regexp
	if rc.c != nil {
		if r, ok := rc.c.Get(pattern); ok && r != nil {
			re = r.(*regexp.Regexp)
			return re, nil
		} else {
			re, err := regexp.Compile(pattern)
			if err != nil {
				return nil, err
			}
			rc.c.Set(pattern, re, 1)
			rc.c.Wait()
			return re, nil
		}
	} else {
		return regexp.Compile(pattern)
	}
}

var regexpCache = func() *RegexpCache {
	v, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7, // number of keys to track frequency of (10M).
		MaxCost:     1e7, // maximum cost of cache (10M * ~260B = 260M).
		BufferItems: 64,  // number of keys per Get buffer.
	})

	return &RegexpCache{c: v}
}()

func FnMatch(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	val, err := runtimev2.GetParamString(ctx, funcExpr, FnMatchDesc.Params, 0)
	if err != nil {
		return err
	}
	pattern, err := runtimev2.GetParamString(ctx, funcExpr, FnMatchDesc.Params, 1)
	if err != nil {
		return err
	}

	n, err := runtimev2.GetParamInt(ctx, funcExpr, FnMatchDesc.Params, 2)
	if err != nil {
		return err
	}

	var re *regexp.Regexp
	if r, err := regexpCache.Get(pattern); err != nil {
		return runtimev2.NewRunError(ctx, err.Error(), funcExpr.NamePos)
	} else {
		re = r
	}

	var result []any
	if r := re.FindAllStringSubmatch(val, int(n)); len(r) > 0 {
		for i := range r {
			rLi := make([]any, len(r[i]))
			for j := range r[i] {
				rLi[j] = r[i][j]
			}
			result = append(result, rLi)
		}
	}

	ctx.Regs.ReturnAppend(
		runtimev2.V{
			V: result,
			T: ast.List,
		},
		runtimev2.V{
			V: len(result) > 0,
			T: ast.Bool,
		},
	)
	return nil
}
