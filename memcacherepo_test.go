package common_test

import (
	"context"
	"os"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
	common "github.com/papawattu/cleanlog-common"
)

type MockMemcacheClient struct {
	store map[string][]byte
}

func (m *MockMemcacheClient) Set(item *memcache.Item) error {
	m.store[item.Key] = item.Value
	return nil
}

func (m *MockMemcacheClient) Get(key string) (*memcache.Item, error) {
	if val, ok := m.store[key]; ok {
		return &memcache.Item{
			Key:   key,
			Value: val,
			Flags: 0,
		}, nil
	}
	return nil, nil
}

func (m *MockMemcacheClient) Delete(key string) error {
	delete(m.store, key)
	return nil
}
func TestMemcacheRepository(t *testing.T) {
	// Create a new MemcacheRepository

	mr := common.NewMemcacheRepository[*common.BaseEntity[string]]("localhost:11211", "test", &MockMemcacheClient{
		store: make(map[string][]byte),
	})

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
func TestMemcacheRepositoryWithRealMemcache(t *testing.T) {
	// Create a new MemcacheRepository

	if os.Getenv("MEMCACHE") == "" {
		t.Skip("MEMCACHE environment variable not set")
	}
	mr := common.NewMemcacheRepository[*common.BaseEntity[string]]("localhost:11211", "test", nil)

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
