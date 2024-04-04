package main

import (
	"github.com/lucacox/event-sourcing/backend"
	"github.com/lucacox/event-sourcing/registry"
)

type EventStore struct {
	name string
	rs   backend.Backend
	er   *registry.EventRegistry
}

func NewEventStore(name string, rs backend.Backend, er *registry.EventRegistry, replicationFactor int) *EventStore {
	rs.Setup(name, replicationFactor)
	return &EventStore{name: name, rs: rs, er: er}
}

func (es *EventStore) NewEvent(name string) *registry.Event {
	e := es.er.NewEvent(name)
	e.Registry = es.er // this is ugly
	return e
}
