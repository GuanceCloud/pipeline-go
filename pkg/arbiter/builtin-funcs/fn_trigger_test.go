package funcs

import (
	"testing"

	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/trigger"
)

func TestFnTrigger(t *testing.T) {
	cases := []ProgCase{}
	cases = append(cases, cTrigger.Progs...)

	for _, tc := range cases {
		runCase(t, tc)
	}
}

func TestFnTriggerConvertsDimensionTagsToStrings(t *testing.T) {
	runCase(t, ProgCase{
		Name: "trigger_converts_dimension_tags_to_strings",
		Script: `trigger(
    result="ok",
    status="critical",
    dimension_tags={
        "str": "value",
        "int": 12,
        "float": 1.25,
        "bool": true,
        "nil": nil
    },
    related_data={}
)`,
		TriggerResult: []trigger.Data{
			{
				Result: "ok",
				Status: "critical",
				DimensionTags: map[string]string{
					"str":   "value",
					"int":   "12",
					"float": "1.25",
					"bool":  "true",
					"nil":   "null",
				},
				RelatedData: map[string]any{},
			},
		},
	})
}
