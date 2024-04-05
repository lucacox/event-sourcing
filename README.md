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
  Token:           "nat-server-token",
  DefaultReplicas: 1,
})
err := natsBE.Connect()
```

Then an Event Registry:

```go
er := registry.NewEventRegistry()
```

and optionally a Key Store:

```go
ks := keystore.NewMemoryKeyStore()
```

Finally create an EventStore for "test" entity:

```go
es := NewEventStore("test", natsBE, er, 0)
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

er.Register(registry.NewEventType("my-event", registry.NewJsonCodec(ks), func() *registry.Event {
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
evt := es.NewEvent("my-event")
evt.Payload.Field1 = "test"
// or 
evt.Payload = MyEventPayload{
  Field1: "field-1",
  Field2: 10,
}
```

For a full example check the `example` directory.

## API

### EventRegistry

TODO

### Codec

TODO

### KeyStore

TODO

### Backend

TODO

### EventStore

TODO

