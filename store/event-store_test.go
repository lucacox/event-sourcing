package store

import (
	"errors"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"

	"github.com/lucacox/event-sourcing/backend"
	"github.com/lucacox/event-sourcing/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewEventStore(t *testing.T) {
	convey.Convey("Given a name, backend, event registry and replication factor", t, func() {
		name := "test-store"
		be := new(backend.MockBackend)
		er := new(registry.EventRegistry)
		replicationFactor := 3

		be.On("SetEventRegistry", er).Return()

		convey.Convey("A new event store should be created", func() {
			store := NewEventStore(name, be, er, replicationFactor)
			assert.Equal(t, name, store.name, "Expected store name to be %s, but got %s", name, store.name)
			assert.Equal(t, be, store.be, "Expected backend to be equal")
			assert.Equal(t, er, store.er, "Expected event registry to be equal")
			assert.Equal(t, replicationFactor, store.rf, "Expected replication factor to be %d, but got %d", replicationFactor, store.rf)
		})
	})
}

func TestEventStore_Start(t *testing.T) {
	convey.Convey("Given a connected backend and event store", t, func() {
		name := "test-store"
		be := new(backend.MockBackend)
		er := new(registry.EventRegistry)
		replicationFactor := 3
		be.On("SetEventRegistry", er).Return()
		be.On("Connect").Return(nil)
		be.On("Setup", name, replicationFactor).Return(nil)

		store := NewEventStore(name, be, er, replicationFactor)

		convey.Convey("The event store should start successfully", func() {
			err := store.Start()
			assert.NoError(t, err, "Expected no error, but got %v", err)
		})
	})

	convey.Convey("Given a backend connection error", t, func() {
		name := "test-store"
		be := new(backend.MockBackend)
		er := new(registry.EventRegistry)
		replicationFactor := 3
		be.On("SetEventRegistry", er).Return()
		expectedErr := errors.New("backend connection error")
		be.On("Connect").Return(expectedErr)

		store := NewEventStore(name, be, er, replicationFactor)

		convey.Convey("The event store should return the connection error", func() {
			err := store.Start()
			assert.EqualError(t, err, expectedErr.Error(), "Expected error %v, but got %v", expectedErr, err)
		})
	})

	convey.Convey("Given a backend setup error", t, func() {
		name := "test-store"
		be := new(backend.MockBackend)
		er := new(registry.EventRegistry)
		replicationFactor := 3
		be.On("SetEventRegistry", er).Return()
		expectedErr := errors.New("backend setup error")
		be.On("Connect").Return(nil)
		be.On("Setup", name, replicationFactor).Return(expectedErr)

		store := NewEventStore(name, be, er, replicationFactor)

		convey.Convey("The event store should return the setup error", func() {
			err := store.Start()
			assert.EqualError(t, err, expectedErr.Error(), "Expected error %v, but got %v", expectedErr, err)
		})
	})
}

func TestEventStore_Stop(t *testing.T) {
	convey.Convey("Given a connected backend and event store", t, func() {
		name := "test-store"
		be := new(backend.MockBackend)
		er := new(registry.EventRegistry)
		replicationFactor := 3
		be.On("SetEventRegistry", er).Return()
		be.On("Connect").Return(nil)
		be.On("Setup", name, replicationFactor).Return(nil)
		be.On("Close").Return(nil)

		store := NewEventStore(name, be, er, replicationFactor)

		convey.Convey("The event store should stop successfully", func() {
			err := store.Stop()
			assert.NoError(t, err, "Expected no error, but got %v", err)
		})
	})

	convey.Convey("Given a backend close error", t, func() {
		name := "test-store"
		be := new(backend.MockBackend)
		er := registry.NewEventRegistry()
		replicationFactor := 3
		be.On("SetEventRegistry", er).Return()
		be.On("Connect").Return(nil)
		be.On("Setup", name, replicationFactor).Return(nil)
		expectedErr := errors.New("backend close error")
		be.On("Close").Return(expectedErr)

		store := NewEventStore(name, be, er, replicationFactor)

		convey.Convey("The event store should return the close error", func() {
			err := store.Stop()
			assert.EqualError(t, err, expectedErr.Error(), "Expected error %v, but got %v", expectedErr, err)
		})
	})

	convey.Convey("Given a nil backend", t, func() {
		store := NewEventStore("test-store", nil, nil, 3)

		convey.Convey("The event store should stop successfully", func() {
			err := store.Stop()
			assert.NoError(t, err, "Expected no error, but got %v", err)
		})
	})
}

