package lang

import (
	"testing"
	"time"

	"github.com/GuanceCloud/cliutils/point"
	"github.com/GuanceCloud/pipeline-go/constants"
	"github.com/GuanceCloud/pipeline-go/ptinput"
	"github.com/stretchr/testify/assert"
)

type pt4t struct {
	Fields map[string]interface{}
	Tags   map[string]string
	Drop   bool
}

func BenchmarkGTags(b *testing.B) {
	b.Run("test-list", func(b *testing.B) {
		s := [][2]string{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}, {"6"}}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for idx := range s {
				_ = s[idx][0]
				_ = s[idx][1]
			}
		}
	})

	b.Run("test-map", func(b *testing.B) {
		s := map[string]string{
			"1": "",
			"2": "",
			"3": "",
			"4": "",
			"5": "",
			"6": "",
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for k, v := range s {
				_ = k
				_ = v
			}
		}
	})
}

func TestStatus(t *testing.T) {
	for k, v := range statusMap {
		outp := &pt4t{
			Fields: map[string]interface{}{
				constants.FieldStatus: k,
			},
		}
		pt := ptinput.NewPlPt(point.Logging, "", nil, outp.Fields, time.Now())
		ProcLoggingStatus(pt, false, nil)
		assert.Equal(t, v, pt.Fields()[constants.FieldStatus])
	}

	{
		outp := &pt4t{
			Fields: map[string]interface{}{
				constants.FieldStatus:  "x",
				constants.FieldMessage: "1234567891011",
			},
		}
		pt := ptinput.NewPlPt(point.Logging, "", nil, outp.Fields, time.Now())
		ProcLoggingStatus(pt, false, nil)
		assert.Equal(t, "x", pt.Fields()[constants.FieldStatus])
		assert.Equal(t, "1234567891011", pt.Fields()[constants.FieldMessage])
	}

	{
		outp := &pt4t{
			Fields: map[string]interface{}{
				constants.FieldStatus:  "x",
				constants.FieldMessage: "1234567891011",
			},
			Tags: map[string]string{
				"xxxqqqddd": "1234567891011",
			},
		}
		pt := ptinput.NewPlPt(point.Logging, "", outp.Tags, outp.Fields, time.Now())
		ProcLoggingStatus(pt, false, nil)
		assert.Equal(t, map[string]interface{}{
			constants.FieldStatus:  "x",
			constants.FieldMessage: "1234567891011",
		}, pt.Fields())
		assert.Equal(t, map[string]string{
			"xxxqqqddd": "1234567891011",
		}, pt.Tags())
	}

	{
		outp := &pt4t{
			Fields: map[string]interface{}{
				constants.FieldStatus:  "n",
				constants.FieldMessage: "1234567891011",
			},
			Tags: map[string]string{
				"xxxqqqddd": "1234567891011",
			},
		}
		pt := ptinput.NewPlPt(point.Logging, "", outp.Tags, outp.Fields, time.Now())
		ProcLoggingStatus(pt, false, nil)
		assert.Equal(t, map[string]interface{}{
			constants.FieldStatus:  "notice",
			constants.FieldMessage: "1234567891011",
		}, pt.Fields())
		assert.Equal(t, map[string]string{
			"xxxqqqddd": "1234567891011",
		}, pt.Tags())
	}
}

func TestGetSetStatus(t *testing.T) {
	out := &pt4t{
		Tags: map[string]string{
			"status": "n",
		},
		Fields: make(map[string]interface{}),
	}

	pt := ptinput.NewPlPt(point.Logging, "", out.Tags, out.Fields, time.Now())
	ProcLoggingStatus(pt, false, nil)
	assert.Equal(t, map[string]string{
		"status": "notice",
	}, pt.Tags())
	assert.Equal(t, make(map[string]interface{}), pt.Fields())

	out.Fields = map[string]interface{}{
		"status": "n",
	}
	out.Tags = make(map[string]string)
	pt = ptinput.NewPlPt(point.Logging, "", out.Tags, out.Fields, time.Now())

	ProcLoggingStatus(pt, false, nil)
	assert.Equal(t, map[string]interface{}{
		"status": "notice",
	}, pt.Fields())
	assert.Equal(t, make(map[string]string), pt.Tags())

	out.Tags = map[string]string{
		"status": "n",
	}

	pt = ptinput.NewPlPt(point.Logging, "", out.Tags, out.Fields, time.Now())
	ProcLoggingStatus(pt, false, nil)
	assert.Equal(t, map[string]string{
		"status": "notice",
	}, pt.Tags())

	pt = ptinput.NewPlPt(point.Logging, "", out.Tags, out.Fields, time.Now())
	ProcLoggingStatus(pt, false, []string{"notice"})
	assert.Equal(t, map[string]string{
		"status": "notice",
	}, pt.Tags())
}
