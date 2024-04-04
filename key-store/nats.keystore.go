package keystore

import "github.com/nats-io/nats.go"

type NATSKeyStore struct {
	nc *nats.Conn
}

func NewNATSKeyStore(nc *nats.Conn) *NATSKeyStore {
	return &NATSKeyStore{nc: nc}
}

func (nks *NATSKeyStore) GetKey(id string) ([]byte, error) {
	return nil, nil
}

func (nks *NATSKeyStore) SetKey(id string, value []byte) error {
	return nil
}

func (nks *NATSKeyStore) DeleteKey(id string) error {
	return nil
}
