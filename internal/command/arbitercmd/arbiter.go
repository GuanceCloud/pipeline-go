package arbitercmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/GuanceCloud/pipeline-go/pkg/arbiter"
	funcs "github.com/GuanceCloud/pipeline-go/pkg/arbiter/builtin-funcs"
	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/request"
	"github.com/GuanceCloud/pipeline-go/pkg/arbiter/trigger"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/spf13/cobra"
)

type config struct {
	openapiEndpoint string
	openapiKey      string
	programStr      string
	duration        string

	listFn      bool
	outFnFormat string
}

func NewCommand(stdout, stderr io.Writer) *cobra.Command {
	cfg := &config{
		openapiEndpoint: "https://openapi.guance.com",
		duration:        "15m",
	}

	root := &cobra.Command{
		Use:   "arbiter",
		Short: "Arbiter command line tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cfg, stdout, args)
		},
	}
	root.SetOut(stdout)
	root.SetErr(stderr)

	runCommand := &cobra.Command{
		Use:   "run",
		Short: "Run arbiter program",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cfg, stdout, args)
		},
	}
	runCommand.Flags().StringVarP(
		&cfg.openapiEndpoint, "guance", "e", cfg.openapiEndpoint, "GuanceCloud openapi endpoint")
	runCommand.Flags().StringVarP(
		&cfg.openapiKey, "guance-key", "k", "", "GuanceCloud openapi key")
	runCommand.Flags().StringVarP(
		&cfg.programStr, "cmd", "c", "", "program passed in as string")
	runCommand.Flags().StringVarP(
		&cfg.duration, "duration", "d", cfg.duration, "query time range, such as 1h, 15m, 60s")

	funcCommand := &cobra.Command{
		Use:   "fn",
		Short: "Arbiter built-in functions",
		Run: func(cmd *cobra.Command, args []string) {
			fn(cfg, stdout)
		},
	}
	funcCommand.Flags().BoolVarP(
		&cfg.listFn, "list", "l", false, "list functions")
	funcCommand.Flags().StringVarP(
		&cfg.outFnFormat, "output", "o", "", "output format, one of: (wide, json)")

	root.AddCommand(runCommand)
	root.AddCommand(funcCommand)
	return root
}

func Execute(args []string, stdout, stderr io.Writer) int {
	cmd := NewCommand(stdout, stderr)
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}
	return 0
}

func run(cfg *config, stdout io.Writer, args []string) error {
	tr := trigger.NewTr()
	var name, script string
	if len(args) == 1 {
		b, err := os.ReadFile(args[0])
		if err != nil {
			return err
		}
		name = args[0]
		script = string(b)
	} else {
		script = cfg.programStr
	}

	if script == "" {
		return fmt.Errorf("no program passed")
	}

	ts := time.Now().Unix()
	d, err := time.ParseDuration(cfg.duration)
	if err != nil {
		return err
	}

	timeRange := []int64{
		int64(ts-d.Milliseconds()/1e3) * 1e3,
		ts * 1e3,
	}

	runStdout := bytes.NewBuffer([]byte{})
	if err := arbiter.Run(name, script,
		arbiter.WithDQLOpenAPI(cfg.openapiEndpoint, cfg.openapiKey, timeRange),
		arbiter.WithFuncs(funcs.Funcs),
		arbiter.WithStdout(runStdout),
		arbiter.WithTrigger(tr),
		arbiter.WithHTTPClient(request.NewHTTPClient(0)),
	); err != nil {
		return err
	}

	b := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(b)
	enc.SetIndent("", "  ")
	_ = enc.Encode(tr.Result())
	fmt.Fprintf(stdout, "=== stdout:\n%s\n=== program run result:\ntrigger output:\n%s\n",
		runStdout.String(), b.String())

	return nil
}

func fn(cfg *config, stdout io.Writer) {
	if !cfg.listFn {
		return
	}

	switch cfg.outFnFormat {
	case "json":
		var fnLi []runtimev2.Desc
		for _, fn := range funcs.Funcs {
			fnLi = append(fnLi, fn.Desc.OStruct())
		}
		b := bytes.NewBuffer([]byte{})
		enc := json.NewEncoder(b)
		enc.SetIndent("", "    ")
		_ = enc.Encode(fnLi)
		fmt.Fprintln(stdout, b.String())
	case "wide":
		for _, fn := range funcs.Funcs {
			fmt.Fprintln(stdout, fn.Desc.OMarkdown("    "))
		}
	default:
		for _, fn := range funcs.Funcs {
			fmt.Fprintln(stdout, fn.Desc.String())
		}
	}
}
