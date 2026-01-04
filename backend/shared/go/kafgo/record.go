package kafgo

import (
	"context"
	"errors"
	"fmt"
	tele "social-network/shared/go/telemetry"
	"sync/atomic"

	"github.com/twmb/franz-go/pkg/kgo"
)

// Record is a type that helps with commiting after processing a record
// AFTER PROCESSING THE RECORD MAKE SURE TO COMMIT!!!
type Record struct {
	monotinicId    uint64
	rec            *kgo.Record
	commitChannel  chan<- (*Record)
	confirmChannel chan (struct{})
	Context        context.Context //HANDLER OF RECORDS MUST LISTEN TO CONTEXT AND IF IT EXPIRES ROLLBACK THE TRANSACTION
}

var ErrBadArgs = errors.New("bro, you passed bad arguments")

// newRecord creates a new Record instance
func newRecord(ctx context.Context, record *kgo.Record, commitChannel chan<- (*Record), monotonicId uint64) (*Record, error) {
	if record == nil {
		tele.Error(ctx, "new record")
		return nil, fmt.Errorf("%w record: %v", ErrBadArgs, record)
	}
	return &Record{
		rec:            record,
		commitChannel:  commitChannel,
		confirmChannel: make(chan struct{}),
		Context:        ctx,
		monotinicId:    monotonicId,
	}, nil
}

// Data returns the data payload
func (r *Record) Data() []byte {
	if r.rec == nil {
		//log?
		return []byte("no data found")
	}
	return r.rec.Value
}

var ErrEmptyRecord = errors.New("empty record")

// TODO commit must return a confirmation
var a atomic.Int64

// Commit marks the record as processed in the Kafka client.
// MAKE SURE THIS IS AT THE END OF A TRANSACTION, DONT BE COMMITING THINGS YOU LATER UNDO!!
func (r *Record) Commit(ctx context.Context) error {
	if r.rec == nil {
		tele.Error(ctx, "record commit record")
		return ErrEmptyRecord
	}
	select {
	case r.commitChannel <- r:
	case <-ctx.Done():
		// optionally log or ignore
		tele.Error(ctx, "record context done")
	}

	a.Add(1)
	tele.Info(ctx, "pre  confirmation of @1, others waiting: @2", "offset", r.rec.Offset, "count", a.Load())
	//wait for the commit routine to confirm this record
	<-r.confirmChannel
	tele.Info(ctx, "post confirmation of @1, others waiting: @2", "offset", r.rec.Offset, "count", a.Load())
	a.Add(-1)

	return nil
}
