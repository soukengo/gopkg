package env

import (
	"fmt"
	"testing"
)
func TestRootDir(t *testing.T) {
	fmt.Println(RootDir())
}

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Println(RootDir())
	}
}