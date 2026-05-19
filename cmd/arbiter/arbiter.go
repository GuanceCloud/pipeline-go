package main

import (
	"os"

	"github.com/GuanceCloud/pipeline-go/internal/command/arbitercmd"
)

func main() {
	os.Exit(arbitercmd.Execute(os.Args[1:], os.Stdout, os.Stderr))
}
