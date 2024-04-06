package backend

import (
	"fmt"

	"github.com/lucacox/event-sourcing/registry"
)

type Backend interface {
	Connect() error
	Close() error
	SetEventRegistry(*registry.EventRegistry)
	Setup(string, int) error
	Save([]*registry.Event, uint64) (uint64, error)
	// returns all events in the store in a map of entity id to events
	Load() (map[string]*registry.Event, error)
	// returns all events for a given entity id
	LoadByEntityId(string) ([]*registry.Event, error)
	// returns all events for a given event type
	LoadByEventType(string) ([]*registry.Event, error)
}

type ErrWrongSequence struct {
	Expected uint64
	Actual   uint64
}

func (e *ErrWrongSequence) Error() string {
	return fmt.Sprintf("wrong sequence: expected %d, got %d", e.Expected, e.Actual)
}
