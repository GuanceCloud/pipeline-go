package trigger

import (
	"sync"
)

type Data struct {
	Result             any               `json:"result"`
	Status             string            `json:"status"`
	DimensionTags      map[string]string `json:"dimension_tags"`
	RelatedData        map[string]any    `json:"related_data"`
	CheckWorkspaceUUID string            `json:"check_workspace_uuid"`
}

type Trigger struct {
	vals    []Data
	rwMutex sync.RWMutex
}

func NewTr() *Trigger {
	return &Trigger{}
}

func (tr *Trigger) Trigger(result any, status string, dimTags, relatedData map[string]any, check_workspace_uuid string) {
	tr.rwMutex.Lock()
	defer tr.rwMutex.Unlock()

	tags := map[string]string{}

	for k, v := range dimTags {
		if v, ok := v.(string); ok {
			tags[k] = v
		}
	}

	data := Data{
		Result:        result,
		Status:        status,
		DimensionTags: tags,
		RelatedData:   relatedData,
	}
	if check_workspace_uuid != "" {
		data.CheckWorkspaceUUID = check_workspace_uuid
	}

	tr.vals = append(tr.vals, data)
}

func (tr *Trigger) Result() []Data {
	tr.rwMutex.RLock()
	defer tr.rwMutex.RUnlock()

	return tr.vals
}
