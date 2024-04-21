package main

import (
	"testing"
)

func BenchmarkHandle(b *testing.B) {
	filename := "../../../../measurements1000000.txt"

	for i := 0; i < b.N; i++ {

		Handle(filename)
	}

}
