package registry

type Codec interface {
	Name() string
	Encode(*Event) ([]byte, error)
	Decode([]byte, *Event) error
}
