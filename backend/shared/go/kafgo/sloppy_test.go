package kafgo

import (
	"context"
	"fmt"
	"math/rand/v2"
	tele "social-network/shared/go/telemetry"
	"time"
)

func Spam() {
	ctx := context.Background()
	//
	//
	//
	consumer, err := NewKafkaConsumer([]string{"kafka:9092"}, "test")
	if err != nil {
		tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
	}

	ch, err := consumer.RegisterTopic("test_topic")
	if err != nil {
		tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
	}

	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				tele.Error(ctx, "panic occured in consumer loop!")
			}
		}()
		for {
			tele.Info(ctx, "consume loop start")
			for record := range ch {
				record.Commit(ctx)
			}
			time.Sleep(time.Second * 1)
		}

	}()

	_, err = consumer.StartConsuming(ctx)
	if err != nil {
		tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
	}
	tele.Info(ctx, "started consuming")

	//
	//
	//
	//
	producer, _, err := NewKafkaProducer([]string{"kafka:9092"})
	if err != nil {
		tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
	}

	type X struct {
		Name string `json:"name"`
	}

	fmt.Println(producer)

	go func() {
		for i := range 1000000 {
			dur := min(rand.IntN(30), rand.IntN(30), rand.IntN(30), rand.IntN(30), rand.IntN(30))
			time.Sleep(time.Millisecond * time.Duration((1*dur)+2))
			err := producer.Send(ctx, "test_topic", X{fmt.Sprint("alex:", i, " and slept for: ", (1*dur)+2, "ms sadfl;jkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsdajkhasdl;jkfhadlshfl;kasdhlfhsda")})
			if err != nil {
				tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
			}

		}
	}()

	time.Sleep(time.Minute)
}
