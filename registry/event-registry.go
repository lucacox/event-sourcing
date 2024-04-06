package registry

import "github.com/google/uuid"

type EventRegistry struct {
	types  map[string]*EventType
	codecs map[string]Codec
}

func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		types:  make(map[string]*EventType),
		codecs: make(map[string]Codec),
	}
}

func (er *EventRegistry) Register(et *EventType) {
	er.types[et.Name] = et
}

func (er *EventRegistry) GetType(name string) *EventType {
	return er.types[name]
}

func (er *EventRegistry) NewEvent(name string) *Event {
	et := er.GetType(name)
	if et == nil {
		return nil
	}
	evt := et.Init()
	evt.Id = uuid.New().String()
	evt.Registry = er
	return evt
}

func (er *EventRegistry) RegisterCodec(codec Codec) {
	er.codecs[codec.Name()] = codec
}

func (er *EventRegistry) GetCodec(name string) Codec {
	return er.codecs[name]
}
