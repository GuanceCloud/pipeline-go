package funcs

import (
	"testing"
)

func TestAppend(t *testing.T) {
	cases := append([]ProgCase{}, cAppend.Progs...)
	for _, tc := range cases {
		runCase(t, tc)
	}
}
