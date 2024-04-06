package store

import "github.com/lucacox/event-sourcing/projector"

type Entity interface {
	projector.Projector
	Id() string
}
