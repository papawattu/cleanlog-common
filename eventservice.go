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
	Connect(context.Context) error
	PostEvent(Event) error
	NextEvent() (*Event, error)
}
type EventService[T any, S comparable] interface {
	Repository[T, S]
	Transport
	SetPrefix(prefix string)
	SetHandlers(handlers EventHandlers)
	HandleEvent(event Event) error
	StartEventRunner(ctx context.Context)
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
	slog.Info("EventService", "HandleEvent", event, "EventType", event.EventType)
	handler, ok := es.Handlers[event.EventType]
	if !ok {
		return nil
	}
	return handler(event)
}

func (es *EventServiceImpl[T, S]) StartEventRunner(ctx context.Context) {
	go func() {

		err := es.Connect(ctx)

		if err != nil {
			log.Fatal(err)
		}
		for {
			select {
			case <-ctx.Done():
				return
			default:
				ev, err := es.NextEvent()
				if err != nil {
					log.Fatal(err)
				}

				if ev != nil {
					es.HandleEvent(*ev)
				}
			}
		}
	}()
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

		slog.Info("EventService", "Create", event.EventData)

		var e T = decodeEntity[T](event.EventData)

		repo.Create(context.Background(), e)
		return nil
	}

	handlers[prefix+Updated] = func(event Event) error {

		slog.Info("EventService", "Update", event.EventData)

		var e T = decodeEntity[T](event.EventData)

		repo.Save(context.Background(), e)
		return nil
	}

	handlers[prefix+Deleted] = func(event Event) error {

		slog.Info("EventService", "Delete", event.EventData)

		var e T = decodeEntity[T](event.EventData)

		repo.Delete(context.Background(), e)
		return nil
	}

	es.SetHandlers(handlers)

	return &es
}
