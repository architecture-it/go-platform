package log

import (
	"testing"
)

func TestBenchmark(t *testing.T) {
	defer Benchmark(Info)
}
