package arbiter

import (
	"io"

	funcs "github.com/GuanceCloud/pipeline-go/pkg/arbiter/builtin-funcs"
	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/dql"
	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/opt"
	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/request"
	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/trigger"
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

func WithDQLKodo(url, wsUUID string, timeRange []int64) Opt {
	return func(c *Config) {
		dql := dql.NewDQLKodo(url, wsUUID, timeRange)
		c.Private[funcs.PDQLCli] = dql
	}
}

func WithDQLOpenAPI(endpoint string, apikey string, timeRange []int64) Opt {
	return func(c *Config) {
		dqlCli := dql.NewDQLOpenAPI(
			endpoint, dql.OpenAPIPath, apikey, timeRange)
		c.Private[funcs.PDQLCli] = dqlCli
	}
}

func WithStdout(writer io.Writer) Opt {
	return func(c *Config) {
		c.Private[funcs.PStdout] = writer
	}
}

func WithOptions(options *opt.Option) Opt {
	return func(c *Config) {
		c.Private[funcs.POptions] = options
	}
}

func WithTrigger(tr *trigger.Trigger) Opt {
	return func(c *Config) {
		c.Private[funcs.PTrigger] = tr
	}
}

func WithHTTPClient(client *request.HTTPClient) Opt {
	return func(c *Config) {
		c.Private[funcs.PHTTPClient] = client
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
