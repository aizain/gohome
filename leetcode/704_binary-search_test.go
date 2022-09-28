package leetcode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var nums []int

func init() {
	nums = make([]int, 0, 100000)
	for i := 0; i < 100000; i++ {
		nums = append(nums, i)
	}
}

func Benchmark_directlySearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		directlySearch(nums, 876291)
	}
}

func Benchmark_binarySearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		binarySearch(nums, 876291)
	}
}

func Test_search(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   int
	}{
		{
			name:   "exist",
			nums:   []int{-1, 0, 3, 5, 9, 12},
			target: 9,
			want:   4,
		},
		{
			name:   "not exist",
			nums:   []int{-1, 0, 3, 5, 9, 12},
			target: 2,
			want:   -1,
		},
		{
			name:   "nil",
			nums:   nil,
			target: 2,
			want:   -1,
		},
		{
			name:   "-n",
			nums:   []int{-1, 0, 5},
			target: -1,
			want:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := directlySearch(tt.nums, tt.target); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
		t.Run(tt.name, func(t *testing.T) {
			if got := binarySearch(tt.nums, tt.target); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
