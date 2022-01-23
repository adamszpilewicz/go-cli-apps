package main

import (
	"io"
	"log"
	"path/filepath"
	"testing"
)

func BenchmarkRun(b *testing.B) {
	filenames, err := filepath.Glob("../testFiles/benchmark/*.csv")
	if err != nil {
		log.Fatal(err)
	}

	// to ignore any time for the benchmark execution
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := run(filenames, "sum", 2, io.Discard); err != nil {
			b.Error(err)
		}
	}
}
