package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayList_Add(t *testing.T) {
	tests := []struct {
		name      string
		list      *ArrayList[any]
		newValues []any
		wantErr   error
		wantItems []any
	}{
		{
			name:      "add one element",
			list:      NewArrayList([]any{1, 2, 3}),
			newValues: []any{1},
			wantErr:   nil,
			wantItems: []any{1, 2, 3, 1},
		},
		{
			name:      "add more element",
			list:      NewArrayList([]any{1, 2, 3}),
			newValues: []any{1, 2, 3},
			wantErr:   nil,
			wantItems: []any{1, 2, 3, 1, 2, 3},
		},
		{
			name:      "add different element",
			list:      NewArrayList([]any{1}),
			newValues: []any{[]any{1}, 2, struct{ name string }{name: "name"}},
			wantErr:   nil,
			wantItems: []any{1, []any{1}, 2, struct{ name string }{name: "name"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.list.Add(tt.newValues...)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantItems, tt.list.items)
		})
	}
}

func TestArrayList_Del(t *testing.T) {
	tests := []struct {
		name      string
		list      *ArrayList[any]
		index     int
		wantValue any
		wantItems []any
		wantErr   error
	}{
		{
			name:      "del on element",
			list:      NewArrayList([]any{1, 2}),
			index:     1,
			wantValue: 2,
			wantItems: []any{1},
			wantErr:   nil,
		},
		{
			name:      "del out of bound",
			list:      NewArrayList([]any{1, 2}),
			index:     2,
			wantValue: nil,
			wantItems: []any{1, 2},
			wantErr:   errors.New("数组越界"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.list.Del(tt.index)
			assert.Equal(t, err, tt.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, got, tt.wantValue)
			assert.Equal(t, tt.list.items, tt.wantItems)
		})
	}
}

func TestArrayList_Get(t *testing.T) {
	tests := []struct {
		name      string
		list      *ArrayList[any]
		index     int
		wantValue any
		wantErr   error
	}{
		{
			name:      "get on element",
			list:      NewArrayList([]any{1, 2}),
			index:     0,
			wantValue: 1,
			wantErr:   nil,
		},
		{
			name:      "get out of bound",
			list:      NewArrayList([]any{1, 2}),
			index:     4,
			wantValue: nil,
			wantErr:   errors.New("数组越界"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.list.Get(tt.index)
			assert.Equal(t, err, tt.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, got, tt.wantValue)
		})
	}
}

func TestArrayList_Set(t *testing.T) {
	tests := []struct {
		name      string
		list      *ArrayList[any]
		index     int
		value     any
		wantValue any
		wantItems []any
		wantErr   error
	}{
		{
			name:      "set on element",
			list:      NewArrayList([]any{1, 2}),
			index:     1,
			value:     3,
			wantValue: 2,
			wantItems: []any{1, 3},
		},
		{
			name:      "set different element",
			list:      NewArrayList([]any{1, 2}),
			index:     1,
			value:     []any{4},
			wantValue: 2,
			wantItems: []any{1, []any{4}},
		},
		{
			name:    "set out of bound",
			list:    NewArrayList([]any{1, 2}),
			index:   4,
			wantErr: errors.New("数组越界"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.list.Set(tt.index, tt.value)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantValue, got)
			assert.Equal(t, tt.wantItems, tt.list.items)
		})
	}
}

func TestArrayList_Size(t *testing.T) {
	tests := []struct {
		name      string
		list      *ArrayList[any]
		wantValue int
	}{
		{
			name:      "has size",
			list:      NewArrayList([]any{1, 2}),
			wantValue: 2,
		},
		{
			name:      "empty size",
			list:      NewArrayList([]any{}),
			wantValue: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.list.Size()
			assert.Equal(t, tt.wantValue, got)
		})
	}
}
