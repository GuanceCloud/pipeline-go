package trigger

import (
	"sync"
)

type Data struct {
	Result      any               `json:"result"`
	Status      string            `json:"status"`
	DimTags     map[string]string `json:"dim_tags"`
	RelatedData map[string]any    `json:"related_data"`
}

type Trigger struct {
	vals    []Data
	rwMutex sync.RWMutex
}

func NewTr() *Trigger {
	return &Trigger{}
}

func (tr *Trigger) Trigger(result any, status string, dimTags, relatedData map[string]any) {
	tr.rwMutex.Lock()
	defer tr.rwMutex.Unlock()

	tags := map[string]string{}

	for k, v := range dimTags {
		if v, ok := v.(string); ok {
			tags[k] = v
		}
	}

	tr.vals = append(tr.vals, Data{
		Result:      result,
		Status:      status,
		DimTags:     tags,
		RelatedData: relatedData,
	})
}

func (tr *Trigger) Result() []Data {
	tr.rwMutex.RLock()
	defer tr.rwMutex.RUnlock()

	return tr.vals
}
