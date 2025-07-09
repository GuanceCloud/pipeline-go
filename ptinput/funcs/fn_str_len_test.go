package funcs

import (
	"testing"
	"time"

	"github.com/GuanceCloud/cliutils/point"
	tu "github.com/GuanceCloud/cliutils/testutil"
	"github.com/GuanceCloud/pipeline-go/ptinput"
)

func TestStrlen(t *testing.T) {
	cases := []struct {
		name, pl, in string
		expected     any
		fail         bool
		key          string
	}{
		{
			name: "t1",
			pl: `
			add_key("k1", strlen("你好"))
			`,
			expected: int64(2),
			key:      "k1",
		},
		{
			name: "t2",
			pl: `
			add_key("k1", strlen("hello"))
			`,
			expected: int64(5),
			key:      "k1",
		},
		{
			name: "t3",
			pl: `
			add_key("k1", strlen("你好hello"))
			`,
			expected: int64(7),
			key:      "k1",
		},
		{
			name: "t4",
			pl: `
			v = []
			v = append(v, strlen("hello你好"))
			v = append(v, len("hello你好"))
			add_key("v", v)
			`,
			expected: "[7,11]",
			key:      "v",
		},
	}

	for idx, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runner, err := NewTestingRunner(tc.pl)
			if err != nil {
				if tc.fail {
					t.Logf("[%d]expect error: %s", idx, err)
				} else {
					t.Errorf("[%d] failed: %s", idx, err)
				}
				return
			}
			pt := ptinput.NewPlPt(
				point.Logging, "test", nil, map[string]any{"message": tc.in}, time.Now())
			errR := runScript(runner, pt)

			if errR != nil {
				t.Fatal(errR.Error())
			}

			v, _, _ := pt.Get(tc.key)
			// tu.Equals(t, nil, err)
			tu.Equals(t, tc.expected, v)
			t.Logf("[%d] PASS", idx)
		})
	}
}
