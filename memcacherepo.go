package common

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"log"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheRepository[T any, S string] struct {
	client *memcache.Client
	host   string
}

func (mr *MemcacheRepository[T, S]) Create(ctx context.Context, entity T) error {

	id, err := mr.GetId(ctx, entity)
	if err != nil {
		return err
	}
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err = enc.Encode(entity)
	if err != nil {
		return err
	}

	h := sha256.New()

	h.Write(b.Bytes())

	err = mr.client.Set(&memcache.Item{
		Key:   string(id),
		Value: b.Bytes(),
		Flags: 0,
	})
	if err != nil {
		return err
	}
	log.Printf("Created entity with id: %v val: %v", id, entity)
	return nil
}
func (mr *MemcacheRepository[T, S]) Save(ctx context.Context, entity T) error {
	id, err := mr.GetId(ctx, entity)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err = enc.Encode(entity)
	if err != nil {
		return err
	}
	mr.client.Set(&memcache.Item{
		Key:   string(id),
		Value: b.Bytes(),
	})
	return nil
}

func (mr *MemcacheRepository[T, S]) Get(ctx context.Context, id S) (T, error) {

	var entity T

	item, err := mr.client.Get(string(id))
	if err != nil {
		var entity T
		return entity, err
	}

	dec := gob.NewDecoder(bytes.NewReader(item.Value))
	err = dec.Decode(&entity)
	if err != nil {
		return entity, err
	}

	return entity, nil
}
func (mr *MemcacheRepository[T, S]) GetAll(ctx context.Context) ([]T, error) {
	return nil, errors.New("Not implemented")
}
func (mr *MemcacheRepository[T, S]) Delete(ctx context.Context, e T) error {
	id, err := mr.GetId(ctx, e)
	if err != nil {
		return err
	}

	mr.client.Delete(string(id))

	return nil
}
func (mr *MemcacheRepository[T, S]) Exists(ctx context.Context, ID S) (bool, error) {
	_, err := mr.Get(ctx, ID)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (mr *MemcacheRepository[T, S]) GetId(ctx context.Context, entity T) (S, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(entity)
	if err != nil {
		return "", err
	}

	h := sha256.New()

	h.Write(b.Bytes())

	id := fmt.Sprintf("%x", h.Sum(nil))
	return S(id), nil
}

func (mr *MemcacheRepository[T, S]) SetHost(host string) error {
	mr.host = host
	return nil
}

func (mr *MemcacheRepository[T, S]) GetHost() string {
	return mr.host
}

func (mr *MemcacheRepository[T, S]) SetClient(client *memcache.Client) error {
	mr.client = client
	return nil
}

func (mr *MemcacheRepository[T, S]) GetClient() *memcache.Client {
	return mr.client
}

func NewMemcacheRepository[T any, S string](host string) Repository[T, S] {
	return &MemcacheRepository[T, S]{
		client: memcache.New(host),
	}
}
