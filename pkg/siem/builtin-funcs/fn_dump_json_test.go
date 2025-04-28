package funcs

import "testing"

func TestDumpJSON(t *testing.T) {
	cases := append([]ProgCase{}, cDumpJSON.Progs...)
	for _, tc := range cases {
		runCase(t, tc)
	}
}
