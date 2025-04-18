package script

import (
	"github.com/GuanceCloud/platypus/pkg/engine"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
)

func Parse(name, script string, fn map[string]*runtimev2.Fn) (*runtimev2.Script, error) {
	s, err := engine.ParseV2(name, script, fn)
	if err != nil {
		return nil, err
	}
	if err := s.Check(); err != nil {
		return nil, err
	}
	return s, nil
}
