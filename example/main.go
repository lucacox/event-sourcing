package main

import (
	"fmt"
	"time"

	"github.com/lucacox/event-sourcing/backend"
	"github.com/lucacox/event-sourcing/keystore"
	"github.com/lucacox/event-sourcing/registry"
	"github.com/lucacox/event-sourcing/store"
)

type NewDeviceEventPayload struct {
	DeviceId     string   `json:"device_id"`
	Serial       string   `json:"serial"`
	MACAddresses []string `json:"mac_addresses"`
	CreatedAt    string   `json:"created_at"`
	UpadatedAt   string   `json:"updated_at"`
}

func main() {
	ks := keystore.NewMemoryKeyStore()
	er := registry.NewEventRegistry()
	natsBE := backend.NewNATSBackend(backend.NATSBackendConfig{
		Connection:      "nats://localhost:4222",
		Token:           "Alifax-RMC-Secret",
		DefaultReplicas: 1,
	})
	err := natsBE.Connect()
	if err != nil {
		panic(err)
	}
	es := store.NewEventStore("test", natsBE, er, 0)
	defer natsBE.Close()

	newDeviceEventType := registry.NewEventType("new-device", registry.NewJsonCodec(ks), func() *registry.Event {
		return &registry.Event{
			Type:      "new-device",
			Timestamp: time.Now(),
			Meta:      map[string]string{},
		}
	})
	er.Register(newDeviceEventType)

	evt := es.NewEvent("new-device")
	evt.Payload = NewDeviceEventPayload{
		DeviceId:     "123",
		Serial:       "123456",
		MACAddresses: []string{"00:00:00:00:00:00"},
		CreatedAt:    "2020-01-01T00:00:00Z",
		UpadatedAt:   "2020-01-01T00:00:00Z",
	}
	evt.Meta["subject"] = "devices.123"

	data, _ := evt.Serialize()
	fmt.Printf("%s\n", data)
}
