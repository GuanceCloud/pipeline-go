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

var docVarbs = func() []*DocVarb {
	var r []*DocVarb

	m := map[string]*FuncExample{}
	for _, v := range FnExps {
		m[v.FnName] = v
	}

	for name, fnD := range Funcs {
		d := DocVarb{
			Name:   name,
			FnDesc: &fnD.Desc,
		}
		if v, ok := m[name]; ok {
			d.FnExp = v
		}

		r = append(r, &d)
	}
	sort.Sort(docs(r))
	return r
}()

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
	err = temp.Execute(b, docVarbs)
	if err != nil {
		return "", err
	}
	r = b.String()
	return r, nil
}
