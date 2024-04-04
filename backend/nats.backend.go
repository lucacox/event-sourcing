package backend

import (
	"fmt"
	"strings"

	"github.com/lucacox/event-sourcing/registry"
	"github.com/nats-io/nats.go"
)

type NATSBackendConfig struct {
	Connection      string
	Token           string
	DefaultReplicas int
}

type NATSBackend struct {
	opts      NATSBackendConfig
	storeName string

	nc *nats.Conn
	js nats.JetStreamContext
}

func NewNATSBackend(opt NATSBackendConfig) *NATSBackend {
	return &NATSBackend{opts: opt}
}

func (n *NATSBackend) Connect() error {
	nc, err := nats.Connect(n.opts.Connection, nats.Token(n.opts.Token))
	if err != nil {
		return err
	}

	js, err := nc.JetStream()
	if err != nil {
		return err
	}

	n.nc = nc
	n.js = js

	return nil
}

func (n *NATSBackend) Close() error {
	n.nc.Close()
	n.nc = nil
	n.js = nil
	return nil
}

func (n *NATSBackend) Setup(storeName string, replicas int) error {
	if replicas == 0 {
		replicas = n.opts.DefaultReplicas
	}
	n.storeName = storeName
	_, err := n.js.AddStream(&nats.StreamConfig{
		Name:     storeName,
		Subjects: []string{storeName + ".>"},
		Replicas: replicas,
	})
	return err
}

func (n *NATSBackend) Save(events []*registry.Event, expectedSequence uint64) (uint64, error) {
	pubOpts := []nats.PubOpt{
		nats.ExpectStream(n.storeName),
	}

	var ack *nats.PubAck
	for i, event := range events {
		subject := fmt.Sprintf("%s.%s.%s", n.storeName, event.EntityId, event.Type)
		event.Meta["nats_subject"] = subject

		if i == 0 && expectedSequence != 0 {
			pubOpts = append(pubOpts, nats.ExpectLastSequencePerSubject(expectedSequence))
		}

		data, err := event.Serialize()
		if err != nil {
			return 0, err
		}
		ack, err = n.js.Publish(subject, data, pubOpts...)
		if err != nil {
			if strings.Contains(err.Error(), "wrong last sequence") {
				return 0, &ErrWrongSequence{
					Expected: expectedSequence,
					Actual:   ack.Sequence,
				}
			}
			return 0, err
		}
		event.Meta["nats_stream_seq"] = fmt.Sprintf("%d", ack.Sequence)
	}
	return ack.Sequence, nil
}

func (n *NATSBackend) Load() ([]*registry.Event, error) {
	return nil, nil
}
