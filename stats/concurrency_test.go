// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

// Package stats used to record pl metrics
package stats

import (
	"sync"
	"testing"
	"time"
)

type reentrantStats struct{ Stats }

func (s *reentrantStats) WriteMetric(map[string]string, float64, float64, float64, time.Duration) {
	SetStats(s)
}

func TestConcurrentStatsReplacementAndReentry(t *testing.T) {
	old := currentStats()
	defer SetStats(old)
	recorder := &reentrantStats{}
	SetStats(recorder)
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				SetStats(recorder)
				WriteMetric(nil, 1, 0, 0, 0)
				SetStats(nil)
			}
		}()
	}
	wg.Wait()
}
