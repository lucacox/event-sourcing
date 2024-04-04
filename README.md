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

