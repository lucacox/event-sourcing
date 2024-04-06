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

type UpdateDeviceEventPayload struct {
	DeviceId     string   `json:"device_id"`
	Serial       string   `json:"serial"`
	MACAddresses []string `json:"mac_addresses"`
	UpdatedAt    string   `json:"updated_at"`
}

type Device struct {
	DeviceId     string
	Serial       string
	MACAddresses []string
	CreatedAt    time.Time
	UpadatedAt   time.Time
}

func (d *Device) Id() string {
	return d.DeviceId
}

func (d *Device) Project(e *registry.Event) error {
	if e.Type == "new-device" {
		payload := e.Payload.(map[string]interface{})
		d.DeviceId = payload["device_id"].(string)
		d.Serial = payload["serial"].(string)
		macAddresses := payload["mac_addresses"].([]interface{})
		d.MACAddresses = make([]string, len(macAddresses))
		for i, addr := range macAddresses {
			d.MACAddresses[i] = addr.(string)
		}
		d.CreatedAt, _ = time.Parse(time.RFC3339, payload["created_at"].(string))
		d.UpadatedAt, _ = time.Parse(time.RFC3339, payload["updated_at"].(string))
	}

	if e.Type == "update-device" {
		payload := e.Payload.(map[string]interface{})
		d.DeviceId = payload["device_id"].(string)
		d.Serial = payload["serial"].(string)
		macAddresses := payload["mac_addresses"].([]interface{})
		d.MACAddresses = make([]string, len(macAddresses))
		for i, addr := range macAddresses {
			d.MACAddresses[i] = addr.(string)
		}
		d.UpadatedAt, _ = time.Parse(time.RFC3339, payload["updated_at"].(string))
	}

	return nil
}

func main() {
	ks := keystore.NewMemoryKeyStore()
	er := registry.NewEventRegistry()
	natsBE := backend.NewNATSBackend(backend.NATSBackendConfig{
		Connection:      "nats://localhost:4222",
		DefaultReplicas: 1,
	})
	es := store.NewEventStore("devices", natsBE, er, 0)
	err := es.Start()
	if err != nil {
		panic(err)
	}
	defer es.Stop()

	jsonCodec := registry.NewJsonCodec(ks)
	er.RegisterCodec(jsonCodec)

	newDeviceEventType := registry.NewEventType("new-device", jsonCodec.Name(), func() *registry.Event {
		return &registry.Event{
			Type:      "new-device",
			Timestamp: time.Now(),
			Meta:      map[string]string{},
			Payload:   NewDeviceEventPayload{},
		}
	})
	er.Register(newDeviceEventType)

	updateDeviceEventType := registry.NewEventType("update-device", jsonCodec.Name(), func() *registry.Event {
		return &registry.Event{
			Type:      "update-device",
			Timestamp: time.Now(),
			Meta:      map[string]string{},
			Payload:   UpdateDeviceEventPayload{},
		}
	})
	er.Register(updateDeviceEventType)

	evt := er.NewEvent("new-device")
	evt.EntityId = "123"
	evt.Payload = NewDeviceEventPayload{
		DeviceId:     "123",
		Serial:       "123456",
		MACAddresses: []string{"00:00:00:00:00:00"},
		CreatedAt:    "2020-01-01T00:00:00Z",
		UpadatedAt:   "2020-01-01T00:00:00Z",
	}

	seq, err := es.AddEvent(evt, 0)
	if err != nil {
		panic(err)
	}

	evt1 := es.NewEvent("update-device")
	evt1.EntityId = "123"
	evt1.Payload = UpdateDeviceEventPayload{
		DeviceId:     "123",
		Serial:       "123456",
		MACAddresses: []string{"00:00:00:00:00:00", "11:11:11:11:11:11"},
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}

	_, err = es.AddEvent(evt1, seq)
	if err != nil {
		panic(err)
	}

	device := &Device{
		DeviceId: "123",
	}
	_, err = es.Project(device)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Device: %+v\n", device)
}
