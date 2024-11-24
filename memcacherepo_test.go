package common_test

import (
	"context"
	"testing"

	common "github.com/papawattu/cleanlog-common"
)

func TestMemcacheRepository(t *testing.T) {
	// Create a new MemcacheRepository
	mr := common.NewMemcacheRepository[*common.BaseEntity[string], string]("localhost:11211")

	// Create a new context
	ctx := context.Background()

	// Create a new entity

	err := mr.Create(ctx, &common.BaseEntity[string]{
		ID: "1",
	})

	if err != nil {
		t.Errorf("Error creating entity: %v", err)
	}

	// Test Get
	id, err := mr.GetId(context.Background(), &common.BaseEntity[string]{
		ID: "1",
	})
	if err != nil {
		t.Errorf("Error getting entity id: %v", err)
	}
	_, err = mr.Get(ctx, id)
	if err != nil {
		t.Errorf("Error getting entity: %v", err)
	}

	// Test Exists
	e, err := mr.Exists(ctx, id)
	if err != nil {
		t.Errorf("Error checking if entity exists: %v", err)
	}

	if !e {
		t.Errorf("Entity should exist")
	}

	// Test Delete
	err = mr.Delete(ctx, &common.BaseEntity[string]{
		ID: "1",
	})
	if err != nil {
		t.Errorf("Error deleting entity: %v", err)
	}
	e, err = mr.Exists(ctx, id)
	if err != nil {
		t.Errorf("Error checking if entity exists: %v", err)
	}

	if e {
		t.Errorf("Entity should not exist")
	}
	// Test GetId

	_, err = mr.GetId(ctx, &common.BaseEntity[string]{
		ID: "1",
	})
	if err != nil {
		t.Errorf("Error getting entity id: %v", err)
	}
}
