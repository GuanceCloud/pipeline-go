// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2026-present Guance, Inc.

package main

import (
	"fmt"
	"os"

	"github.com/GuanceCloud/pipeline-go/internal/command/arbitercheck"
	"github.com/GuanceCloud/pipeline-go/internal/command/arbitercmd"
	"github.com/GuanceCloud/pipeline-go/internal/command/pipelinecheck"
	"github.com/spf13/cobra"
)

func main() {
	root := newRootCommand()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "pipeline-go",
		Short: "GuanceCloud pipeline and arbiter toolkit",
	}
	root.SetOut(os.Stdout)
	root.SetErr(os.Stderr)

	root.AddCommand(newPipelineCommand())
	root.AddCommand(newArbiterCommand())
	root.AddCommand(newPassthroughCommand(
		"pipeline-check",
		"Validate pipeline scripts",
		func(args []string) int {
			return pipelinecheck.Run(args, os.Stdout, os.Stderr)
		},
	))
	root.AddCommand(newPassthroughCommand(
		"arbiter-check",
		"Validate arbiter scripts",
		func(args []string) int {
			return arbitercheck.Run(args, os.Stdout, os.Stderr)
		},
	))

	return root
}

func newPipelineCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pipeline",
		Short: "Pipeline tools",
	}
	cmd.AddCommand(newPassthroughCommand(
		"check",
		"Validate pipeline scripts",
		func(args []string) int {
			return pipelinecheck.Run(args, os.Stdout, os.Stderr)
		},
	))
	return cmd
}

func newArbiterCommand() *cobra.Command {
	cmd := arbitercmd.NewCommand(os.Stdout, os.Stderr)
	checkCmd := newPassthroughCommand(
		"check",
		"Validate arbiter scripts",
		func(args []string) int {
			return arbitercheck.Run(args, os.Stdout, os.Stderr)
		},
	)
	cmd.AddCommand(checkCmd)
	return cmd
}

func newPassthroughCommand(use, short string, run func(args []string) int) *cobra.Command {
	return &cobra.Command{
		Use:                use,
		Short:              short,
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			os.Exit(run(args))
		},
	}
}
