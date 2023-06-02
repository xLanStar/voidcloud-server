package tests

import "testing"

//go:noinline
func passArray(arr []int) int {
	return len(arr)
}

//go:noinline
func passDots(arr ...int) int {
	return len(arr)
}

func BenchmarkArray(b *testing.B) {
	arr := make([]int, 10000)
	ans := 0
	for i := 0; i < b.N; i++ {
		ans += passArray(arr)
	}
}

func BenchmarkDots(b *testing.B) {
	arr := make([]int, 10000)
	ans := 0
	for i := 0; i < b.N; i++ {
		ans += passDots(arr...)
	}
}
