package registry

type EventType struct {
	Name  string
	Codec Codec
	Init  func() *Event
}

func NewEventType(name string, codec Codec, init func() *Event) *EventType {
	return &EventType{Name: name, Codec: codec, Init: init}
}
