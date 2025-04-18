package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/GuanceCloud/pipeline-go/pkg/siem"
	funcs "github.com/GuanceCloud/pipeline-go/pkg/siem/builtin-funcs"
	"github.com/GuanceCloud/pipeline-go/pkg/siem/trigger"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/spf13/cobra"
)

var (
	openapiEndpoint string
	openapiKey      string
	programStr      string

	listFn      bool
	outFnFormat string
)

var rootCommad = &cobra.Command{
	Use:   "siem run -e https://openapi.guance.com -k xxxxxx script.p",
	Short: "SIEM command line tool",
	RunE:  run,
}

var runCommand = &cobra.Command{
	Use:   "run",
	Short: "Run SIEM program",
	RunE:  run,
}

var funcCommand = &cobra.Command{
	Use:   "fn",
	Short: "SIEM built-in functions",
	Run:   fn,
}

func init() {
	rootCommad.AddCommand(runCommand)
	rootCommad.AddCommand(funcCommand)

	runCommand.Flags().StringVarP(
		&openapiEndpoint, "guance", "e", "https://openapi.guance.com", "GuanceCloud openapi endpoint")
	runCommand.Flags().StringVarP(
		&openapiKey, "guance-key", "k", "", "GuanceCloud openapi key")
	runCommand.Flags().StringVarP(
		&programStr, "cmd", "c", "", "program passed in as string")

	funcCommand.Flags().BoolVarP(
		&listFn, "list", "l", false, "list functions")
	funcCommand.Flags().StringVarP(
		&outFnFormat, "output", "o", "", "output format, one of: (wide, json)")
}

func main() {
	if err := rootCommad.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
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
		script = programStr
	}

	if script == "" {
		return fmt.Errorf("no program passed")
	}

	if err := siem.Run(name, script,
		siem.WithDQLOpenAPI(openapiEndpoint, openapiKey),
		siem.WithFuncs(funcs.Funcs),
		siem.WithStdout(os.Stdout),
		siem.WithTrigger(tr),
	); err != nil {
		return err
	} else {
		b := bytes.NewBuffer([]byte{})
		enc := json.NewEncoder(b)
		enc.SetIndent("", "  ")
		_ = enc.Encode(tr.Result())
		fmt.Fprintf(os.Stdout, "\n=== program run result:\ntrigger output:\n"+b.String()+"\n")
	}

	return nil
}

func fn(cmd *cobra.Command, args []string) {

	if listFn {
		switch outFnFormat {
		case "json":
			var fnLi []runtimev2.Desc
			for _, fn := range funcs.Funcs {
				fnLi = append(fnLi, fn.Desc.OStruct())
			}
			b := bytes.NewBuffer([]byte{})
			enc := json.NewEncoder(b)
			enc.SetIndent("", "    ")
			_ = enc.Encode(fnLi)
			fmt.Println(b.String())
		case "wide":
			for _, fn := range funcs.Funcs {
				fmt.Println(fn.Desc.OMarkdown("    "))
			}
		default:
			for _, fn := range funcs.Funcs {
				fmt.Println(fn.Desc.String())
			}
		}
	}
}
