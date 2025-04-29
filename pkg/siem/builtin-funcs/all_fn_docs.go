package funcs

import (
	"bytes"
	"encoding/json"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/GuanceCloud/pipeline-go/pkg/siem/trigger"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
)

type DocVarb struct {
	Name   string
	FnDesc *runtimev2.FnDesc
	FnExp  *FuncExample
}

// docs 类型定义
type docs []*DocVarb

// 实现 sort.Interface 接口的 Len 方法
func (d docs) Len() int {
	return len(d)
}

// 实现 sort.Interface 接口的 Less 方法
func (d docs) Less(i, j int) bool {
	return d[i].Name < d[j].Name
}

// 实现 sort.Interface 接口的 Swap 方法
func (d docs) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func GenDocs(e map[string]*FuncExample, f map[string]*runtimev2.Fn) []*DocVarb {
	var r []*DocVarb
	for name, fnD := range f {
		d := DocVarb{
			Name:   name,
			FnDesc: &fnD.Desc,
		}
		if v, ok := e[name]; ok {
			d.FnExp = v
		}

		r = append(r, &d)
	}
	sort.Sort(docs(r))
	return r
}

func GenerateDocs() (string, error) {
	var r string
	docBuf, err := os.ReadFile("./docs/FnDocs.tmpl")
	if err != nil {
		return "", err
	}

	temp := template.New("docs")
	temp = temp.Funcs(template.FuncMap{
		"signature": func(d *runtimev2.FnDesc) string {
			return d.Signature()
		},
		"typestr": func(p *runtimev2.Param) string {
			return p.TypStr()
		},
		"trigger_output": func(d []trigger.Data) string {
			b := bytes.NewBuffer([]byte{})
			enc := json.NewEncoder(b)
			enc.SetIndent("", "    ")
			_ = enc.Encode(d)
			return b.String()
		},
		"endsWithNewline": func(s string) bool {
			return strings.HasSuffix(s, "\n")
		},
	})

	if temp, err = temp.Parse(string(docBuf)); err != nil {
		return "", err
	}

	b := bytes.NewBuffer([]byte{})
	err = temp.Execute(b, GenDocs(FnExps, Funcs))
	if err != nil {
		return "", err
	}
	r = b.String()
	return r, nil
}

type FuncElem struct {
	Name      string           `json:"name"`
	Signature string           `json:"signature"`
	Desc      string           `json:"desc"`
	Params    []*FuncElemParam `json:"params"`
	Returns   []*FuncElemRet   `json:"returns"`
	Examples  []*FuncElemExp   `json:"examples"`
}

type FuncElemParam struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	Type string `json:"type"`
}

type FuncElemRet struct {
	Desc string `json:"desc"`
	Type string `json:"type"`
}

type FuncElemExp struct {
	Code       string `json:"code"`
	Stdout     string `json:"stdout"`
	TriggerOut string `json:"trigger_out"`
}

func GenerateDocs2() ([]*FuncElem, error) {
	docs := GenDocs(FnExps, Funcs)
	r := make([]*FuncElem, 0, len(docs))
	for _, d := range docs {
		e := &FuncElem{
			Name:      d.Name,
			Signature: d.FnDesc.Signature(),
			Desc:      d.FnDesc.Desc,
		}
		for _, p := range d.FnDesc.Params {
			e.Params = append(e.Params, &FuncElemParam{
				Name: p.Name,
				Desc: p.Desc,
				Type: p.TypStr(),
			})
		}
		for _, p := range d.FnDesc.Returns {
			e.Returns = append(e.Returns, &FuncElemRet{
				Desc: p.Desc,
				Type: p.TypStr(),
			})
		}
		for _, p := range d.FnExp.Progs {
			var trOut string
			if len(p.TriggerResult) > 0 {
				b := bytes.NewBuffer([]byte{})
				enc := json.NewEncoder(b)
				enc.SetIndent("", "    ")
				_ = enc.Encode(p.TriggerResult)
				trOut = b.String()
			}

			e.Examples = append(e.Examples, &FuncElemExp{
				Code:       p.Script,
				Stdout:     p.Stdout,
				TriggerOut: trOut,
			})
		}

		r = append(r, e)
	}
	return r, nil
}
