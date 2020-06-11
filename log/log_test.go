package main

import (
	"testing"
)

func TestBenchmark(t *testing.T) {
	defer Benchmark("Info")
}
