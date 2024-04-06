package registry

import (
	"time"
)

type Event struct {
	Id        string            `json:"id"`
	Timestamp time.Time         `json:"timestamp"`
	EntityId  string            `json:"-"`
	Type      string            `json:"type"`
	Payload   any               `json:"payload"`
	Meta      map[string]string `json:"meta"`
	Sequence  uint64            `json:"-"`
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
	evtType := e.Registry.GetType(e.Type)
	codec := e.Registry.GetCodec(evtType.CodecName)
	return codec.Encode(e)
}

func (e *Event) Deserialize(data []byte) error {
	evtType := e.Registry.GetType(e.Type)
	codec := e.Registry.GetCodec(evtType.CodecName)
	return codec.Decode(data, e)
}
