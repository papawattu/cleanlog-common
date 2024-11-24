package common

import (
	"context"
	"errors"
)

type InMemoryRepository[T Entity[S], S comparable] struct {
	entities map[S]*T
}

func (wri *InMemoryRepository[T, S]) Create(ctx context.Context, e T) error {

	id := e.GetID()

	if _, ok := wri.entities[id]; ok {
		return errors.New("entity already exists")
	}

	wri.entities[id] = &e

	return nil
}
func (wri *InMemoryRepository[T, S]) Save(ctx context.Context, e T) error {

	id := e.GetID()

	if _, ok := wri.entities[id]; !ok {
		return errors.New("entity not found")
	}

	wri.entities[id] = &e

	return nil
}

func (wri *InMemoryRepository[T, S]) Get(ctx context.Context, id S) (T, error) {

	var zero T
	wl, ok := wri.entities[id]
	if !ok {
		return zero, nil
	}

	return *wl, nil
}

func (wri *InMemoryRepository[T, S]) GetAll(ctx context.Context) ([]T, error) {

	es := []T{}
	for _, e := range wri.entities {
		es = append(es, *e)
	}

	return es, nil
}

func (wri *InMemoryRepository[T, S]) Delete(ctx context.Context, e T) error {

	id, err := wri.GetId(ctx, e)
	if err != nil {
		return err
	}
	if _, ok := wri.entities[id]; !ok {
		return errors.New("entity not found")
	}
	delete(wri.entities, id)
	return nil
}

func (wri *InMemoryRepository[T, S]) GetId(ctx context.Context, e T) (S, error) {

	id := e.GetID()

	return id, nil
}

func (wri *InMemoryRepository[T, S]) Exists(ctx context.Context, id S) (bool, error) {

	_, ok := wri.entities[id]
	return ok, nil
}
func NewInMemoryRepository[T Entity[S], S comparable]() Repository[T, S] {
	return &InMemoryRepository[T, S]{
		entities: make(map[S]*T),
	}
}
