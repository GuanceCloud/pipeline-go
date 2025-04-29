package funcs

import "testing"

func TestStrJoin(t *testing.T) {
	cases := append([]ProgCase{}, cStrJoin.Progs...)
	for _, tc := range cases {
		runCase(t, tc)
	}
}
