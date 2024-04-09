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
Returns a registered EventType by its unique name.

#### `func (er *EventRegistry) NewEvent(name string) *Event`
Create a new isntance of the specified event calling `EventType.Init()` function.

#### `func (er *EventRegistry) RegisterCodec(codec Codec)`
Register a Codec instance in the registry.

#### `func (er *EventRegistry) GetCodec(name string) Codec`
Returns a registered codec instance.

--- 

### EventType

#### `func NewEventType(name string, codecName string, init func() *Event) *EventType`
This is the constructor for an EventType. `name` is the unique name of the event type, 
`codecName` is the name of the Codec to be associated with this type and `init` is the
event instance initialization function, used to set event default values.

---

### Event

#### `func (e *Event) EncryptPayload(key []byte) error`
This function will encrypt the event payload using AES256-CBC with the specified `key`

#### `func (e *Event) DecryptPayload(key []byte) error`
This function will decrypt the event payload using AES256-CBC with the specified `key`

#### `func (e *Event) Serialize() ([]byte, error)`
This function will serialize the event according to the Codec associated with its EventType.

#### `func (e *Event) Deserialize(data []byte) error`
This function will deserialize the `data` into the event according to the Codec associated with its EventType.

---

### Codec

This interface define methods for a Codec, used to serialize and deserialize events.

#### JSONCodec

#### `func NewJsonCodec(ks keystore.KeyStore) *JsonCodec`
JSON Codec constructor, if `ks` is not nil it will be used to get AES Keys to encrypt
events payload.

#### `func (jc *JsonCodec) Name() string`
Returns the name of the codec: "JSON Codec".

#### `func (jc *JsonCodec) Decode(data []byte, target *Event) error`
Deserialize `data` into `target`.

#### `func (jc *JsonCodec) Encode(e *Event) ([]byte, error)`
Serialize `e` into a byte array.

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
