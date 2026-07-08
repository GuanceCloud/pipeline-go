package trigger

import (
	"reflect"
	"testing"
)

func TestTriggerConvertsDimensionTagsToStrings(t *testing.T) {
	tr := NewTr()
	tr.Trigger(
		"result",
		"critical",
		map[string]any{
			"str":    "value",
			"int":    int64(12),
			"float":  1.25,
			"bool":   true,
			"nil":    nil,
			"list":   []any{int64(1), "a"},
			"object": map[string]any{"host": "h1"},
		},
		map[string]any{},
		"",
	)

	got := tr.Result()
	if len(got) != 1 {
		t.Fatalf("expected 1 trigger result, got %d", len(got))
	}

	wantTags := map[string]string{
		"str":    "value",
		"int":    "12",
		"float":  "1.25",
		"bool":   "true",
		"nil":    "null",
		"list":   `[1,"a"]`,
		"object": `{"host":"h1"}`,
	}
	if !reflect.DeepEqual(got[0].DimensionTags, wantTags) {
		t.Fatalf("expected dimension tags %#v, got %#v", wantTags, got[0].DimensionTags)
	}
}

func TestDimensionTagValueToStringFallsBackWhenJSONUnsupported(t *testing.T) {
	got := dimensionTagValueToString(func() {})
	if got == "" {
		t.Fatal("expected fallback string for unsupported JSON value")
	}
}
