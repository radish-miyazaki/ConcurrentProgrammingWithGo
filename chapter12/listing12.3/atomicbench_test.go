package main

import (
	"sync/atomic"
	"testing"
)

var total = int64(0)

func BenchmarkNormal(bench *testing.B) {
	for range bench.N {
		total += 1
	}
}

func BenchmarkAtomic(bench *testing.B) {
	for range bench.N {
		atomic.AddInt64(&total, 1)
	}
}
