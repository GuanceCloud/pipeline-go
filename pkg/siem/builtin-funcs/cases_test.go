package funcs

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/GuanceCloud/pipeline-go/pkg/siem/trigger"
	"github.com/GuanceCloud/platypus/pkg/engine"
	"github.com/GuanceCloud/platypus/pkg/engine/runtimev2"
	"github.com/stretchr/testify/assert"
)

func TestDocs(t *testing.T) {
	d, err := GenerateDocs()
	assert.NoError(t, err)
	_ = os.WriteFile("docs/docs.md", []byte(d), 0644)
}

func runCase(t *testing.T, c ProgCase, private ...map[runtimev2.TaskP]any) {
	s, err := engine.ParseV2(c.Name, c.Script, Funcs)
	if err != nil {
		t.Error(err)
		return
	}

	var privateMap map[runtimev2.TaskP]any
	if len(private) > 0 && private[0] != nil {
		privateMap = private[0]
	} else {
		privateMap = map[runtimev2.TaskP]any{}
	}

	stdout := bytes.NewBuffer([]byte{})
	privateMap[PStdout] = stdout
	tr := trigger.NewTr()
	privateMap[PTrigger] = tr
	if err := s.Run(nil, runtimev2.WithPrivate(privateMap)); err != nil {
		t.Error(err.Error())
	}
	o := stdout.String()
	t.Log(o)

	trBuf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(trBuf)
	enc.SetIndent("", "  ")
	_ = enc.Encode(tr.Result())
	t.Log(trBuf.String())
	if c.jsonout {
		var v1 any
		var v2 any
		err = json.Unmarshal([]byte(o), &v1)
		assert.NoError(t, err)
		err = json.Unmarshal([]byte(c.Stdout), &v2)
		assert.NoError(t, err)
		assert.Equal(t, v1, v2)
	} else {
		assert.Equal(t, c.Stdout, o)
	}
	assert.Equal(t, tr.Result(), c.TriggerResult)
}
