package registry

type EventRegistry struct {
	types map[string]*EventType
}

func NewEventRegistry() *EventRegistry {
	return &EventRegistry{types: make(map[string]*EventType)}
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
	return et.Init()
}
