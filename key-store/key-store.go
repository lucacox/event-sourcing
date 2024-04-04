package keystore

type KeyStore interface {
	GetKey(key string) ([]byte, error)
	SetKey(key string, value []byte) error
	DeleteKey(key string) error
}
