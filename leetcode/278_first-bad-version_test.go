package leetcode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Benchmark_firstBadVersion(b *testing.B) {
	benchmarks := []struct {
		name    string
		version int
		fn      func(n int) int
	}{
		{
			name:    "near directly",
			version: 6,
			fn:      directlyFindVersion,
		},
		{
			name:    "far directly",
			version: 6000,
			fn:      directlyFindVersion,
		},
		{
			name:    "near incr",
			version: 6,
			fn:      incrFindVersion,
		},
		{
			name:    "far incr",
			version: 6000,
			fn:      incrFindVersion,
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bm.fn(bm.version)
			}
		})
	}
}

func Test_firstBadVersion(t *testing.T) {
	tests := []struct {
		name        string
		version     int
		wantVersion int
	}{
		{
			name:        ">first bad version",
			version:     5,
			wantVersion: 4,
		},
		{
			name:        "=first bad version",
			version:     4,
			wantVersion: 4,
		},
		{
			name:        "<first bad version",
			version:     3,
			wantVersion: 4,
		},
		{
			name:        ">first far bad version",
			version:     600,
			wantVersion: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantVersion, directlyFindVersion(tt.version), tt.name)
			assert.Equalf(t, tt.wantVersion, incrFindVersion(tt.version), tt.name)
		})
	}
}
