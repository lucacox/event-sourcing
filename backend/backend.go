package backend

import (
	"fmt"

	"github.com/lucacox/event-sourcing/registry"
)

type Backend interface {
	Connect() error
	Close() error
	Setup(storeName string, replicas int) error
	Save(events []*registry.Event, expectedSequence uint64) (uint64, error)
	Load() ([]*registry.Event, error)
}

type ErrWrongSequence struct {
	Expected uint64
	Actual   uint64
}

func (e *ErrWrongSequence) Error() string {
	return fmt.Sprintf("wrong sequence: expected %d, got %d", e.Expected, e.Actual)
}
