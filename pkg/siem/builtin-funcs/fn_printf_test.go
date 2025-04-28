package funcs

import "testing"

func TestPrintf(t *testing.T) {
	cases := append([]ProgCase{}, cPrintf.Progs...)
	for _, tc := range cases {
		runCase(t, tc)
	}
}
