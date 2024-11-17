package common

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"time"
)

const (
	Created = "Created"
	Deleted = "Deleted"
	Updated = "Updated"
	Version = 1
)

type Event struct {
	EventId      string
	EventSHA     string
	EventType    string
	EventData    string
	EventVersion int
	EventTime    time.Time
}

type EventHandler func(event Event) error

type EventHandlers map[string]EventHandler

type Transport interface {
	PostEvent(event Event) error
	NextEvent() (*Event, error)
}
type EventService[T any, S comparable] interface {
	Repository[T, S]
	Transport
	SetPrefix(prefix string)
	SetHandlers(handlers EventHandlers)
	HandleEvent(event Event) error
}

type EventServiceImpl[T any, S comparable] struct {
	Repository[T, S]
	Transport
	Prefix   string
	Handlers EventHandlers
}

func (es *EventServiceImpl[T, S]) SetPrefix(prefix string) {
	es.Prefix = prefix
}

func (es *EventServiceImpl[T, S]) SetHandlers(handlers EventHandlers) {
	es.Handlers = handlers
}

func (es *EventServiceImpl[T, S]) Create(ctx context.Context, e T) error {
	slog.Info("EventService", "Create", e)
	ent, err := json.Marshal(e)

	if err != nil {
		return err
	}

	// Broadcast event
	event := Event{
		EventId:      "1",
		EventType:    es.Prefix + Created,
		EventTime:    time.Now(),
		EventVersion: Version,
		EventData:    string(ent),
	}

	slog.Info("EventBroadcaster", "Create", event.EventData)
	err = es.PostEvent(event)

	if err != nil {
		slog.Error("Error broadcasting event", "error", err)
		return err
	}

	slog.Info("EventBroadcaster", "Create", "Event published")

	return nil

}

func (es *EventServiceImpl[T, S]) Save(ctx context.Context, e T) error {
	slog.Info("EventService", "Save", e)

	ent, err := json.Marshal(e)

	if err != nil {
		slog.Error("Error marshalling event", "error", err)
	}

	// Broadcast event
	event := Event{
		EventId:      "1",
		EventType:    es.Prefix + Updated,
		EventTime:    time.Now(),
		EventVersion: Version,
		EventData:    string(ent),
	}

	slog.Info("EventBroadcaster", "Save", event.EventData)
	err = es.PostEvent(event)

	if err != nil {
		slog.Error("Error broadcasting event", "error", err)
		return err
	}

	slog.Info("EventBroadcaster", "Save", "Event published")

	return nil
}

func (es *EventServiceImpl[T, S]) Delete(ctx context.Context, e T) error {
	slog.Info("EventService", "Delete", e)

	ent, err := json.Marshal(e)

	if err != nil {
		slog.Error("Error marshalling event", "error", err)
	}

	// Broadcast event

	event := Event{
		EventId:      "1",
		EventType:    es.Prefix + Deleted,
		EventTime:    time.Now(),
		EventVersion: Version,
		EventData:    string(ent),
	}

	slog.Info("EventBroadcaster", "Delete", event.EventData)
	err = es.PostEvent(event)

	if err != nil {
		slog.Error("Error broadcasting event", "error", err)
		return err
	}

	slog.Info("EventBroadcaster", "Delete", "Event published")

	return nil
}

func (es *EventServiceImpl[T, S]) Exists(ctx context.Context, ID S) (bool, error) {
	return es.Repository.Exists(ctx, ID)
}

func (es *EventServiceImpl[T, S]) GetId(ctx context.Context, e T) (S, error) {
	return es.Repository.GetId(ctx, e)
}

func (es *EventServiceImpl[T, S]) Get(ctx context.Context, ID S) (T, error) {
	return es.Repository.Get(ctx, ID)

}

func (es *EventServiceImpl[T, S]) GetAll(ctx context.Context) ([]T, error) {
	return es.Repository.GetAll(ctx)
}

func decodeEntity[T any](data string) T {
	var wl T
	err := json.Unmarshal([]byte(data), &wl)
	if err != nil {
		log.Fatalf("Error decoding work log: %v", err)
	}
	return wl
}

func (es *EventServiceImpl[T, S]) HandleEvent(event Event) error {
	handler, ok := es.Handlers[event.EventType]
	if !ok {
		return nil
	}
	return handler(event)
}
func NewEventService[T any, S comparable](repo Repository[T, S], transport Transport, prefix string) EventService[T, S] {

	es := EventServiceImpl[T, S]{
		Repository: repo,
		Transport:  transport,
		Prefix:     prefix,
		Handlers:   make(EventHandlers),
	}
	handlers := make(EventHandlers)

	handlers[prefix+Created] = func(event Event) error {

		var e T = decodeEntity[T](event.EventData)

		id, err := repo.GetId(context.Background(), e)

		if err != nil {
			slog.Error("Error getting ID", "error", err)
		}
		exists, err := repo.Exists(context.Background(), id)
		if err != nil {
			slog.Error("Error checking if ID exists", "error", err)
		}

		if exists {
			slog.Info("EventService", "Create", "Entity already exists")
			return nil
		}

		repo.Create(context.Background(), e)
		return nil
	}

	handlers[prefix+Updated] = func(event Event) error {

		var e T = decodeEntity[T](event.EventData)

		id, err := repo.GetId(context.Background(), e)

		if err != nil {
			slog.Error("Error getting ID", "error", err)
		}
		exists, err := repo.Exists(context.Background(), id)

		if err != nil {
			slog.Error("Error checking if entity exists", "error", err)
		}

		if !exists {
			slog.Info("EventService", "Update", "Entity does not exist")
			return nil
		}

		repo.Save(context.Background(), e)
		return nil
	}

	handlers[prefix+Deleted] = func(event Event) error {

		var e T = decodeEntity[T](event.EventData)

		repo.Delete(context.Background(), e)
		return nil
	}

	es.SetHandlers(handlers)

	return &es
}
