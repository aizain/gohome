package main

import (
	"errors"
)

type List[T any] interface {
	Add(e ...T) error
	Set(i int, e T) (T, error)
	Get(i int) (T, error)
	Del(i int) (T, error)
	Size() int
}

func NewArrayList[T any](items []T) *ArrayList[T] {
	return &ArrayList[T]{
		items: items,
	}
}

type ArrayList[T any] struct {
	items []T
}

func (a *ArrayList[T]) Add(e ...T) error {
	a.items = append(a.items, e...)
	return nil
}

func (a *ArrayList[T]) Set(i int, e T) (T, error) {
	if err := a.checkBound(i); err != nil {
		var empty T
		return empty, err
	}
	old := a.items[i]
	a.items[i] = e
	return old, nil
}

func (a *ArrayList[T]) Get(i int) (T, error) {
	if err := a.checkBound(i); err != nil {
		var empty T
		return empty, err
	}
	return a.items[i], nil
}

func (a *ArrayList[T]) Del(i int) (T, error) {
	if err := a.checkBound(i); err != nil {
		var empty T
		return empty, err
	}
	deleted := a.items[i]
	if i == a.Size()-1 {
		a.items = a.items[:i]
	} else {
		a.items = append(a.items[:i], a.items[i+1:]...)
	}
	return deleted, nil
}

func (a *ArrayList[T]) Size() int {
	if a == nil {
		return 0
	}
	return len(a.items)
}

func (a *ArrayList[T]) checkBound(i int) error {
	if i >= a.Size() || i < 0 {
		return errors.New("数组越界")
	}
	return nil
}
