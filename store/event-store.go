package store

import (
	"github.com/lucacox/event-sourcing/backend"
	"github.com/lucacox/event-sourcing/registry"
)

type EventStore struct {
	name string
	be   backend.Backend
	er   *registry.EventRegistry
	rf   int
}

func NewEventStore(name string, be backend.Backend, er *registry.EventRegistry, replicationFactor int) *EventStore {
	if be == nil || er == nil {
		return nil
	}
	be.SetEventRegistry(er)
	return &EventStore{name: name, be: be, er: er, rf: replicationFactor}
}

// Start connects to the backend and sets up the store
func (es *EventStore) Start() error {
	err := es.be.Connect()
	if err != nil {
		return err
	}
	return es.be.Setup(es.name, es.rf)
}

// Stop closes the connection to the backend
func (es *EventStore) Stop() error {
	if es != nil && es.be != nil {
		return es.be.Close()
	}
	return nil
}

// NewEvent creates a new event with the given name
func (es *EventStore) NewEvent(name string) *registry.Event {
	e := es.er.NewEvent(name)
	e.Registry = es.er // this is ugly
	return e
}

// AddEvent synchronously adds a new event to the store and returns the sequence number
func (es *EventStore) AddEvent(e *registry.Event, expectedSeq uint64) (uint64, error) {
	return es.be.Save([]*registry.Event{e}, expectedSeq)
}

// Project applies all events for a given entity to the model
// and returns the last sequence number applied and the error if any
func (es *EventStore) Project(model Entity) (uint64, error) {
	events, err := es.be.LoadByEntityId(model.Id())
	if err != nil {
		return 0, err
	}

	var lastSeq uint64
	for _, e := range events {
		err = model.Project(e)
		if err != nil {
			return e.Sequence, err
		}
		lastSeq = e.Sequence
	}

	return lastSeq, nil
}

// ProjectAll applies all events for a given list of entities to the models
// and returns a map of the last sequence number applied and a map of errors if any
func (es *EventStore) ProjectAll(models []Entity) (map[string]uint64, map[string]error) {
	projections := make(map[string]uint64)
	errors := make(map[string]error)
	for _, m := range models {
		seq, err := es.Project(m)
		errors[m.Id()] = err
		projections[m.Id()] = seq
	}
	return projections, errors
}
