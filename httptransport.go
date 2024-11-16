package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HttpTransport struct {
	lastId         string
	postUri        string
	streamUri      string
	scanner        *bufio.Scanner
	defaultRetries int
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

		e := ht.scanner.Text()
		switch {
		case strings.HasPrefix(e, "event: "):
			log.Printf("Event: %s\n", strings.TrimLeft(e, "event: "))
		case strings.HasPrefix(e, "data: "):
			log.Printf("Data: %s\n", strings.TrimLeft(e, "data: "))
			event = decodeEvent(strings.TrimLeft(e, "data: "))
		case strings.HasPrefix(e, "id: "):
			log.Printf("Id: %s\n", strings.TrimLeft(e, "id: "))
			ht.lastId = strings.TrimLeft(e, "id: ")
		default:
			break
		}
	}
	return &event, nil
}

func (ht *HttpTransport) Connect() error {

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

	//defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: status code %d", resp.StatusCode)
	}

	ht.scanner = bufio.NewScanner(resp.Body)

	return nil
}
func NewHttpTransport(postUri, streamUri string, retries int) *HttpTransport {

	return &HttpTransport{
		postUri:        postUri,
		streamUri:      streamUri,
		defaultRetries: retries,
	}
}
