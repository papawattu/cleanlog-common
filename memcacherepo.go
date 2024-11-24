package common

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"log"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheClient interface {
	Set(item *memcache.Item) error
	Get(key string) (*memcache.Item, error)
	Delete(key string) error
}

type MemcacheRepository[T Entity[S], S string] struct {
	client MemcacheClient
	host   string
	prefix string
}

func (mr MemcacheRepository[T, S]) Create(ctx context.Context, e T) error {

	id := e.GetID()

	if _, err := mr.Get(ctx, id); err == nil {
		return errors.New("entity already exists")
	}

	var b bytes.Buffer
	enc := gob.NewEncoder(&b)

	err := enc.Encode(e)
	if err != nil {
		return err
	}

	h := sha256.New()

	h.Write(b.Bytes())
	err = mr.client.Set(&memcache.Item{
		Key:   mr.prefix + string(id),
		Value: b.Bytes(),
		Flags: 0,
	})
	if err != nil {
		return err
	}

	log.Printf("Created entity with id: %s val: %+v", id, e)
	return nil
}
func (mr *MemcacheRepository[T, S]) Save(ctx context.Context, e T) error {
	id := e.GetID()

	if _, err := mr.Get(ctx, id); err != nil {
		return errors.New("entity not found")
	}

	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(e)
	if err != nil {
		return err
	}
	mr.client.Set(&memcache.Item{
		Key:   mr.prefix + string(id),
		Value: b.Bytes(),
	})
	return nil
}

func (mr *MemcacheRepository[T, S]) Get(ctx context.Context, id S) (T, error) {

	var entity T

	item, err := mr.client.Get(mr.prefix + string(id))
	if err != nil {

		return entity, err
	}

	if item == nil {
		return entity, errors.New("entity not found")
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

	mr.client.Delete(mr.prefix + string(id))

	return nil
}
func (mr *MemcacheRepository[T, S]) Exists(ctx context.Context, ID S) (bool, error) {
	_, err := mr.Get(ctx, ID)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (mr *MemcacheRepository[T, S]) GetId(ctx context.Context, e T) (S, error) {
	return e.GetID(), nil
}

func (mr *MemcacheRepository[T, S]) SetHost(host string) error {
	mr.host = host
	return nil
}

func (mr *MemcacheRepository[T, S]) GetHost() string {
	return mr.host
}

func (mr *MemcacheRepository[T, S]) SetClient(client MemcacheClient) error {
	mr.client = client
	return nil
}

func (mr *MemcacheRepository[T, S]) GetClient() MemcacheClient {
	return mr.client
}

func NewMemcacheRepository[T Entity[S], S string](host string, prefix string, mc MemcacheClient) Repository[T, S] {
	if mc == nil {
		mc = memcache.New(host)
	}
	return &MemcacheRepository[T, S]{
		client: mc,
		host:   host,
		prefix: prefix,
	}
}
