// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"github.com/GuanceCloud/pipeline-go/ptinput/ipdb"
	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnGeoIPDesc = runtimev2.FnDesc{
	Name: "geoip",
	Desc: "GeoIP",
	Params: []*runtimev2.Param{
		{
			Name: "ip",
			Desc: "IP address.",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "IP geographical information.",
			Typs: []ast.DType{ast.Map},
		},
	},
}

func FnGeoIPCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnGeoIPDesc.Params)
}

func FnGeoIP(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	ipStr, err := runtimev2.GetParamString(ctx, funcExpr, FnGeoIPDesc.Params, 0)
	if err != nil {
		return err
	}

	var db ipdb.IPdb

	if v, ok := ctx.PValue(PGeoIPDB); ok {
		if v, ok := v.(ipdb.IPdb); ok {
			db = v
		}
	}

	dict, _ := GeoIPHandle(db, ipStr)
	ctx.Regs.ReturnAppend(runtimev2.V{V: dict, T: ast.Map})
	return nil
}

func GeoIPHandle(db ipdb.IPdb, ip string) (map[string]any, error) {
	if db == nil {
		return map[string]any{}, nil
	}

	record, err := db.Geo(ip)
	if err != nil {
		return nil, err
	}

	res := map[string]any{}

	if record != nil {
		res["city"] = record.City
		res["province"] = record.Region
		res["country"] = record.Country
		res["isp"] = record.Isp
	} else {
		res["city"] = ""
		res["province"] = ""
		res["country"] = ""
		res["isp"] = "unknown"
	}
	return res, nil
}
