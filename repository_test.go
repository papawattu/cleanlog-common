package common_test

import (
	"context"
	"testing"

	common "github.com/papawattu/cleanlog-common"
)

type TestRepo[T common.Entity[S], S comparable] struct {
	t *testing.T
}

func (t *TestRepo[T, S]) Create(ctx context.Context, e T) error {
	createCalled++
	return nil
}

func (t *TestRepo[T, S]) Save(ctx context.Context, e T) error {
	saveCalled++
	return nil
}

func (t *TestRepo[T, S]) Get(ctx context.Context, id S) (T, error) {
	getCalled++
	var zero T
	return zero, nil
}

func (t *TestRepo[T, S]) GetAll(ctx context.Context) ([]T, error) {
	getAllCalled++
	return nil, nil
}

func (t *TestRepo[T, S]) Delete(ctx context.Context, e T) error {
	deleteCalled++
	return nil
}

func (t *TestRepo[T, S]) Exists(ctx context.Context, id S) (bool, error) {
	existsCalled++
	return existsReturnValue, nil
}

func (t *TestRepo[T, S]) GetId(ctx context.Context, e T) (S, error) {
	getIdCalled++
	var zero S
	return zero, nil
}

func NewTestRepository[T common.Entity[S], S comparable]() common.Repository[T, S] {
	return &TestRepo[T, S]{}
}

func TestEventRepo(t *testing.T) {
	createCalled = 0
	saveCalled = 0
	getCalled = 0
	getAllCalled = 0
	deleteCalled = 0
	existsCalled = 0
	getIdCalled = 0

	existsReturnValue = false

	ctx := context.Background()
	// Create a new test repository
	repo := NewTestRepository[*common.BaseEntity[string]]()

	repo.Create(ctx, &common.BaseEntity[string]{})
	repo.Save(ctx, &common.BaseEntity[string]{})
	repo.Get(ctx, "1")
	repo.GetAll(ctx)
	repo.Delete(ctx, &common.BaseEntity[string]{})
	repo.Exists(ctx, "1")
	repo.GetId(ctx, &common.BaseEntity[string]{})

	if createCalled != 1 {
		t.Errorf("Create not called")
	}

	if saveCalled != 1 {
		t.Errorf("Save not called")
	}

	if getCalled != 1 {
		t.Errorf("Get not called")
	}

	if getAllCalled != 1 {
		t.Errorf("GetAll not called")
	}

	if deleteCalled != 1 {
		t.Errorf("Delete not called")
	}

	if existsCalled != 1 {
		t.Errorf("Exists not called")
	}

	if getIdCalled != 1 {
		t.Errorf("GetId not called")
	}

}
