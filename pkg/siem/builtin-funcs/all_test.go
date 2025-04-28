package funcs

import (
	"testing"
)

func TestFnAll(t *testing.T) {
	lenFns := len(Funcs)
	lenExps := len(FnExps)
	switch {
	case lenFns > lenExps:
		for k := range Funcs {
			if _, ok := FnExps[k]; !ok {
				t.Errorf("func %s not found in FnExps", k)
			}
		}
	case lenFns < lenExps:
		for k := range FnExps {
			if _, ok := Funcs[k]; !ok {
				t.Errorf("func %s not found in Funcs", k)
			}
		}
	}
}
