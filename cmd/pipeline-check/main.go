// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2026-present Guance, Inc.

package main

import (
	"os"

	"github.com/GuanceCloud/pipeline-go/internal/command/pipelinecheck"
)

func main() {
	os.Exit(pipelinecheck.Run(os.Args[1:], os.Stdout, os.Stderr))
}
