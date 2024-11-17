package common

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"
)

type HttpTransport struct {
	ctx            context.Context
	lastId         string
	postUri        string
	streamUri      string
	scanner        *bufio.Scanner
	defaultRetries int
	connected      bool
	body           io.ReadCloser
}

func (ht *HttpTransport) PostEvent(event Event) error {
	ev, err := json.Marshal(event)
	if err != nil {
		return err
	}

	h := sha256.New()

	h.Write([]byte(ev))

	event.EventSHA = fmt.Sprintf("%x", h.Sum(nil))

	ev, err = json.Marshal(event)

	if err != nil {
		return err
	}

	client := NewRetryableClient(10)

	r, err := http.NewRequest("POST", ht.postUri, bytes.NewBuffer(ev))

	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(r)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Error: status code %d", resp.StatusCode)
	}

	return nil
}
func decodeEvent(ev string) Event {
	var event Event
	err := json.Unmarshal([]byte(ev), &event)
	if err != nil {
		log.Fatalf("Error decoding event: %s : %v", ev, err)
	}
	return event
}
func (ht *HttpTransport) NextEvent() (*Event, error) {

	var event Event

	for ht.scanner.Scan() {

		select {
		case <-ht.ctx.Done():
			ht.body.Close()
			return nil, nil
		default:
			e := ht.scanner.Text()
			switch {
			case strings.HasPrefix(e, "event: "):
				slog.Debug("Event Type", "Type", strings.TrimLeft(e, "event: "))
			case strings.HasPrefix(e, "data: "):
				slog.Debug("Event Data", "Data", strings.TrimLeft(e, "data: "))
				event = decodeEvent(strings.TrimLeft(e, "data: "))
				return &event, nil
			case strings.HasPrefix(e, "id: "):
				slog.Debug("Event Id", "Id", strings.TrimLeft(e, "id: "))
				ht.lastId = strings.TrimLeft(e, "id: ")
			default:
				break
			}
		}
	}
	return &event, nil
}

func (ht *HttpTransport) Connect(ctx context.Context) error {

	slog.Info("Connecting to event stream", "URI", ht.streamUri)

	ht.ctx = ctx

	client := NewRetryableClient(ht.defaultRetries)

	req, err := http.NewRequest("GET", ht.streamUri, nil)

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Last-Event-ID", ht.lastId)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error connecting to event stream: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: status code %d", resp.StatusCode)
	}

	slog.Info("Connected to event stream", "URI", req.URL)

	ht.scanner = bufio.NewScanner(resp.Body)

	ht.connected = true

	ht.body = resp.Body

	return nil
}
func NewHttpTransport(postUri, streamUri string, retries int) *HttpTransport {

	return &HttpTransport{
		postUri:        postUri,
		streamUri:      streamUri,
		defaultRetries: retries,
		connected:      false,
	}
}
