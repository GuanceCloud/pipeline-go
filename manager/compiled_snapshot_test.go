// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

package manager

import (
	"github.com/GuanceCloud/cliutils/point"
	"github.com/GuanceCloud/pipeline-go/constants"
	"github.com/GuanceCloud/pipeline-go/lang"
	"github.com/GuanceCloud/pipeline-go/lang/platypus"
	"testing"
)

func TestCompiledSnapshotPreservesInstancesAndPriority(t *testing.T) {
	compile := func(ns string) *platypus.PlScript {
		t.Helper()
		scripts, errs := platypus.NewScripts(map[string]string{"test.p": `add_key(ok,true)`}, lang.WithNS(ns), lang.WithCat(point.Logging), lang.WithCache())
		if len(errs) != 0 {
			t.Fatal(errs)
		}
		p := scripts["test.p"]
		t.Cleanup(p.Cleanup)
		return p
	}
	local, remote := compile(constants.NSDefault), compile(constants.NSRemote)
	for _, input := range [][]*platypus.PlScript{{local, remote}, {remote, local}} {
		snapshot, err := NewCompiledSnapshot(NewManagerCfg(nil, nil), input)
		if err != nil {
			t.Fatal(err)
		}
		got, _ := snapshot.QueryScript(point.Logging, "test.p")
		if got != remote {
			t.Fatal("wrong namespace priority")
		}
	}
	snapshot, err := NewCompiledSnapshot(NewManagerCfg(nil, nil), []*platypus.PlScript{local})
	if err != nil {
		t.Fatal(err)
	}
	got, _ := snapshot.QueryScript(point.Logging, "test.p")
	if got != local {
		t.Fatal("script was recompiled")
	}
	for _, bad := range [][]*platypus.PlScript{{nil}, {local, local}} {
		if _, err := NewCompiledSnapshot(NewManagerCfg(nil, nil), bad); err == nil {
			t.Fatal("invalid input accepted")
		}
	}
	got, _ = snapshot.QueryScript(point.Logging, "test.p")
	if got != local {
		t.Fatal("other candidate changed existing snapshot")
	}
}
