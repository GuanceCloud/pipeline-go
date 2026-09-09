// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

// Package stats used to record pl metrics
package funcs

import (
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
				SetNetFilter(false, nil, []string{"allowed.example"})
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
	SetNetFilter(false, nil, hosts)
	hosts[0] = "mutated.example"
	if requestURLBlocked("https://old.example") {
		t.Fatal("caller mutated policy")
	}
	SetNetFilter(false, nil, []string{"new.example"})
	if !requestURLBlocked("https://old.example") || requestURLBlocked("https://new.example") {
		t.Fatal("policy was appended instead of replaced")
	}
}
