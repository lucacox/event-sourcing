package backend

import (
	"github.com/lucacox/event-sourcing/registry"
	"github.com/stretchr/testify/mock"
)

type MockBackend struct {
	mock.Mock
}

func (m *MockBackend) Connect() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockBackend) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockBackend) SetEventRegistry(er *registry.EventRegistry) {
	m.Called(er)
}

func (m *MockBackend) Setup(name string, replicationFactor int) error {
	args := m.Called(name, replicationFactor)
	return args.Error(0)
}

func (m *MockBackend) Save(events []*registry.Event, expectedSeq uint64) (uint64, error) {
	args := m.Called(events, expectedSeq)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockBackend) Load() (map[string]*registry.Event, error) {
	args := m.Called()
	return args.Get(0).(map[string]*registry.Event), args.Error(1)
}

func (m *MockBackend) LoadByEntityId(id string) ([]*registry.Event, error) {
	args := m.Called(id)
	return args.Get(0).([]*registry.Event), args.Error(1)
}

func (m *MockBackend) LoadByEventType(eventType string) ([]*registry.Event, error) {
	args := m.Called(eventType)
	return args.Get(0).([]*registry.Event), args.Error(1)
}
