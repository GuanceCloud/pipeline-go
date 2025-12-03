package funcs

import (
	"embed"
	"fmt"

	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
	"github.com/GuanceCloud/platypus/pkg/token"
)

//go:embed docs
var FnDocs embed.FS

var Funcs = map[string]*runtimev2.Fn{
	FnAppendDesc.Name: {
		CallCheck: FnAppendCheck,
		Call:      FnAppend,
		Desc:      FnAppendDesc,
	},
	FnCallFuncDesc.Name: {
		CallCheck: FnCallFuncCheck,
		Call:      FnCallFunc,
		Desc:      FnCallFuncDesc,
	},
	FnCastDesc.Name: {
		CallCheck: FnCastCheck,
		Call:      FnCast,
		Desc:      FnCastDesc,
	},
	FnCIDRDesc.Name: {
		CallCheck: FnCIDRCheck,
		Call:      FnCIDR,
		Desc:      FnCIDRDesc,
	},
	FnDeleteDesc.Name: {
		CallCheck: FnDeleteCheck,
		Call:      FnDelete,
		Desc:      FnDeleteDesc,
	},
	FnDQLDesc.Name: {
		CallCheck: FnDQLCheck,
		Call:      FnDQL,
		Desc:      FnDQLDesc,
	},
	FnDQLSeriesGetDesc.Name: {
		CallCheck: FnDQLSeriesGetCheck,
		Call:      FnDQLSeriesGet,
		Desc:      FnDQLSeriesGetDesc,
	},
	FnDQLTimerangeGetDesc.Name: {
		CallCheck: FnDQLTimerangeGetCheck,
		Call:      FnDQLTimerangeGet,
		Desc:      FnDQLTimerangeGetDesc,
	},
	FnB64DecDesc.Name: {
		CallCheck: FnB64DecCheck,
		Call:      FnB64Dec,
		Desc:      FnB64DecDesc,
	},
	FnB64EncDesc.Name: {
		CallCheck: FnB64EncCheck,
		Call:      FnB64Enc,
		Desc:      FnB64EncDesc,
	},
	FnExitDesc.Name: {
		CallCheck: FnExitCheck,
		Call:      FnExit,
		Desc:      FnExitDesc,
	},
	FnGeoIPDesc.Name: {
		CallCheck: FnGeoIPCheck,
		Call:      FnGeoIP,
		Desc:      FnGeoIPDesc,
	},
	FnGJSONDesc.Name: {
		CallCheck: FnGJSONCheck,
		Call:      FnGJSON,
		Desc:      FnGJSONDesc,
	},
	FnGrokDesc.Name: {
		CallCheck: FnGrokCheck,
		Call:      FnGrok,
		Desc:      FnGrokDesc,
	},
	FnHashDesc.Name: {
		CallCheck: FnHashCheck,
		Call:      FnHash,
		Desc:      FnHashDesc,
	},
	FnHTTPRequestDesc.Name: {
		CallCheck: FnHTTPRequestChecking,
		Call:      FnHTTPRequest,
		Desc:      FnHTTPRequestDesc,
	},
	FnLenDesc.Name: {
		CallCheck: FnLenCheck,
		Call:      FnLen,
		Desc:      FnLenDesc,
	},
	FnLoadJSONDesc.Name: {
		CallCheck: FnLoadJSONCheck,
		Call:      FnLoadJSON,
		Desc:      FnLoadJSONDesc,
	},
	FnDumpJSONDesc.Name: {
		CallCheck: FnDumpJSONCheck,
		Call:      FnDumpJSON,
		Desc:      FnDumpJSONDesc,
	},
	FnLowercaseDesc.Name: {
		CallCheck: FnLowercaseCheck,
		Call:      FnLowercase,
		Desc:      FnLowercaseDesc,
	},
	FnMatchDesc.Name: {
		CallCheck: FnMatchCheck,
		Call:      FnMatch,
		Desc:      FnMatchDesc,
	},
	FnParseDateDesc.Name: {
		CallCheck: FnParseDateCheck,
		Call:      FnParseDate,
		Desc:      FnParseDateDesc,
	},
	FnParseDurationDesc.Name: {
		CallCheck: FnParseDurationCheck,
		Call:      FnParseDuration,
		Desc:      FnParseDurationDesc,
	},
	FnParseIntDesc.Name: {
		CallCheck: FnParseIntCheck,
		Call:      FnParseInt,
		Desc:      FnParseIntDesc,
	},
	FnFormatIntDesc.Name: {
		CallCheck: FnFormatIntChecking,
		Call:      FnFormatInt,
		Desc:      FnFormatIntDesc,
	},
	FnPrintfDesc.Name: {
		CallCheck: FnPrintfCheck,
		Call:      FnPrintf,
		Desc:      FnPrintfDesc,
	},
	FnReplaceDesc.Name: {
		CallCheck: FnReplaceCheck,
		Call:      FnReplace,
		Desc:      FnReplaceDesc,
	},
	FnSetoptDesc.Name: {
		CallCheck: FnSetoptCheck,
		Call:      FnSetopt,
		Desc:      *FnSetoptDesc,
	},
	FnSQLCoverDesc.Name: {
		CallCheck: FnSQLCoverCheck,
		Call:      FnSQLCover,
		Desc:      FnSQLCoverDesc,
	},
	FnStrFmtDesc.Name: {
		CallCheck: FnStrFmtCheck,
		Call:      FnStrFmt,
		Desc:      FnStrFmtDesc,
	},
	FnStrJoinDesc.Name: {
		CallCheck: FnStrJoinCheck,
		Call:      FnStrJoin,
		Desc:      FnStrJoinDesc,
	},
	FnStrSplitDesc.Name: {
		CallCheck: FnStrSplitCheck,
		Call:      FnStrSplit,
		Desc:      FnStrSplitDesc,
	},
	FnTimeNowDesc.Name: {
		CallCheck: FnTimeNowCheck,
		Call:      FnTimenow,
		Desc:      FnTimeNowDesc,
	},
	FnTriggerDesc.Name: {
		CallCheck: FnTriggerCheck,
		Call:      FnTrigger,
		Desc:      FnTriggerDesc,
	},
	FnTrimDesc.Name: {
		CallCheck: FnTrimCheck,
		Call:      FnTrim,
		Desc:      FnTrimDesc,
	},
	FnUppercaseDesc.Name: {
		CallCheck: FnUppercaseCheck,
		Call:      FnUppercase,
		Desc:      FnUppercaseDesc,
	},
	FnURLParseDesc.Name: {
		CallCheck: FnURLParseCheck,
		Call:      FnURLParse,
		Desc:      FnURLParseDesc,
	},
	FnURLDecodeDesc.Name: {
		CallCheck: FnURLDecodeCheck,
		Call:      FnURLDecode,
		Desc:      FnURLDecodeDesc,
	},
	FnUserAgentDesc.Name: {
		CallCheck: FnUserAgentCheck,
		Call:      FnUserAgent,
		Desc:      FnUserAgentDesc,
	},
	FnValidJSONDesc.Name: {
		CallCheck: FnValidJSONCheck,
		Call:      FnValidJSON,
		Desc:      FnValidJSONDesc,
	},
	FnValueTypeDesc.Name: {
		CallCheck: FnValueTypeCheck,
		Call:      FnValueType,
		Desc:      FnValueTypeDesc,
	},
	FnXMLDesc.Name: {
		CallCheck: FnXMLChecking,
		Call:      FnXML,
		Desc:      FnXMLDesc,
	},
}

func RunFnPTypErr(ctx *runtimev2.Task, p *runtimev2.Param, pos token.LnColPos) *errchain.PlError {
	return runtimev2.NewRunError(ctx,
		fmt.Sprintf("parameter `%s` expected to be of type `%s`", p.Name, p.TypStr()), pos)
}
