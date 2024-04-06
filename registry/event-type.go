package registry

type EventType struct {
	Name      string
	CodecName string
	Init      func() *Event
}

func NewEventType(name string, codecName string, init func() *Event) *EventType {
	return &EventType{Name: name, CodecName: codecName, Init: init}
}
