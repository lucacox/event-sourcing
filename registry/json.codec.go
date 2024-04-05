package registry

import (
	"encoding/json"

	"github.com/lucacox/event-sourcing/keystore"
)

type JsonCodec struct {
	ks keystore.KeyStore
}

func NewJsonCodec(ks keystore.KeyStore) *JsonCodec {
	return &JsonCodec{ks: ks}
}

func (jc *JsonCodec) Name() string {
	return "JSON Codec"
}

func (jc *JsonCodec) Decode(data []byte, target *Event) error {
	err := json.Unmarshal(data, &target)
	if jc.ks != nil {
		key, err := jc.ks.GetKey(target.EntityId)
		if err == nil && key != nil {
			err = target.DecryptPayload(key)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (jc *JsonCodec) Encode(e *Event) ([]byte, error) {
	if jc.ks != nil {
		key, err := jc.ks.GetKey(e.EntityId)
		if err == nil && key != nil {
			err = e.EncryptPayload(key)
			if err != nil {
				return nil, err
			}
		}
	}
	return json.Marshal(e)
}
