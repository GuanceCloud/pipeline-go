package funcs

import "testing"

func TestB64Enc(t *testing.T) {
	cases := append([]ProgCase{}, cB64Enc.Progs...)
	for _, tc := range cases {
		runCase(t, tc)
	}
}

func TestB64Dec(t *testing.T) {
	cases := append([]ProgCase{}, cB64Dec.Progs...)
	for _, tc := range cases {
		runCase(t, tc)
	}
}
