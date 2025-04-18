package trigger

import "sync"

type Data struct {
	Result    any            `json:"result"`
	Level     string         `json:"level"`
	DimTags   map[string]any `json:"dim_data"`
	ExtraData map[string]any `json:"extra_data"`
}

type Trigger struct {
	vals    []Data
	rwMutex sync.RWMutex
}

func NewTr() *Trigger {
	return &Trigger{}
}

func (tr *Trigger) Trigger(result any, level string, dimTags, extraData map[string]any) {
	tr.rwMutex.Lock()
	defer tr.rwMutex.Unlock()

	tr.vals = append(tr.vals, Data{
		Result:    result,
		Level:     level,
		DimTags:   dimTags,
		ExtraData: extraData,
	})
}

func (tr *Trigger) Result() []Data {
	tr.rwMutex.RLock()
	defer tr.rwMutex.RUnlock()

	return tr.vals
}
