# Event Sourcing Library

**WARNING: this is a work in progress**

This go module implements a simple, configurable event store to implement Event Sourcing.

## Install

```bash
go get github.com/lucacox/event-sourcing
```

## Usage

First istantiate a Backend implementation:

```go
be := backend.NewNATSBackend(backend.NATSBackendConfig{
  Connection:      "nats://localhost:4222",
  DefaultReplicas: 1,
})
```

Then an Event Registry:

```go
er := registry.NewEventRegistry()
```

and optionally a KeyStore:

```go
ks := keystore.NewMemoryKeyStore()
```

Finally create and start an EventStore for "test" entity:

```go
es := NewEventStore("test", natsBE, er, 1)
err := es.Start() // start will call backend Connect() and Setup() methods
```

To register event types, that can be objects of any type, use `Register` method of `EventRegistry`, the function wants a `EventType` object, a `Codec` instance, and an function to create a new instance of the event with default values:

```go
type MyEventPayload struct {
  Field1 string
  Field2 int
}

// if a KeyStore is passed, when encoding/decoding the codec will check if
// a key is associated to the event associated entity, if so the payload 
// will be encrypted/decrypted.
jsonCodec := registry.NewJsonCodec(ks)

// to register a new event you must pass its name, a codec name and an 
// initialization function to set default values
er.Register(registry.NewEventType("my-event", jsonCode.Name(), func() *registry.Event {
	return &registry.Event{
      Type:      "my-event",
      Timestamp: time.Now(),
      Meta:      map[string]string{},
      Payload:   MyEventPayload{},
	}
}))
```

To create a new Event of a registered type:

```go
evt := er.NewEvent("my-event")
evt.Payload.Field1 = "test"
// or 
evt.Payload = MyEventPayload{
  Field1: "field-1",
  Field2: 10,
}
```

To commit an event to the backend use:

```go
// the second param is the expected last message stream sequence id
// if set to 0 the expectation is not used. The function returns
// the event stream sequence number
seq, err := es.AddEvent(evt, 0)
```

To reconstruct the state of an entity you have to create an object that implements the Projector and Entity interfaces

```go
type MyEntity struct {
  Id string
  Field1 string
  Field2 int
}

func (e *MyEntity) Id() string {
  return e.Id
}

func (e *MyEntity) Project(evt *registry.Event) error {
  switch evt.Type {
  case "my-event":
    payload = evt.Payload.(map[string]interface{})
    e.Field1 = payload["field1"].(string)
    ...
  }

  return nil
}

```


For a full example check the `example` directory.

## API

### EventRegistry

The EventRegostry holds all known EventTypes. Its used to create new event instances. 

#### `func NewEventRegistry() *EventRegistry`

This is the EventRegistry constructor.

#### `func (er *EventRegistry) Register(et *EventType)`

This will register a new EventType in the registry.

#### `func (er *EventRegistry) GetType(name string) *EventType`

Returns a registered EventType by its name.

#### `func (er *EventRegistry) NewEvent(name string) *Event`

Create a new isntance of the specified event calling `EventType.Init()` function.

#### `func (er *EventRegistry) RegisterCodec(codec Codec)`

Register a Codec instance in the registry.

#### `func (er *EventRegistry) GetCodec(name string) Codec`

Returns a registered codec instance.

--- 

### EventType

TODO

---

### Event

TODO

---

### Codec

TODO

#### JSONCodec

---

TODO

### KeyStore

TODO

#### MemoryKeyStore

TODO

#### 
NATSKeyStore

TODO

---

### Backend

TODO

#### NATSBackend

TODO

---

### EventsStore

TODO

---

### Projector

TODO
