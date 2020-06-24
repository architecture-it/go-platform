package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBenchmark(t *testing.T) {
	defer Benchmark("Info")
}

func TestPipeline(t *testing.T) {
	res := Trace.Pipeline("PRUEBA")
	assert.Equal(t, res, "2020-06-24 09:27:20.734 | -1 | /tmp/go-build402542673/b001/log.test | TRACE | log.go:18 | PRUEBA")
}

func TestJSON(t *testing.T) {
	res := Trace.JSON("TEST")
	assert.Equal(t, res, "{\"Date\":\"2020-06-24 09:28:02.791\",\"Level\":\"-1\",\"Local\":\"/tmp/go-build488598112/b001/log.test\",\"Name\":\"TRACE\",\"Path\":\"log.go:18\",\"Message\":\"TEST\"}")
}
