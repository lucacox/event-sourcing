package projector

import "github.com/lucacox/event-sourcing/registry"

type Projector interface {
	Project(events *registry.Event) error
}
