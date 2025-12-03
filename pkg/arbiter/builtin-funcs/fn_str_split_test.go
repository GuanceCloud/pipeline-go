package funcs

import "testing"

func TestStrSplit(t *testing.T) {
	cases := append([]ProgCase{}, cStrSplit.Progs...)
	for _, tc := range cases {
		runCase(t, tc)
	}
}
