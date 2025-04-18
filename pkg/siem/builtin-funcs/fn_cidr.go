// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"net/netip"

	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
)

var FnCIDRDesc = runtimev2.FnDesc{
	Name: "cidr",
	Desc: "Check the IP whether in CIDR block",
	Params: []*runtimev2.Param{
		{
			Name: "ip",
			Desc: "The ip address",
			Typs: []ast.DType{ast.String},
		},
		{
			Name: "mask",
			Desc: "The CIDR mask",
			Typs: []ast.DType{ast.String},
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "Whether the IP is in CIDR block",
			Typs: []ast.DType{ast.Bool},
		},
	},
}

func FnCIDRCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnCIDRDesc.Params)
}

func FnCIDR(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	ipAddr, err := runtimev2.GetParamString(ctx, funcExpr, FnCIDRDesc.Params, 0)
	if err != nil {
		return err
	}

	mask, err := runtimev2.GetParamString(ctx, funcExpr, FnCIDRDesc.Params, 1)
	if err != nil {
		return err
	}

	ok, _ := CIDRContains(ipAddr, mask)
	ctx.Regs.ReturnAppend(runtimev2.V{V: ok, T: ast.Bool})
	return nil
}

func CIDRContains(ipAddr, prefix string) (bool, error) {
	network, err := netip.ParsePrefix(prefix)
	if err != nil {
		return false, err
	}

	ip, err := netip.ParseAddr(ipAddr)
	if err != nil {
		return false, err
	}

	return network.Contains(ip), nil
}
