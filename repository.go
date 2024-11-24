package common

import (
	"context"
	"time"
)

type Entity[S comparable] interface {
	GetID() S
	GetLastUpdateDate() time.Time
	GetCreationDate() time.Time
	GetVersion() int
	SetVersion(v int)
	SetLastUpdateDate(t time.Time)
	SetCreationDate(t time.Time)
}

type Repository[T Entity[S], S comparable] interface {
	Create(ctx context.Context, entity T) error
	Save(ctx context.Context, entity T) error
	Get(ctx context.Context, ID S) (T, error)
	GetAll(ctx context.Context) ([]T, error)
	Delete(ctx context.Context, e T) error
	Exists(ctx context.Context, ID S) (bool, error)
	GetId(ctx context.Context, e T) (S, error)
}

type BaseEntity[S comparable] struct {
	ID             S         `json:"id"`
	LastUpdateDate time.Time `json:"lastUpdateDate"`
	CreationDate   time.Time `json:"creationDate"`
	Version        int       `json:"version"`
}

func (b *BaseEntity[S]) GetID() S {
	return b.ID
}

func (b *BaseEntity[S]) GetLastUpdateDate() time.Time {
	return b.LastUpdateDate
}

func (b *BaseEntity[S]) GetCreationDate() time.Time {
	return b.CreationDate
}

func (b *BaseEntity[S]) GetVersion() int {
	return b.Version
}

func (b *BaseEntity[S]) SetVersion(v int) {
	b.Version = v
}

func (b *BaseEntity[S]) SetLastUpdateDate(t time.Time) {
	b.LastUpdateDate = t
}

func (b *BaseEntity[S]) SetCreationDate(t time.Time) {
	b.CreationDate = t
}

func NewBaseEntity[S comparable](id S) Entity[S] {
	return &BaseEntity[S]{ID: id, LastUpdateDate: time.Now(), CreationDate: time.Now(), Version: 1}
}
