// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

// Package stats used to record pl metrics
package funcs

import (
	"fmt"
	"sync"
	"testing"
)

func TestNetPolicyConcurrentReplacement(t *testing.T) {
	old := gNetFilterPolicy.Load()
	defer gNetFilterPolicy.Store(old)
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 500; j++ {
				ReplaceNetFilter(false, nil, []string{"allowed.example"})
				p := gNetFilterPolicy.Load()
				if p != nil {
					filterURL("https://allowed.example", p.disableInternal, p.cidrsWhitelist, p.hostWhitelist)
				}
			}
		}()
	}
	wg.Wait()
}

func TestNetPolicyUnsetStillValidatesURL(t *testing.T) {
	old := gNetFilterPolicy.Load()
	defer gNetFilterPolicy.Store(old)
	gNetFilterPolicy.Store(nil)
	for _, target := range []string{"file:///tmp/not-http", "://invalid"} {
		if !requestURLBlocked(target) {
			t.Fatalf("accepted %q", target)
		}
	}
	if requestURLBlocked("https://example.com") {
		t.Fatal("default policy blocked HTTP")
	}
}

func TestNetPolicyReplacesAndCopiesAllowlists(t *testing.T) {
	old := gNetFilterPolicy.Load()
	defer gNetFilterPolicy.Store(old)
	hosts := []string{"old.example"}
	ReplaceNetFilter(false, nil, hosts)
	hosts[0] = "mutated.example"
	if requestURLBlocked("https://old.example") {
		t.Fatal("caller mutated policy")
	}
	ReplaceNetFilter(false, nil, []string{"new.example"})
	if !requestURLBlocked("https://old.example") || requestURLBlocked("https://new.example") {
		t.Fatal("policy was appended instead of replaced")
	}
}

func TestSetNetFilterPreservesAppendSemantics(t *testing.T) {
	old := gNetFilterPolicy.Load()
	defer gNetFilterPolicy.Store(old)
	gNetFilterPolicy.Store(nil)
	hosts := []string{"old.example"}
	cidrs := []string{"192.0.2.0/24"}
	SetNetFilter(true, cidrs, hosts)
	hosts[0] = "changed.example"
	cidrs[0] = "198.51.100.0/24"
	SetNetFilter(false, nil, []string{"new.example"})
	p := gNetFilterPolicy.Load()
	if p.disableInternal || len(p.cidrsWhitelist) != 1 || p.cidrsWhitelist[0] != "192.0.2.0/24" {
		t.Fatalf("unexpected policy: %+v", p)
	}
	if requestURLBlocked("https://old.example") || requestURLBlocked("https://new.example") {
		t.Fatal("append lost hosts or aliased caller memory")
	}
	ReplaceNetFilter(false, nil, nil)
	if p := gNetFilterPolicy.Load(); len(p.hostWhitelist) != 0 || len(p.cidrsWhitelist) != 0 {
		t.Fatal("replacement did not clear allowlists")
	}
}

func TestSetNetFilterConcurrentAppend(t *testing.T) {
	old := gNetFilterPolicy.Load()
	defer gNetFilterPolicy.Store(old)
	gNetFilterPolicy.Store(nil)
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			SetNetFilter(false, nil, []string{fmt.Sprintf("host-%d.example", i)})
		}(i)
	}
	wg.Wait()
	p := gNetFilterPolicy.Load()
	if len(p.hostWhitelist) != 32 {
		t.Fatalf("lost concurrent entries: %d", len(p.hostWhitelist))
	}
	for i := 0; i < 32; i++ {
		if requestURLBlocked(fmt.Sprintf("https://host-%d.example", i)) {
			t.Fatalf("missing host %d", i)
		}
	}
}