func TestEventStore_NewEvent(t *testing.T) {
	convey.Convey("Given an event store", t, func() {
		name := "test-store"
		be := new(backend.MockBackend)
		er := registry.NewEventRegistry()
		replicationFactor := 3
		be.On("SetEventRegistry", er).Return()

		eventTypeName := "test-event"
		er.Register(registry.NewEventType(eventTypeName, "test-codec", func() *registry.Event {
			return &registry.Event{
				Type:      eventTypeName,
				Timestamp: time.Now(),
				Meta:      map[string]string{},
				Payload:   nil,
			}
		}))

		store := NewEventStore(name, be, er, replicationFactor)

		convey.Convey("When calling NewEvent", func() {
			event := store.NewEvent(eventTypeName)

			convey.Convey("The event should be created with the correct name and event registry", func() {
				assert.Equal(t, eventTypeName, event.Type, "Expected event name to be %s, but got %s", eventTypeName, event.Type)
				assert.Equal(t, er, event.Registry, "Expected event registry to be equal")
			})
		})
	})
}

func TestEventStore_AddEvent(t *testing.T) {
	convey.Convey("Given an event store and an event", t, func() {
		name := "test-store"
		be := new(backend.MockBackend)
		er := registry.NewEventRegistry()
		replicationFactor := 3
		be.On("SetEventRegistry", er).Return()

		eventTypeName := "test-event"
		er.Register(registry.NewEventType(eventTypeName, "test-codec", func() *registry.Event {
			return &registry.Event{
				Type:      eventTypeName,
				Timestamp: time.Now(),
				Meta:      map[string]string{},
				Payload:   nil,
			}
		}))

		store := NewEventStore(name, be, er, replicationFactor)

		event := store.NewEvent(eventTypeName)
		be.On("Save", []*registry.Event{event}, uint64(10)).Return(uint64(10), nil)

		convey.Convey("When calling AddEvent with the event and expected sequence number", func() {
			expectedSeq := uint64(10)
			seq, err := store.AddEvent(event, expectedSeq)

			convey.Convey("The event should be saved and the sequence number should be returned", func() {
				assert.NoError(t, err, "Expected no error, but got %v", err)
				assert.Equal(t, expectedSeq, seq, "Expected sequence number to be %d, but got %d", expectedSeq, seq)
			})
		})
	})
}

func TestEventStore_Project(t *testing.T) {
	convey.Convey("Given an event store and a model", t, func() {
		name := "test-store"
		be := new(backend.MockBackend)
		er := registry.NewEventRegistry()
		replicationFactor := 3
		be.On("SetEventRegistry", er).Return()

		store := NewEventStore(name, be, er, replicationFactor)

		model := new(MockEntity)
		entityID := "test-entity"
		model.On("Id").Return(entityID)
		model.On("Project", mock.Anything).Return(nil)

		convey.Convey("When calling Project", func() {
			events := []*registry.Event{
				{
					Sequence: 1,
					Payload:  []byte("event1"),
				},
				{
					Sequence: 2,
					Payload:  []byte("event2"),
				},
				{
					Sequence: 3,
					Payload:  []byte("event3"),
				},
			}
			be.On("LoadByEntityId", entityID).Return(events, nil)

			lastSeq, err := store.Project(model)

			convey.Convey("The events should be loaded and projected onto the model", func() {
				assert.NoError(t, err, "Expected no error, but got %v", err)
				assert.Equal(t, uint64(3), lastSeq, "Expected last sequence number to be 3, but got %d", lastSeq)
				model.AssertCalled(t, "Project", events[0])
				model.AssertCalled(t, "Project", events[1])
				model.AssertCalled(t, "Project", events[2])
			})
		})
	})

	convey.Convey("Given an event store with a backend error", t, func() {
		name := "test-store"
		be := new(backend.MockBackend)
		er := registry.NewEventRegistry()
		replicationFactor := 3
		be.On("SetEventRegistry", er).Return()

		store := NewEventStore(name, be, er, replicationFactor)

		model := new(MockEntity)
		entityID := "test-entity"
		model.On("Id").Return(entityID)
		model.On("Project", mock.Anything).Return(nil)

		convey.Convey("When calling Project", func() {
			expectedErr := errors.New("backend error")
			be.On("LoadByEntityId", entityID).Return([]*registry.Event{}, expectedErr)

			lastSeq, err := store.Project(model)

			convey.Convey("The event store should return the backend error", func() {
				assert.EqualError(t, err, expectedErr.Error(), "Expected error %v, but got %v", expectedErr, err)
				assert.Equal(t, uint64(0), lastSeq, "Expected last sequence number to be 0, but got %d", lastSeq)
				model.AssertNotCalled(t, "Project")
			})
		})
	})

	convey.Convey("Given a model that fails to project an event", t, func() {
		name := "test-store"
		be := new(backend.MockBackend)
		er := registry.NewEventRegistry()
		replicationFactor := 3
		be.On("SetEventRegistry", er).Return()

		store := NewEventStore(name, be, er, replicationFactor)

		model := new(MockEntity)
		entityID := "test-entity"
		model.On("Id").Return(entityID)

		convey.Convey("When calling Project", func() {
			events := []*registry.Event{
				{
					Sequence: 1,
					Payload:  []byte("event1"),
				},
				{
					Sequence: 2,
					Payload:  []byte("event2"),
				},
				{
					Sequence: 3,
					Payload:  []byte("event3"),
				},
			}
			expectedErr := errors.New("projection error")
			model.On("Project", events[0]).Return(nil).Once()
			model.On("Project", events[1]).Return(expectedErr)

			be.On("LoadByEntityId", entityID).Return(events, nil)

			lastSeq, err := store.Project(model)

			convey.Convey("The event store should return the projection error", func() {
				assert.EqualError(t, err, expectedErr.Error(), "Expected error %v, but got %v", expectedErr, err)
				assert.Equal(t, uint64(2), lastSeq, "Expected last sequence number to be 2, but got %d", lastSeq)
				model.AssertCalled(t, "Project", events[0])
				model.AssertCalled(t, "Project", events[1])
				model.AssertNotCalled(t, "Project", events[2])
			})
		})
	})
}

