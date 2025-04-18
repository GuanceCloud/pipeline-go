package funcs

import (
	"testing"
)

func TestFnTrigger(t *testing.T) {
	cases := []ProgCase{}
	cases = append(cases, cTrigger.Progs...)

	for _, tc := range cases {
		runCase(t, tc)
	}
}
