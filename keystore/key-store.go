package keystore

type KeyStore interface {
	GetKey(string) ([]byte, error)
	SetKey(string, []byte) error
	DeleteKey(string) error
}