func TestEventStore_ProjectAll(t *testing.T) {
	convey.Convey("Given an event store and a list of models", t, func() {
		name := "test-store"
		be := new(backend.MockBackend)
		er := registry.NewEventRegistry()
		replicationFactor := 3
		be.On("SetEventRegistry", er).Return()

		store := NewEventStore(name, be, er, replicationFactor)

		model1 := new(MockEntity)
		entityID1 := "test-entity-1"
		model1.On("Id").Return(entityID1)
		model1.On("Project", mock.Anything).Return(nil)

		model2 := new(MockEntity)
		entityID2 := "test-entity-2"
		model2.On("Id").Return(entityID2)
		model2.On("Project", mock.Anything).Return(nil)

		models := []Entity{model1, model2}

		convey.Convey("When calling ProjectAll", func() {
			be.On("LoadByEntityId", entityID1).Return([]*registry.Event{{}}, nil)
			be.On("LoadByEntityId", entityID2).Return([]*registry.Event{{}}, nil)

			projections, errors := store.ProjectAll(models)

			convey.Convey("The events should be loaded and projected onto the models", func() {
				assert.Len(t, projections, 2, "Expected 2 projections")
				assert.Len(t, errors, 2, "Expected 2 errors")
				assert.NoError(t, errors[entityID1], "Expected no error for model 1")
				assert.NoError(t, errors[entityID2], "Expected no error for model 2")
				model1.AssertCalled(t, "Project", mock.Anything)
				model2.AssertCalled(t, "Project", mock.Anything)
			})
		})
	})

	convey.Convey("Given an event store with a backend error", t, func() {
		name := "test-store"
		be := new(backend.MockBackend)
		er := registry.NewEventRegistry()
		replicationFactor := 3
		be.On("SetEventRegistry", er).Return()

		store := NewEventStore(name, be, er, replicationFactor)

		model1 := new(MockEntity)
		entityID1 := "test-entity-1"
		model1.On("Id").Return(entityID1)
		model1.On("Project", mock.Anything).Return(nil)

		model2 := new(MockEntity)
		entityID2 := "test-entity-2"
		model2.On("Id").Return(entityID2)
		model2.On("Project", mock.Anything).Return(nil)

		models := []Entity{model1, model2}

		convey.Convey("When calling ProjectAll", func() {
			expectedErr := errors.New("backend error")
			be.On("LoadByEntityId", entityID1).Return([]*registry.Event{}, expectedErr)
			be.On("LoadByEntityId", entityID2).Return([]*registry.Event{}, expectedErr)

			projections, errors := store.ProjectAll(models)

			convey.Convey("The event store should return the backend error for each model", func() {
				assert.Len(t, projections, 2, "Expected 2 projections")
				assert.Len(t, errors, 2, "Expected 2 errors")
				assert.EqualError(t, errors[entityID1], expectedErr.Error(), "Expected backend error for model 1")
				assert.EqualError(t, errors[entityID2], expectedErr.Error(), "Expected backend error for model 2")
				model1.AssertNotCalled(t, "Project")
				model2.AssertNotCalled(t, "Project")
			})
		})
	})
}
