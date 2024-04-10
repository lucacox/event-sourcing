package store

import (
	"github.com/lucacox/event-sourcing/registry"
	"github.com/stretchr/testify/mock"
)

type MockEntity struct {
	mock.Mock
}

func (m *MockEntity) Id() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockEntity) Project(e *registry.Event) error {
	args := m.Called(e)
	return args.Error(0)
}
