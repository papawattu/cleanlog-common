package common_test

import (
	"context"
	"testing"
	"time"

	common "github.com/papawattu/cleanlog-common"
)

var (
	nextEventCalled,
	postEventCalled,
	createCalled,
	saveCalled,
	getCalled,
	getAllCalled,
	deleteCalled,
	existsCalled,
	getIdCalled int
)

var existsReturnValue = false
var getIdReturnValue int

var mockEvent *common.Event

type testTransport struct {
}

func (t *testTransport) Connect(context.Context) error {
	return nil
}

func (t *testTransport) PostEvent(e common.Event) error {
	postEventCalled++
	mockEvent = &e
	return nil
}

func (t *testTransport) NextEvent() (*common.Event, error) {
	nextEventCalled++
	return mockEvent, nil
}

type Object struct {
	Data string `json:"data"`
}

func TestEventServiceCreate(t *testing.T) {
	createCalled = 0
	postEventCalled = 0
	nextEventCalled = 0

	existsReturnValue = false

	ctx := context.Background()
	// Create a new test repository

	repo := NewTestRepository[*common.BaseEntity[string]]()

	trans := &testTransport{}
	// Create a new event service
	es := common.NewEventService(repo, trans, "test")

	//te := common.NewBaseEntity("1")
	// Create the event
	err := es.Create(ctx, &common.BaseEntity[string]{ID: "1"})
	if err != nil {
		t.Errorf("Error creating event: %v", err)
	}

	// Get the next event
	nextEvent, err := es.NextEvent()
	if err != nil {
		t.Errorf("Error getting next event: %v", err)
	}

	if nextEvent.EventType != "testCreated" {
		t.Errorf("Event type is not correct: %s", nextEvent.EventType)
	}

	// if nextEvent.EventData != "{\"ID\":1}" {
	// 	t.Errorf("Event data is not correct: %s", nextEvent.EventData)
	// }

	if nextEvent.EventVersion != 1 {
		t.Errorf("Event version is not correct: %d", nextEvent.EventVersion)
	}

	if nextEvent.EventId != "1" {
		t.Errorf("Event id is not correct: %s", nextEvent.EventId)
	}

	if postEventCalled != 1 {
		t.Errorf("PostEvent was not called")
	}

	if nextEventCalled != 1 {
		t.Errorf("NextEvent was not called")
	}

	err = es.HandleEvent(*nextEvent)

	if err != nil {
		t.Errorf("Error handling event: %v", err)
	}

	if createCalled != 1 {
		t.Errorf("Create was not called")
	}

}
func TestEventServiceDelete(t *testing.T) {
	deleteCalled = 0
	postEventCalled = 0
	nextEventCalled = 0

	ctx := context.Background()
	// Create a new test repository
	repo := NewTestRepository[*common.BaseEntity[int], int]()

	trans := &testTransport{}
	// Create a new event service
	es := common.NewEventService(repo, trans, "test")

	//te := common.NewBaseEntity(int(1))
	// Create the event
	err := es.Delete(ctx, &common.BaseEntity[int]{ID: 1})
	if err != nil {
		t.Errorf("Error deleting event: %v", err)
	}

	// Get the next event
	nextEvent, err := es.NextEvent()
	if err != nil {
		t.Errorf("Error getting next event: %v", err)
	}

	if nextEvent.EventType != "testDeleted" {
		t.Errorf("Event type is not correct: %s", nextEvent.EventType)
	}

	if nextEvent.EventVersion != 1 {
		t.Errorf("Event version is not correct: %d", nextEvent.EventVersion)
	}

	if nextEvent.EventId != "1" {
		t.Errorf("Event id is not correct: %s", nextEvent.EventId)
	}

	if postEventCalled != 1 {
		t.Errorf("PostEvent was not called")
	}

	if nextEventCalled != 1 {
		t.Errorf("NextEvent was not called")
	}

	err = es.HandleEvent(*nextEvent)

	if err != nil {
		t.Errorf("Error handling event: %v", err)
	}

	if deleteCalled != 1 {
		t.Errorf("Delete was not called")
	}

}
func TestEventServiceSave(t *testing.T) {
	postEventCalled = 0
	nextEventCalled = 0
	saveCalled = 0
	existsReturnValue = true

	ctx := context.Background()
	// Create a new test repository
	repo := NewTestRepository[*common.BaseEntity[int]]()
	trans := &testTransport{}
	// Create a new event service
	//en := common.NewBaseEntity(int(1))
	es := common.NewEventService(repo, trans, "test")

	// Create the event
	err := es.Save(ctx, &common.BaseEntity[int]{ID: 1})
	if err != nil {
		t.Errorf("Error saving event: %v", err)
	}

	// Get the next event
	nextEvent, err := es.NextEvent()
	if err != nil {
		t.Errorf("Error getting next event: %v", err)
	}

	if nextEvent.EventType != "testUpdated" {
		t.Errorf("Event type is not correct: %s", nextEvent.EventType)
	}

	if nextEvent.EventVersion != 1 {
		t.Errorf("Event version is not correct: %d", nextEvent.EventVersion)
	}

	if nextEvent.EventId != "1" {
		t.Errorf("Event id is not correct: %s", nextEvent.EventId)
	}

	if postEventCalled != 1 {
		t.Errorf("PostEvent was not called")
	}

	if nextEventCalled != 1 {
		t.Errorf("NextEvent was not called")
	}

	err = es.HandleEvent(*nextEvent)

	if err != nil {
		t.Errorf("Error handling event: %v", err)
	}

	if saveCalled != 1 {
		t.Errorf("Save was not called")
	}
}
func TestEventServiceGetAll(t *testing.T) {

	ctx := context.Background()
	// Create a new test repository
	repo := NewTestRepository[*common.BaseEntity[int]]()

	trans := &testTransport{}
	// Create a new event service
	es := common.NewEventService(repo, trans, "test")

	// Create the event
	_, err := es.GetAll(ctx)
	if err != nil {
		t.Errorf("Error getting event: %v", err)
	}
	if getAllCalled == 0 {
		t.Errorf("GetAll was not called")
	}
}
func TestEventServiceGet(t *testing.T) {
	postEventCalled = 0
	nextEventCalled = 0

	ctx := context.Background()
	// Create a new test repository
	repo := NewTestRepository[*common.BaseEntity[int]]()

	trans := &testTransport{}
	// Create a new event service
	es := common.NewEventService(repo, trans, "test")

	// Create the event
	_, err := es.Get(ctx, 1)
	if err != nil {
		t.Errorf("Error getting event: %v", err)
	}
	if getCalled == 0 {
		t.Errorf(`Get was not called`)
	}
}

