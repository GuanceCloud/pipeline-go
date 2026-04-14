// Unless explicitly stated otherwise all files in this repository are licensed
// under the MIT License.
// This product includes software developed at Guance Cloud (https://www.guance.com/).
// Copyright 2021-present Guance, Inc.

// Package manager for managing pipeline scripts
package manager

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/GuanceCloud/cliutils/point"
	"github.com/stretchr/testify/assert"
)

func TestRelation(t *testing.T) {
	rl := NewPipelineRelation()

	rl.UpdateRelation(0, map[point.Category]map[string]string{
		point.Logging: {
			"abc": "a1.p",
		},
	})
	p, ok := rl.Query(point.Logging, "abc")
	assert.True(t, ok)
	assert.Equal(t, "a1.p", p)

	p, ok = rl.Query(point.Logging, "def")
	assert.False(t, ok)
	assert.Equal(t, "", p)

	name, ok := ScriptName(rl, point.Logging, point.NewPoint("abc", point.NewKVs(map[string]interface{}{"message@json": "a"})), nil)
	assert.True(t, ok)
	assert.Equal(t, "a1.p", name)

	name, ok = ScriptName(rl, point.Logging, point.NewPoint("abcd", point.NewKVs(map[string]interface{}{"message@json": "a"})), map[string]string{"abcd": "a2.p"})
	assert.True(t, ok)
	assert.Equal(t, "a2.p", name)
}

func TestRelationWildcard(t *testing.T) {
	rl := NewPipelineRelation()

	rl.UpdateRelation(0, map[point.Category]map[string]string{
		point.Logging: {
			"abc":  "exact.p",
			"ab*":  "star.p",
			"ab??": "two-char.p",
			"foo?": "single-char.p",
		},
	})

	p, ok := rl.Query(point.Logging, "abc")
	assert.True(t, ok)
	assert.Equal(t, "exact.p", p)

	p, ok = rl.Query(point.Logging, "ab12")
	assert.True(t, ok)
	assert.Equal(t, "two-char.p", p)

	p, ok = rl.Query(point.Logging, "ab123")
	assert.True(t, ok)
	assert.Equal(t, "star.p", p)

	p, ok = rl.Query(point.Logging, "foo1")
	assert.True(t, ok)
	assert.Equal(t, "single-char.p", p)

	p, ok = rl.Query(point.Logging, "foo12")
	assert.False(t, ok)
	assert.Equal(t, "", p)

	name, ok := ScriptName(rl, point.Logging, point.NewPoint("ab12", point.NewKVs(nil)), nil)
	assert.True(t, ok)
	assert.Equal(t, "two-char.p", name)
}

func TestRelationWildcardGenericAnchor(t *testing.T) {
	rl := NewPipelineRelation()

	rl.UpdateRelation(0, map[point.Category]map[string]string{
		point.Logging: {
			"*prod*api*": "prod-api.p",
			"*error":     "tail.p",
		},
	})

	p, ok := rl.Query(point.Logging, "aa-prod-zz-api-bb")
	assert.True(t, ok)
	assert.Equal(t, "prod-api.p", p)

	p, ok = rl.Query(point.Logging, "fatalerror")
	assert.True(t, ok)
	assert.Equal(t, "tail.p", p)
}

func TestRelationWildcardDirectAssignment(t *testing.T) {
	rl := NewPipelineRelation()
	rl.relation = map[point.Category]map[string]string{
		point.Logging: {
			"svc-*": "svc-star.p",
			"svc-?": "svc-one.p",
		},
	}

	p, ok := rl.Query(point.Logging, "svc-a")
	assert.True(t, ok)
	assert.Equal(t, "svc-one.p", p)

	p, ok = rl.Query(point.Logging, "svc-abc")
	assert.True(t, ok)
	assert.Equal(t, "svc-star.p", p)
}

func TestRelationQueryConcurrent(t *testing.T) {
	rl := NewPipelineRelation()
	rl.UpdateRelation(0, map[point.Category]map[string]string{
		point.Logging: {
			"svc-prod-*":  "prod-star.p",
			"*prod*api*":  "prod-api.p",
			"svc-prod-??": "prod-two.p",
			"*error":      "tail.p",
		},
	})

	cases := map[string]string{
		"svc-prod-ab":        "prod-two.p",
		"svc-prod-abc":       "prod-star.p",
		"aa-prod-zz-api":     "prod-api.p",
		"service-fatalerror": "tail.p",
	}

	var wg sync.WaitGroup
	for source, expected := range cases {
		source := source
		expected := expected

		for i := 0; i < 32; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 1000; j++ {
					name, ok := rl.Query(point.Logging, source)
					assert.True(t, ok)
					assert.Equal(t, expected, name)
				}
			}()
		}
	}

	wg.Wait()
}

func TestRelationUpdateAndQueryConcurrent(t *testing.T) {
	rl := NewPipelineRelation()

	stateA := map[point.Category]map[string]string{
		point.Logging: {
			"svc-prod-*":  "prod-star-a.p",
			"*prod*api*":  "prod-api-a.p",
			"svc-prod-??": "prod-two-a.p",
			"*error":      "tail-a.p",
		},
	}

	stateB := map[point.Category]map[string]string{
		point.Logging: {
			"svc-prod-*":  "prod-star-b.p",
			"*prod*api*":  "prod-api-b.p",
			"svc-prod-??": "prod-two-b.p",
			"*error":      "tail-b.p",
		},
	}

	allowed := map[string]map[string]struct{}{
		"svc-prod-ab": {
			"prod-two-a.p": {},
			"prod-two-b.p": {},
		},
		"svc-prod-abc": {
			"prod-star-a.p": {},
			"prod-star-b.p": {},
		},
		"aa-prod-zz-api": {
			"prod-api-a.p": {},
			"prod-api-b.p": {},
		},
		"service-fatalerror": {
			"tail-a.p": {},
			"tail-b.p": {},
		},
	}

	rl.UpdateRelation(1, stateA)

	var (
		stop atomic.Bool
		wg   sync.WaitGroup
		errC = make(chan error, 1)
	)

	reportErr := func(err error) {
		select {
		case errC <- err:
		default:
		}
		stop.Store(true)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 2000 && !stop.Load(); i++ {
			if i%2 == 0 {
				rl.UpdateRelation(int64(i+2), stateB)
			} else {
				rl.UpdateRelation(int64(i+2), stateA)
			}
		}
		stop.Store(true)
	}()

	for source, expected := range allowed {
		source := source
		expected := expected

		for i := 0; i < 16; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for !stop.Load() {
					name, ok := rl.Query(point.Logging, source)
					if !ok {
						reportErr(fmt.Errorf("query miss for source %q", source))
						return
					}
					if _, exists := expected[name]; !exists {
						reportErr(fmt.Errorf("unexpected result for source %q: %q", source, name))
						return
					}
				}
			}()
		}
	}

	wg.Wait()

	select {
	case err := <-errC:
		t.Fatal(err)
	default:
	}
}
