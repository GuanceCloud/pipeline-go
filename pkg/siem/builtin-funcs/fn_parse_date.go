// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package funcs

import (
	"fmt"
	"regexp"
	"time"

	"github.com/GuanceCloud/platypus/pkg/ast"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/GuanceCloud/platypus/pkg/errchain"
	"github.com/araddon/dateparse"
)

var FnParseDateDesc = runtimev2.FnDesc{
	Name: "parse_date",
	Desc: "Parses a date string to a nanoseconds timestamp, support multiple date formats. " +
		"If the date string not include timezone and no timezone is provided, the local timezone is used.",
	Params: []*runtimev2.Param{
		{
			Name: "date",
			Desc: "The key to use for parsing.",
			Typs: []ast.DType{ast.String},
		},
		{
			Name: "timezone",
			Desc: "The timezone to use for parsing. If ",
			Typs: []ast.DType{ast.String},
			Val:  func() any { return "" },
		},
	},
	Returns: []*runtimev2.Param{
		{
			Desc: "The parsed timestamp in nanoseconds.",
			Typs: []ast.DType{ast.Int},
		},
		{
			Desc: "Whether the parsing was successful.",
			Typs: []ast.DType{ast.Bool},
		},
	},
}

func FnParseDateCheck(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	return runtimev2.CheckPassParam(ctx, funcExpr, FnParseDateDesc.Params)
}

func FnParseDate(ctx *runtimev2.Task, funcExpr *ast.CallExpr) *errchain.PlError {
	value, err := runtimev2.GetParamString(ctx, funcExpr, FnParseDateDesc.Params, 0)
	if err != nil {
		return err
	}

	tz, err := runtimev2.GetParamString(ctx, funcExpr, FnParseDateDesc.Params, 1)
	if err != nil {
		return err
	}

	if nanots, err := TimestampHandle(value, tz); err != nil {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: int64(0), T: ast.Int},
			runtimev2.V{V: false, T: ast.Bool},
		)
	} else {
		ctx.Regs.ReturnAppend(
			runtimev2.V{V: nanots, T: ast.Int},
			runtimev2.V{V: true, T: ast.Bool},
		)
	}
	return nil
}

func TimestampHandle(value, tz string) (int64, error) {
	var t time.Time
	var err error
	timezone := time.Local

	if tz != "" {
		// parse timezone as +x or -x
		if tz[0] == '+' || tz[0] == '-' {
			if _, has := timezoneList[tz]; !has {
				return 0, fmt.Errorf("fail to parse timezone %s", tz)
			}
			tz = timezoneList[tz]
		}
		if timezone, err = time.LoadLocation(tz); err != nil {
			return 0, err
		}
	}

	// pattern match first
	unixTime, err := parseDatePattern(value, timezone)

	if unixTime > 0 && err == nil {
		return unixTime, nil
	}

	if t, err = dateparse.ParseIn(value, timezone); err != nil {
		return 0, err
	}

	// l.Debugf("parse `%s' -> %v(nano: %d)", value, t, t.UnixNano())

	return t.UnixNano(), nil
}

func parseDatePattern(value string, loc *time.Location) (int64, error) {
	valueCpy := value
	for _, p := range datePattern {
		if p.defaultYear {
			ty := time.Now()
			year := ty.Year()
			value = fmt.Sprintf("%s %d", value, year)
		} else {
			value = valueCpy
		}

		// 默认定义的规则能匹配，不匹配的则由 dataparse 处理
		if tm, err := time.ParseInLocation(p.goFmt, value, loc); err != nil {
			continue
		} else {
			unixTime := tm.UnixNano()
			return unixTime, nil
		}
	}
	return 0, fmt.Errorf("no match")
}

