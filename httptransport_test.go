package common_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	common "github.com/papawattu/cleanlog-common"
)

func TestHttpTransportPostEvent(t *testing.T) {
	t.Log("Testing HttpTransport")

	var event common.Event

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			ev := r.Body

			err := json.NewDecoder(ev).Decode(&event)

			t.Logf("Event: %v", event)

			if err != nil {
				t.Fatalf("Error decoding event: %s : %v", ev, err)
			}
			w.WriteHeader(201)
		}
	}))

	defer server.Close()

	// Create a new HttpTransport
	ht := common.NewHttpTransport(server.URL, "http://localhost:8080", 1)

	// Create a new Event

	tm := time.Now()

	ht.PostEvent(common.Event{
		EventType:    "test",
		EventTime:    tm,
		EventVersion: 1,
		EventData:    "event",
		EventId:      "1",
	})

	if event.EventId != "1" {
		t.Errorf("Event id is not correct: %s", event.EventId)
	}
	if event.EventData != "event" {
		t.Errorf("Event data is not correct: %s", event.EventData)
	}

	if event.EventVersion != 1 {
		t.Errorf("Event version is not correct: %d", event.EventVersion)
	}

	if event.EventType != "test" {
		t.Errorf("Event type is not correct: %s", event.EventType)
	}

	t.Log("HttpTransport test complete")

}

func TestHttpTransportNextEvent(t *testing.T) {
	t.Log("Testing HttpTransport")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "text/event-stream")
			w.Write([]byte("event: test\n"))
			w.Write([]byte("data: { \"EventType\":\"test\", \"EventId\" :\"1\", \"EventVersion\": 1, \"EventData\": \"test\" }\n"))
			w.Write([]byte("id: 1\n"))
			w.Write([]byte("\n"))
		}
	}))

	defer server.Close()

	// Create a new HttpTransport
	ht := common.NewHttpTransport("", server.URL, 1)

	// Create a new Event

	err := ht.Connect(context.Background())

	if err != nil {
		t.Fatalf("Error connecting: %v", err)
	}

	e, err := ht.NextEvent()

	if err != nil {
		t.Fatalf("Error getting event: %v", err)
	}

	if e == nil {
		t.Fatalf("Event is nil")
	}

	if e.EventId != "1" {
		t.Errorf("Event id is not correct: %s", e.EventId)
	}
	if e.EventData != "test" {
		t.Errorf("Event data is not correct: %s", e.EventData)
	}

	t.Log("HttpTransport test complete")

}
