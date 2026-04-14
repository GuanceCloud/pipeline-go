package manager

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/GuanceCloud/cliutils/point"
)

type wildcardBenchmarkCase struct {
	name           string
	rel            map[string]string
	exactSource    string
	wildcardSource string
}

func BenchmarkRelationQueryWildcard(b *testing.B) {
	for _, bc := range wildcardBenchmarkCases() {
		b.Run(bc.name, func(b *testing.B) {
			benchmarkWildcardDataset(b, bc)
		})
	}
}

func BenchmarkRelationQueryWildcardParallel(b *testing.B) {
	for _, bc := range wildcardBenchmarkCases() {
		b.Run(bc.name, func(b *testing.B) {
			benchmarkWildcardDatasetParallel(b, bc)
		})
	}
}

func BenchmarkRelationQueryWildcardParallelWithUpdates(b *testing.B) {
	for _, bc := range wildcardBenchmarkCases() {
		b.Run(bc.name, func(b *testing.B) {
			benchmarkWildcardDatasetParallelWithUpdates(b, bc)
		})
	}
}

func benchmarkWildcardDataset(b *testing.B, bc wildcardBenchmarkCase) {
	rl := NewPipelineRelation()
	rl.UpdateRelation(0, map[point.Category]map[string]string{
		point.Logging: bc.rel,
	})

	b.Run("exact", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = rl.Query(point.Logging, bc.exactSource)
		}
	})

	b.Run("wildcard_prod", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = rl.Query(point.Logging, bc.wildcardSource)
		}
	})
}

func benchmarkWildcardDatasetParallel(b *testing.B, bc wildcardBenchmarkCase) {
	rl := NewPipelineRelation()
	rl.UpdateRelation(0, map[point.Category]map[string]string{
		point.Logging: bc.rel,
	})

	b.Run("wildcard_prod", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = rl.Query(point.Logging, bc.wildcardSource)
			}
		})
	})
}

func benchmarkWildcardDatasetParallelWithUpdates(b *testing.B, bc wildcardBenchmarkCase) {
	rl := NewPipelineRelation()
	stateA := map[point.Category]map[string]string{
		point.Logging: cloneRelationMap(bc.rel, "a"),
	}
	stateB := map[point.Category]map[string]string{
		point.Logging: cloneRelationMap(bc.rel, "b"),
	}
	rl.UpdateRelation(1, stateA)

	var (
		stop atomic.Bool
		wg   sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		version := int64(2)
		for !stop.Load() {
			rl.UpdateRelation(version, stateB)
			version++
			if stop.Load() {
				return
			}
			rl.UpdateRelation(version, stateA)
			version++
		}
	}()

	defer func() {
		stop.Store(true)
		wg.Wait()
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = rl.Query(point.Logging, bc.wildcardSource)
		}
	})
}

func wildcardBenchmarkCases() []wildcardBenchmarkCase {
	return []wildcardBenchmarkCase{
		buildPrefixHeavyBenchmarkCase(),
		buildGenericHeavyBenchmarkCase(),
		buildOverlapHeavyBenchmarkCase(),
	}
}

func buildPrefixHeavyBenchmarkCase() wildcardBenchmarkCase {
	rel := map[string]string{
		"svc-precise": "exact.p",
		"svc-prod-*":  "prod-star.p",
		"svc-prod-??": "prod-two.p",
		"svc-*-api":   "api-star.p",
	}

	for i := 0; i < 512; i++ {
		rel[fmt.Sprintf("other-%03d-*", i)] = "other.p"
	}

	return wildcardBenchmarkCase{
		name:           "prefix_heavy",
		rel:            rel,
		exactSource:    "svc-precise",
		wildcardSource: "svc-prod-ab",
	}
}

func buildGenericHeavyBenchmarkCase() wildcardBenchmarkCase {
	rel := map[string]string{
		"svc-precise": "exact.p",
		"*prod*api*":  "prod-api.p",
		"*ab*cd*":     "ab-cd.p",
		"*xy?z*":      "xyz.p",
	}

	for i := 0; i < 512; i++ {
		rel[fmt.Sprintf("*mid%03d*", i)] = "other.p"
	}

	return wildcardBenchmarkCase{
		name:           "generic_heavy",
		rel:            rel,
		exactSource:    "svc-precise",
		wildcardSource: "aa-prod-zz-api-bb",
	}
}

func buildOverlapHeavyBenchmarkCase() wildcardBenchmarkCase {
	rel := map[string]string{
		"svc-precise":       "exact.p",
		"svc-prod-*":        "prod-star.p",
		"svc-prod-????????": "prod-8.p",
	}

	for i := 1; i <= 64; i++ {
		rel["svc-prod-"+strings.Repeat("?", i)] = fmt.Sprintf("len-%02d.p", i)
	}

	for i := 1; i <= 64; i++ {
		rel["svc-prod-"+strings.Repeat("?", i)+"*"] = fmt.Sprintf("len-star-%02d.p", i)
	}

	for i := 0; i < 256; i++ {
		rel[fmt.Sprintf("svc-prod-%s*", strings.Repeat("a", i%16+1))] = "prefix-overlap.p"
	}

	return wildcardBenchmarkCase{
		name:           "overlap_heavy",
		rel:            rel,
		exactSource:    "svc-precise",
		wildcardSource: "svc-prod-abcdefgh",
	}
}

func cloneRelationMap(rel map[string]string, suffix string) map[string]string {
	cloned := make(map[string]string, len(rel))
	for pattern, name := range rel {
		if name == "" {
			cloned[pattern] = name
			continue
		}
		cloned[pattern] = name + "-" + suffix
	}
	return cloned
}
