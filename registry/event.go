package registry

import "time"

type Event struct {
	Id        string            `json:"id"`
	Timestamp time.Time         `json:"timestamp"`
	EntityId  string            `json:"entity_id"`
	Type      string            `json:"type"`
	Payload   any               `json:"payload"`
	Meta      map[string]string `json:"meta"`
	Registry  *EventRegistry    `json:"-"`
}

func (e *Event) EncryptPayload(key []byte) error {
	// TODO: implement me
	return nil
}

func (e *Event) DecryptPayload(key []byte) error {
	// TODO: implement me
	return nil
}

func (e *Event) Serialize() ([]byte, error) {
	return e.Registry.GetType(e.Type).Codec.Encode(e)
}

func (e *Event) Deserialize(data []byte) error {
	return e.Registry.GetType(e.Type).Codec.Decode(data, e)
}
