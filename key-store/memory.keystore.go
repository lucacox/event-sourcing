package keystore

import "fmt"

type MemoryKeyStore struct {
	store map[string][]byte
}

func NewMemoryKeyStore() *MemoryKeyStore {
	return &MemoryKeyStore{store: make(map[string][]byte)}
}

func (mks *MemoryKeyStore) GetKey(id string) ([]byte, error) {
	key, ok := mks.store[id]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}
	return key, nil
}

func (mks *MemoryKeyStore) SetKey(id string, key []byte) error {
	mks.store[id] = key
	return nil
}

func (mks *MemoryKeyStore) DeleteKey(id string) error {
	delete(mks.store, id)
	return nil
}
