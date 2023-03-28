package idgen

import (
	"fmt"
	"testing"
)

func Benchmark_sonyflakeGenerator_GenID(b *testing.B) {
	var generator = NewSonyflakeGenerator()
	for i := 0; i < b.N; i++ {
		_, err := generator.GenID()
		if err != nil {
			b.Fatal(err)
		}
	}
}
func Benchmark_mongoGenerator_GenID(b *testing.B) {
	var generator = NewMongoGenerator()
	for i := 0; i < b.N; i++ {
		_, err := generator.GenID()
		if err != nil {
			b.Fatal(err)
		}
	}
}
func Benchmark_xidGenerator_GenID(b *testing.B) {
	var generator = NewXidGenerator()
	for i := 0; i < b.N; i++ {
		_, err := generator.GenID()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Test_xidGenerator_GenID(t *testing.T) {
	var generator = NewXidGenerator()
	id, err := generator.GenID()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(id)
}
