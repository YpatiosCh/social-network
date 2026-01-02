package kafgo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"social-network/shared/go/ct"

	"github.com/twmb/franz-go/pkg/kgo"
)

type kafkaProducer struct {
	client *kgo.Client
}

// seeds are used for finding the server, just as many kafka ip's you have
func NewKafkaProducer(ctx context.Context, seeds []string) (*kafkaProducer, func(), error) {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
	)
	cl.Close()

	if err != nil {
		return nil, cl.Close, err
	}

	kfc := &kafkaProducer{
		client: cl,
	}

	return kfc, cl.Close, nil
}

var ErrMarshalFail = errors.New("can't unmarshal")
var ErrProduceFail = errors.New("failed to produce")

func (kfc *kafkaProducer) Send(ctx context.Context, topic string, payload ...any) error {
	records := make([]*kgo.Record, len(payload))
	for i, p := range payload {
		bytes, err := json.Marshal(p)
		if err != nil {
			return errors.Join(ErrMarshalFail, err)
		}
		records[i] = &kgo.Record{Topic: topic, Value: bytes}
	}

	results := kfc.client.ProduceSync(ctx, records...)
	if results.FirstErr() != nil {
		return errors.Join(ErrProduceFail, results.FirstErr())
	}
	return nil
}

type kafkaConsumer struct {
	client *kgo.Client
	topics []ct.KafkaTopic
}

// seeds are used for finding the server, just as many kafka ip's you have
// enter the topics you want to consume, if any
// enter you group identifier
func NewKafkaConsumer(ctx context.Context, seeds []string, topics []string, group string) (*kafkaConsumer, func(), error) {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topics...),
	)
	cl.Close()

	if err != nil {
		return nil, cl.Close, err
	}

	kfc := &kafkaConsumer{
		client: cl,
	}

	return kfc, cl.Close, nil
}

var ErrFetch = errors.New("error when fetching")
var ErrConsumerFunc = errors.New("consumer function error")

type Record struct {
	rec           *kgo.Record
	commitChannel chan<- (*kgo.Record)
}

var ErrBadArgs = errors.New("bad arguments passed")

func newRecord(record *kgo.Record, commitChannel chan<- (*kgo.Record)) (*Record, error) {
	if record == nil {
		return nil, fmt.Errorf("%w record: %v", ErrBadArgs, record)
	}
	return &Record{
		rec:           record,
		commitChannel: commitChannel,
	}, nil
}

// String returns a human-readable representation of the record.
func (r *Record) Data() []byte {
	if r.rec == nil {
		//log?
		return []byte{}
	}
	return r.rec.Value
}

// Commit marks the record as processed in the Kafka client.
func (r *Record) Commit(ctx context.Context) {
	if r.rec == nil {
		return
	}

	r.commitChannel <- r.rec
}

func (kfc *kafkaConsumer) NewConsumer(ctx context.Context, topic ct.KafkaTopic) (<-chan *Record, error) {
	outputChan := make(chan *Record)
	return outputChan, nil
}

func (kfc *kafkaConsumer) StartConsuming(ctx context.Context) (<-chan (string), func(), error) {
	output := make(chan (string))

	go func() {
		for {
			fetches := kfc.client.PollFetches(ctx)
			if errs := fetches.Errors(); len(errs) > 0 {
				// All errors are retried internally when fetching, but non-retriable errors are
				// returned from polls so that users can notice and take action.
				close(output)
			}
			kfc.client.CommitRecords()
			// We can iterate through a record iterator...
			iter := fetches.RecordIter()
			for !iter.Done() {
				record := iter.Next()
				record.
					fmt.Println(string(record.Value), "from an iterator!")
				if err != nil {
					return errors.Join(ErrConsumerFunc, err)
				}
			}
		}
	}()
}
