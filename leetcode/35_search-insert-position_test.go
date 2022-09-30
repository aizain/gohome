package leetcode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_searchInsert(t *testing.T) {
	type args struct {
		nums   []int
		target int
	}
	tests := []struct {
		name      string
		nums      []int
		target    int
		wantIndex int
	}{
		{
			name:      "empty nums",
			nums:      []int{},
			target:    1,
			wantIndex: 0,
		},
		{
			name:      "mid insert",
			nums:      []int{1, 2, 4, 5},
			target:    3,
			wantIndex: 2,
		},
		{
			name:      "dup insert",
			nums:      []int{1, 2, 4, 5},
			target:    4,
			wantIndex: 2,
		},
		{
			name:      "dup 2 insert",
			nums:      []int{1, 2, 4, 4, 5},
			target:    4,
			wantIndex: 2,
		},
		{
			name:      "head insert",
			nums:      []int{1, 2, 4, 4, 5},
			target:    0,
			wantIndex: 0,
		},
		{
			name:      "tail insert",
			nums:      []int{1, 2, 4, 4, 5},
			target:    8,
			wantIndex: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantIndex, searchInsert(tt.nums, tt.target), tt.name)
		})
	}
}
