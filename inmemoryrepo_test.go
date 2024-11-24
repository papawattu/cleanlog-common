package common_test

import (
	"context"
	"testing"

	common "github.com/papawattu/cleanlog-common"
)

func TestInMemoryRepo(t *testing.T) {
	// Create a new InMemoryRepository
	repo := common.NewInMemoryRepository[*common.BaseEntity[int]]()

	// Create a new context
	ctx := context.Background()

	// Create a new entity

	// Test Create

	te := common.NewBaseEntity(1)
	err := repo.Create(ctx, te.(*common.BaseEntity[int]))

	if err != nil {
		t.Errorf("Error creating entity: %v", err)
	}

	// Test Get
	id, err := repo.GetId(context.Background(), te.(*common.BaseEntity[int]))
	if err != nil {
		t.Errorf("Error getting entity id: %v", err)
	}
	_, err = repo.Get(ctx, id)
	if err != nil {
		t.Errorf("Error getting entity: %v", err)
	}

	// Test Exists
	e, err := repo.Exists(ctx, id)
	if err != nil {
		t.Errorf("Error checking if entity exists: %v", err)
	}

	if !e {
		t.Errorf("Entity should exist")
	}

	// Test Delete
	err = repo.Delete(ctx, te.(*common.BaseEntity[int]))
	if err != nil {
		t.Errorf("Error deleting entity: %v", err)
	}
	e, err = repo.Exists(ctx, id)
	if err != nil {
		t.Errorf("Error checking if entity exists: %v", err)
	}

	if e {
		t.Errorf("Entity should not exist")
	}
}
