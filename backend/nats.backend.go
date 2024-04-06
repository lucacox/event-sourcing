package backend

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lucacox/event-sourcing/registry"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NATSBackendConfig struct {
	Connection      string
	Token           string
	DefaultReplicas int
}

type NATSBackend struct {
	opts      NATSBackendConfig
	storeName string

	er *registry.EventRegistry

	nc     *nats.Conn
	js     jetstream.JetStream
	stream jetstream.Stream
}

func NewNATSBackend(opt NATSBackendConfig) *NATSBackend {
	return &NATSBackend{opts: opt}
}

func (n *NATSBackend) Connect() error {
	opts := []nats.Option{}
	if n.opts.Token != "" {
		opts = append(opts, nats.Token(n.opts.Token))
	}
	nc, err := nats.Connect(n.opts.Connection, opts...)
	if err != nil {
		return err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return err
	}

	n.nc = nc
	n.js = js

	return nil
}

func (n *NATSBackend) SetEventRegistry(er *registry.EventRegistry) {
	n.er = er
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
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	n.stream, err = n.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     storeName,
		Subjects: []string{storeName + ".>"},
		Replicas: replicas,
	})
	return err
}

func (n *NATSBackend) Save(events []*registry.Event, expectedSequence uint64) (uint64, error) {
	pubOpts := []jetstream.PublishOpt{
		jetstream.WithExpectStream(n.storeName),
	}

	var ack *jetstream.PubAck
	for i, event := range events {
		subject := fmt.Sprintf("%s.%s.%s", n.storeName, event.EntityId, event.Type)
		event.Meta["nats_subject"] = subject

		if i == 0 && expectedSequence != 0 {
			pubOpts = append(pubOpts, jetstream.WithExpectLastSequence(expectedSequence))
		}

		data, err := event.Serialize()
		if err != nil {
			return 0, err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		// ack, err = n.js.Publish(ctx, subject, data, pubOpts...)
		ack, err = n.js.PublishMsg(ctx, &nats.Msg{
			Subject: subject,
			Data:    data,
			Header: nats.Header{
				"Event-Type": []string{event.Type},
				"Codec":      []string{event.Registry.GetType(event.Type).CodecName},
			},
		}, pubOpts...)
		if err != nil {
			if strings.Contains(err.Error(), "wrong last sequence") {
				msg, _ := n.stream.GetLastMsgForSubject(ctx, subject)
				return 0, &ErrWrongSequence{
					Expected: expectedSequence,
					Actual:   msg.Sequence,
				}
			}
			return 0, err
		}
		event.Meta["nats_stream_seq"] = fmt.Sprintf("%d", ack.Sequence)
	}
	return ack.Sequence, nil
}

func (n *NATSBackend) Load() (map[string]*registry.Event, error) {
	return nil, nil
}

func (n *NATSBackend) LoadByEntityId(id string) ([]*registry.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	subject := fmt.Sprintf("%s.%s.>", n.storeName, id)

	c, err := n.stream.OrderedConsumer(ctx, jetstream.OrderedConsumerConfig{
		DeliverPolicy:  jetstream.DeliverAllPolicy,
		FilterSubjects: []string{subject},
	})
	defer n.stream.DeleteConsumer(ctx, c.CachedInfo().Name)

	if err != nil {
		return nil, err
	}
	info, err := n.stream.Info(ctx)
	if err != nil {
		return nil, err
	}
	num := info.State.Msgs
	if num == 0 {
		return []*registry.Event{}, nil
	}
	msgs, err := c.Fetch(int(num))
	if err != nil {
		return nil, err
	}
	events := make([]*registry.Event, 0, num)
	for msg := range msgs.Messages() {
		etype := msg.Headers().Get("Event-Type")
		codecName := msg.Headers().Get("Codec")

		event := n.er.NewEvent(etype)
		n.er.GetCodec(codecName).Decode(msg.Data(), event)
		if err != nil {
			return nil, err
		}
		md, err := msg.Metadata()
		if err != nil {
			return nil, err
		}
		event.Sequence = md.Sequence.Stream
		events = append(events, event)
	}
	return events, nil
}

func (n *NATSBackend) LoadByEventType(evType string) ([]*registry.Event, error) {
	return nil, nil
}
