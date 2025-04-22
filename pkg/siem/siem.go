package siem

import (
	"io"

	funcs "github.com/GuanceCloud/pipeline-go/pkg/siem/builtin-funcs"
	"github.com/GuanceCloud/pipeline-go/pkg/siem/dql"
	"github.com/GuanceCloud/pipeline-go/pkg/siem/trigger"
	"github.com/GuanceCloud/platypus/pkg/engine"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
)

type Config struct {
	Fn      map[string]*runtimev2.Fn
	Private map[runtimev2.TaskP]any
	Signal  runtimev2.Signal
}

type Opt func(*Config)

func WithFuncs(fn map[string]*runtimev2.Fn) Opt {
	return func(c *Config) {
		c.Fn = fn
	}
}

func WithDQLKodo(endpoint, path, wsUUID string) Opt {
	return func(c *Config) {
		dql := dql.NewDQLKodo(endpoint, path, wsUUID)
		c.Private[funcs.PDQLCli] = dql
	}
}

func WithDQLOpenAPI(endpoint string, apikey string) Opt {
	return func(c *Config) {
		dqlCli := dql.NewDQLOpenAPI(
			endpoint, dql.OpenAPIPath, apikey)
		c.Private[funcs.PDQLCli] = dqlCli
	}
}

func WithStdout(writer io.Writer) Opt {
	return func(c *Config) {
		c.Private[funcs.PStdout] = writer
	}
}

func WithTrigger(tr *trigger.Trigger) Opt {
	return func(c *Config) {
		c.Private[funcs.PTrigger] = tr
	}
}

func Run(name, script string, opt ...Opt) error {
	cfg := &Config{
		Private: map[runtimev2.TaskP]any{},
	}
	for _, o := range opt {
		o(cfg)
	}
	s, err := engine.ParseV2(name, script, cfg.Fn)
	if err != nil {
		return err
	}

	if err := s.Run(
		cfg.Signal,
		runtimev2.WithPrivate(cfg.Private),
	); err != nil {
		return err
	}
	return nil
}
