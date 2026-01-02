package kafgo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	tele "social-network/shared/go/telemetry"

	"github.com/twmb/franz-go/pkg/kgo"
)

// How to use this. Create a kafka producer. User Send() to send payloads. The payload should be a struct with json tags, cause it will get marshaled.

type KafkaProducer struct {
	client *kgo.Client
}

// seeds are used for finding the server, just as many kafka ip's you have
func NewKafkaProducer(seeds []string) (producer *KafkaProducer, close func(), err error) {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, nil, err
	}
	kfc := &KafkaProducer{
		client: cl,
	}
	return kfc, cl.Close, nil
}

var ErrProduceFail = errors.New("failed to produce")

// Send sends payload(s) to the specified topic
func (kfc *KafkaProducer) Send(ctx context.Context, topic string, payload ...any) error {
	tele.Info(ctx, "sending")
	records := make([]*kgo.Record, len(payload))
	for i, p := range payload {
		bytes, err := json.Marshal(p)
		if err != nil {
			return err
		}

		records[i] = &kgo.Record{Topic: topic, Value: bytes}
	}

	tele.Info(ctx, "right before produce")

	// var wg sync.WaitGroup
	// wg.Add(1)
	// record := &kgo.Record{Topic: "test_topic", Value: []byte("bar")}
	// kfc.client.Produce(ctx, record, func(_ *kgo.Record, err error) {
	// 	defer wg.Done()
	// 	if err != nil {
	// 		fmt.Printf("record had a produce error: %v\n", err)
	// 	}

	// })
	// wg.Wait()

	results := kfc.client.ProduceSync(ctx, &kgo.Record{Topic: topic, Value: []byte("test_data")})
	tele.Info(ctx, "right after produce")
	if results.FirstErr() != nil {
		return fmt.Errorf("failed to produce %w", results.FirstErr())
	}
	tele.Info(ctx, "finished sending")
	return nil
}
