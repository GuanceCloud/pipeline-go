// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"github.com/GuanceCloud/grok"
	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtime"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
	"github.com/dgraph-io/ristretto"
)

var FnGrokDesc = runtimev2.FnDesc{
	Name: "grok",
	Desc: "Extracts data from a string using a Grok pattern. " +
		"Grok is based on regular expression syntax, and using regular (named) capture groups in a pattern is equivalent to using a pattern in a pattern. " +
		"A valid regular expression is also a valid Grok pattern.",
	Params: []*runtimev2.Param{
		{
			Name: "input",
			Desc: "The input string used to extract data.",
			Typs: []ast.DType{ast.String},
		},
		{
			Name: "pattern",
			Desc: "The pattern used to extract data.",
			Typs: []ast.DType{ast.String},
		},
		{
			Name: "extra_patterns",
			Desc: "Additional patterns for parsing patterns.",
			Typs: []ast.DType{ast.Map},
			Val:  func() any { return map[string]any{} },
		},
		{
			Name: "trim_space",
			Desc: "Whether to trim leading and trailing spaces from the parsed value.",
			Typs: []ast.DType{ast.Bool},
			Val:  func() any { return true },
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "The parsed result.",
			Typs: []ast.DType{ast.Map},
		},
		{
			Desc: "Whether the parsing was successful.",
			Typs: []ast.DType{ast.Bool},
		},
	},
}

func FnGrokCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnGrokDesc.Params)
}

var defaultPatterns = runtime.DenormalizedGlobalPatterns

func FnGrok(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	val, err := runtimev2.GetParamString(ctx, funcExpr, FnGrokDesc.Params, 0)
	if err != nil {
		return err
	}

	pattern, err := runtimev2.GetParamString(ctx, funcExpr, FnGrokDesc.Params, 1)
	if err != nil {
		return err
	}

	extraPatterns, err := runtimev2.GetParamMap(ctx, funcExpr, FnGrokDesc.Params, 2)
	if err != nil {
		return err
	}

	extra := map[string]string{}
	for k, v := range extraPatterns {
		if s, ok := v.(string); ok {
			extra[k] = s
		}
	}

	trim, err := runtimev2.GetParamBool(ctx, funcExpr, FnGrokDesc.Params, 3)
	if err != nil {
		return err
	}

	var norP PatternStorages

	if nor, _ := grok.DenormalizePatternsFromMap(extra, defaultPatterns); len(nor) > 0 {
		norP[0] = nor
		norP[1] = defaultPatterns
	} else {
		norP[0] = defaultPatterns
	}

	var grokRe *grok.GrokRegexp

	if p, err := grok.DenormalizePattern(pattern, norP); err == nil {
		if grokCache != nil {
			if v, ok := grokCache.Get(p.Denormalized()); ok {
				grokRe = v.(*grok.GrokRegexp)
			} else {
				if r, err := grok.CompilePattern2(p, norP); err == nil {
					grokRe = r
					grokCache.Set(p.Denormalized(), r, 1)
					grokCache.Wait()
				}
			}
		}
	}

	r := map[string]any{}
	if grokRe == nil {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: map[string]any(nil), T: ast.Map},
			runtimev2.V{V: false, T: ast.Bool},
		)
	} else {
		if grokRe.WithTypeInfo() {
			result, err := grokRe.RunWithTypeInfo(val, trim)
			if err != nil {
				ctx.Regs.ReturnAppend(
					runtimev2.V{V: map[string]any(nil), T: ast.Map},
					runtimev2.V{V: false, T: ast.Bool},
				)
				return nil
			}
			for i, name := range grokRe.MatchNames() {
				r[name] = result[i]
			}
		} else {
			result, err := grokRe.Run(val, trim)
			if err != nil {
				ctx.Regs.ReturnAppend(
					runtimev2.V{V: map[string]any(nil), T: ast.Map},
					runtimev2.V{V: false, T: ast.Bool},
				)
				return nil
			}
			for i, name := range grokRe.MatchNames() {
				r[name] = result[i]
			}
		}
	}

	ctx.Regs.ReturnAppend(
		runtimev2.V{V: r, T: ast.Map},
		runtimev2.V{V: true, T: ast.Bool},
	)
	return nil
}

var grokCache = func() *ristretto.Cache {
	v, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7, // number of keys to track frequency of (10M).
		MaxCost:     1e7, // maximum cost of cache (10M * ~500B = 500M).
		BufferItems: 64,  // number of keys per Get buffer.
	})

	return v
}()

type PatternStorages [2]map[string]*grok.GrokPattern

func (p PatternStorages) GetPattern(str string) (*grok.GrokPattern, bool) {
	for _, v := range p {
		if v, ok := v[str]; ok {
			return v, ok
		}
	}
	return nil, false
}

func (p PatternStorages) SetPattern(string, *grok.GrokPattern) {}