func TestEventServiceExists(t *testing.T) {
	postEventCalled = 0
	nextEventCalled = 0
	existsCalled = 0

	ctx := context.Background()
	// Create a new test repository
	repo := NewTestRepository[*common.BaseEntity[int]]()

	es := common.NewEventService(repo, &testTransport{}, "test")

	_, err := es.Exists(ctx, 1)
	if err != nil {
		t.Errorf("Error getting event: %v", err)
	}

	if existsCalled == 0 {
		t.Errorf("Exists was not called")
	}
}

func TestEventServiceGetId(t *testing.T) {
	postEventCalled = 0
	nextEventCalled = 0
	getIdCalled = 0

	ctx := context.Background()
	// Create a new test repository
	repo := NewTestRepository[*common.BaseEntity[int]]()

	es := common.NewEventService(repo, &testTransport{}, "test")

	_, err := es.GetId(ctx, &common.BaseEntity[int]{ID: 1})

	if err != nil {
		t.Errorf("Error getting event: %v", err)
	}

	if getIdCalled == 0 {
		t.Errorf("GetId was not called")
	}
}

func TestEventStartEventRunner(t *testing.T) {
	postEventCalled = 0
	nextEventCalled = 0
	getIdCalled = 0
	existsCalled = 0
	createCalled = 0
	saveCalled = 0
	getCalled = 0
	getAllCalled = 0
	deleteCalled = 0

	ctx := context.Background()
	// Create a new test repository
	repo := NewTestRepository[*common.BaseEntity[int]]()

	es := common.NewEventService(repo, &testTransport{}, "test")

	newCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	es.StartEventRunner(newCtx)

	for {
		select {
		case <-newCtx.Done():
			t.Log("Event runner stopped")
			return
		case <-time.After(2 * time.Second):
		}
	}
}
