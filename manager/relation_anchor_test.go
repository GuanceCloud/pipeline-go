package manager

import (
	"math/rand"
	"testing"

	"github.com/GuanceCloud/cliutils/point"
	"github.com/stretchr/testify/assert"
)

func TestLiteralAnchorIndexQuery(t *testing.T) {
	rel := map[string]string{
		"svc-prod-??": "prod-two.p",
		"svc-prod-*":  "prod-star.p",
		"*prod*api*":  "prod-api.p",
		"*error":      "tail.p",
		"plain":       "plain.p",
	}

	idx := buildLiteralAnchorIndex(rel)

	cases := []struct {
		source string
		name   string
		ok     bool
	}{
		{source: "plain", name: "plain.p", ok: true},
		{source: "svc-prod-ab", name: "prod-two.p", ok: true},
		{source: "svc-prod-abc", name: "prod-star.p", ok: true},
		{source: "aa-prod-zz-api-bb", name: "prod-api.p", ok: true},
		{source: "fatalerror", name: "tail.p", ok: true},
		{source: "nomatch", ok: false},
	}

	for _, tc := range cases {
		name, ok := idx.query(tc.source)
		assert.Equal(t, tc.ok, ok, tc.source)
		assert.Equal(t, tc.name, name, tc.source)
	}
}

func TestLiteralAnchorIndexEquivalence(t *testing.T) {
	rng := rand.New(rand.NewSource(2))
	alphabet := []rune("abcdxyz012")

	for round := 0; round < 40; round++ {
		rel := make(map[string]string)
		wildcards := 0

		for len(rel) < 64 {
			pattern := randomPattern(rng, alphabet)
			if _, exists := rel[pattern]; exists {
				continue
			}

			if hasWildcard(pattern) {
				wildcards++
			}
			rel[pattern] = pattern + ".p"
		}

		if wildcards == 0 {
			rel["*"] = "wild-any.p"
		}

		rl := NewPipelineRelation()
		rl.UpdateRelation(0, map[point.Category]map[string]string{
			point.Logging: rel,
		})
		idx := buildLiteralAnchorIndex(rel)

		for i := 0; i < 200; i++ {
			source := randomSource(rng, alphabet)

			expectedName, expectedOK := rl.Query(point.Logging, source)
			actualName, actualOK := idx.query(source)

			assert.Equal(t, expectedOK, actualOK, "round=%d source=%q", round, source)
			assert.Equal(t, expectedName, actualName, "round=%d source=%q", round, source)
		}
	}
}

func randomPattern(rng *rand.Rand, alphabet []rune) string {
	n := 1 + rng.Intn(8)
	buf := make([]rune, 0, n)
	hasWildcardRune := false

	for i := 0; i < n; i++ {
		switch rng.Intn(8) {
		case 0:
			buf = append(buf, '*')
			hasWildcardRune = true
		case 1:
			buf = append(buf, '?')
			hasWildcardRune = true
		default:
			buf = append(buf, alphabet[rng.Intn(len(alphabet))])
		}
	}

	if !hasWildcardRune && rng.Intn(2) == 0 {
		pos := rng.Intn(len(buf))
		if rng.Intn(2) == 0 {
			buf[pos] = '*'
		} else {
			buf[pos] = '?'
		}
	}

	return string(buf)
}

func randomSource(rng *rand.Rand, alphabet []rune) string {
	n := rng.Intn(10)
	buf := make([]rune, n)
	for i := range buf {
		buf[i] = alphabet[rng.Intn(len(alphabet))]
	}
	return string(buf)
}
