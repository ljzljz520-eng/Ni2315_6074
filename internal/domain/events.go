package domain

import (
	"encoding/json"
	"time"
)

type EventType string

const (
	EventCreated   EventType = "entry.created"
	EventSubmitted EventType = "entry.submitted"
	EventReviewed  EventType = "entry.reviewed"
	EventArchived  EventType = "entry.archived"
)

type DomainEvent struct {
	ID      string            `json:"id"`
	Type    EventType         `json:"type"`
	EntryID string            `json:"entry_id"`
	Actor   string            `json:"actor"`
	At      time.Time         `json:"at"`
	Payload map[string]string `json:"payload"`
}

func (e DomainEvent) Encode() ([]byte, error) {
	return json.Marshal(e)
}

func DecodeEvent(data []byte) (DomainEvent, error) {
	var event DomainEvent
	err := json.Unmarshal(data, &event)
	return event, err
}

func NewEvent(id string, kind EventType, entryID string, actor string, payload map[string]string, at time.Time) DomainEvent {
	return DomainEvent{ID: id, Type: kind, EntryID: entryID, Actor: actor, Payload: payload, At: at}
}