var timezoneList = map[string]string{
	"-11":    "Pacific/Midway",
	"-10":    "Pacific/Honolulu",
	"-9:30":  "Pacific/Marquesas",
	"-9":     "America/Anchorage",
	"-8":     "America/Los_Angeles",
	"-7":     "America/Phoenix",
	"-6":     "America/Chicago",
	"-5":     "America/New_York",
	"-4":     "America/Santiago",
	"-3:30":  "America/St_Johns",
	"-3":     "America/Sao_Paulo",
	"-2":     "America/Noronha",
	"-1":     "America/Scoresbysund",
	"+0":     "Europe/London",
	"+1":     "Europe/Vatican",
	"+2":     "Europe/Kiev",
	"+3":     "Europe/Moscow",
	"+3:30":  "Asia/Tehran",
	"+4":     "Asia/Dubai",
	"+4:30":  "Asia/Kabul",
	"+5":     "Asia/Samarkand",
	"+5:30":  "Asia/Kolkata",
	"+5:45":  "Asia/Kathmandu",
	"+6":     "Asia/Almaty",
	"+6:30":  "Asia/Yangon",
	"+7":     "Asia/Jakarta",
	"+8":     "Asia/Shanghai",
	"+8:45":  "Australia/Eucla",
	"+9":     "Asia/Tokyo",
	"+9:30":  "Australia/Darwin",
	"+10":    "Australia/Sydney",
	"+10:30": "Australia/Lord_Howe",
	"+11":    "Pacific/Guadalcanal",
	"+12":    "Pacific/Auckland",
	"+12:45": "Pacific/Chatham",
	"+13":    "Pacific/Apia",
	"+14":    "Pacific/Kiritimati",

	"CST": "Asia/Shanghai",
	"UTC": "Europe/London",
	// TODO: add more...
}

var datePattern = func() []struct {
	desc        string
	pattern     *regexp.Regexp
	goFmt       string
	defaultYear bool
} {
	dataPatternSource := []struct {
		desc        string
		pattern     string
		goFmt       string
		defaultYear bool
	}{
		{
			desc:    "nginx log datetime, 02/Jan/2006:15:04:05 -0700",
			pattern: `\d{2}/\w+/\d{4}:\d{2}:\d{2}:\d{2} \+\d{4}`,
			goFmt:   "02/Jan/2006:15:04:05 -0700",
		},
		{
			desc:    "redis log datetime, 14 May 2019 19:11:40.164",
			pattern: `\d{2} \w+ \d{4} \d{2}:\d{2}:\d{2}.\d{3}`,
			goFmt:   "02 Jan 2006 15:04:05.000",
		},
		{
			desc:        "redis log datetime, 14 May 19:11:40.164",
			pattern:     `\d{2} \w+ \d{2}:\d{2}:\d{2}.\d{3}`,
			goFmt:       "02 Jan 15:04:05.000 2006",
			defaultYear: true,
		},
		{
			desc:    "mysql, 171113 14:14:20",
			pattern: `\d{6} \d{2}:\d{2}:\d{2}`,
			goFmt:   "060102 15:04:05",
		},

		{
			desc:    "gin, 2021/02/27 - 14:14:20",
			pattern: `\d{4}/\d{2}/\d{2} - \d{2}:\d{2}:\d{2}`,
			goFmt:   "2006/01/02 - 15:04:05",
		},
		{
			desc:    "apache,  Tue May 18 06:25:05.176170 2021",
			pattern: `\w+ \w+ \d{2} \d{2}:\d{2}:\d{2}.\d{6} \d{4}`,
			goFmt:   "Mon Jan 2 15:04:05.000000 2006",
		},
		{
			desc:    "postgresql, 2021-05-27 06:54:14.760 UTC",
			pattern: `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3} UTC`,
			goFmt:   "2006-01-02 15:04:05.000 UTC",
		},
	}

	dst := []struct {
		desc        string
		pattern     *regexp.Regexp
		goFmt       string
		defaultYear bool
	}{}

	for _, p := range dataPatternSource {
		if c, err := regexp.Compile(p.pattern); err == nil {
			dst = append(dst, struct {
				desc        string
				pattern     *regexp.Regexp
				goFmt       string
				defaultYear bool
			}{
				desc:        p.desc,
				pattern:     c,
				goFmt:       p.goFmt,
				defaultYear: p.defaultYear,
			})
		}
	}
	return dst
}()
