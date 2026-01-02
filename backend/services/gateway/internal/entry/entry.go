package entry

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"social-network/services/gateway/internal/handlers"
	"social-network/shared/gen-go/chat"
	"social-network/shared/gen-go/media"
	"social-network/shared/gen-go/posts"
	"social-network/shared/gen-go/users"
	configutil "social-network/shared/go/configs"
	"social-network/shared/go/ct"
	"social-network/shared/go/gorpc"
	"social-network/shared/go/kafgo"
	redis_connector "social-network/shared/go/redis"
	tele "social-network/shared/go/telemetry"
	"syscall"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

var producer = kafgo.KafkaProducer{}

func Run() {
	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	cfgs := getConfigs()

	// Inject envs to custom types
	ct.InitCustomTypes(cfgs.PassSecret, cfgs.EncrytpionKey)

	// go mytest()

	// time.Sleep(time.Second * 5)
	// os.Exit(1)

	//
	//
	//
	// TELEMETRY
	closeTelemetry, err := tele.InitTelemetry(ctx, "api-gateway", "API", cfgs.TelemetryCollectorAddress, ct.CommonKeys(), cfgs.EnableDebugLogs, cfgs.SimplePrint)
	if err != nil {
		tele.Fatalf("failed to init telemetry: %s", err.Error())
	}
	defer closeTelemetry()
	tele.Info(ctx, "initialized telemetry")

	//
	//
	//
	// CACHE
	CacheService := redis_connector.NewRedisClient(
		cfgs.RedisAddr,
		cfgs.RedisPassword,
		cfgs.RedisDB,
	)
	if err := CacheService.TestRedisConnection(); err != nil {
		// tele.Fatalf("connection test failed, ERROR: %v", err)
	}
	tele.Info(ctx, "Cache service connection started correctly")

	//
	//
	//

	consumer, err := kafgo.NewKafkaConsumer([]string{"kafka:9092"}, "test")
	if err != nil {
		tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
	}

	ch, err := consumer.RegisterTopic("test_topic")
	if err != nil {
		tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
	}

	go func() {
		for {

			for record := range ch {
				tele.Info(ctx, "RECORD ARRIVED! @1", "record", record)
				record.Commit(ctx)
			}
			time.Sleep(time.Second)
			tele.Info(ctx, "consume loop")
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
	producer, _, err := kafgo.NewKafkaProducer([]string{"kafka:9092"})
	if err != nil {
		tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
	}

	type X struct {
		Name string `json:"name"`
	}

	fmt.Println(producer)

	go func() {
		for i := range 10000 {
			err := producer.Send(ctx, "test_topic", X{fmt.Sprint("alex:", i)})
			if err != nil {
				tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
			}
			tele.Info(ctx, "producer loop")
			time.Sleep(time.Second)
		}
	}()

	time.Sleep(time.Minute)

	//
	//
	//
	// GRPC CLIENTS
	UsersService, err := gorpc.GetGRpcClient(
		users.NewUserServiceClient,
		cfgs.UsersGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		tele.Fatalf("failed to connect to users service: %v", err)
	}

	PostsService, err := gorpc.GetGRpcClient(
		posts.NewPostsServiceClient,
		cfgs.PostsGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		tele.Fatalf("failed to connect to posts service: %v", err)
	}

	ChatService, err := gorpc.GetGRpcClient(
		chat.NewChatServiceClient,
		cfgs.ChatGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		tele.Fatalf("failed to connect to chat service: %v", err)
	}

	MediaService, err := gorpc.GetGRpcClient(
		media.NewMediaServiceClient,
		cfgs.MediaGRPCAddr,
		ct.CommonKeys(),
	)
	if err != nil {
		tele.Fatalf("failed to connect to media service: %v", err)
	}

	//
	//
	//
	// HANDLER
	apiMux := handlers.NewHandlers(
		"gateway",
		CacheService,
		UsersService,
		PostsService,
		ChatService,
		MediaService,
	)

	//
	//
	//
	// SERVER
	server := &http.Server{
		Handler:     apiMux,
		Addr:        cfgs.HTTPAddr,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	srvErr := make(chan error, 1)
	go func() {
		tele.Info(ctx, "Starting server at @1", "address", server.Addr)
		srvErr <- server.ListenAndServe()
	}()

	//
	//
	//
	// SHUTDOWN
	select {
	case err = <-srvErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			tele.Fatalf("Failed to listen and serve: %v", err)
		}
	case <-ctx.Done():
		stopSignal()
	}

	tele.Info(ctx, "Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(cfgs.ShutdownTimeout)*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		tele.Fatalf("Graceful server Shutdown Failed: %v", err)
	}

	tele.Info(ctx, "Server stopped")
}

type configs struct {
	RedisAddr     string `env:"REDIS_ADDR"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB"`

	UsersGRPCAddr string `env:"USERS_GRPC_ADDR"`
	PostsGRPCAddr string `env:"POSTS_GRPC_ADDR"`
	ChatGRPCAddr  string `env:"CHAT_GRPC_ADDR"`
	MediaGRPCAddr string `env:"MEDIA_GRPC_ADDR"`

	HTTPAddr        string `env:"HTTP_ADDR"`
	ShutdownTimeout int    `env:"SHUTDOWN_TIMEOUT_SECONDS"`

	EnableDebugLogs bool `env:"ENABLE_DEBUG_LOGS"`
	SimplePrint     bool `env:"ENABLE_SIMPLE_PRINT"`

	OtelResourceAttributes    string `env:"OTEL_RESOURCE_ATTRIBUTES"`
	TelemetryCollectorAddress string `env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	PassSecret                string `env:"PASSWORD_SECRET"`
	EncrytpionKey             string `env:"ENC_KEY"`
}

func getConfigs() configs { // sensible defaults
	cfgs := configs{
		RedisAddr:     "redis:6379",
		RedisPassword: "",
		RedisDB:       0,

		UsersGRPCAddr: "users:50051",
		PostsGRPCAddr: "posts:50051",
		ChatGRPCAddr:  "chat:50051",
		MediaGRPCAddr: "media:50051",

		HTTPAddr:        "0.0.0.0:8081",
		ShutdownTimeout: 5,

		EnableDebugLogs:           true,
		SimplePrint:               true,
		OtelResourceAttributes:    "service.name=api-gateway,service.namespace=social-network,deployment.environment=dev",
		TelemetryCollectorAddress: "alloy:4317",
		PassSecret:                "a2F0LWFsZXgtdmFnLXlwYXQtc3RhbS16b25lMDEtZ28=",
		EncrytpionKey:             "a2F0LWFsZXgtdmFnLXlwYXQtc3RhbS16b25lMDEtZ28=",
	}

	// load environment variables if present
	_, err := configutil.LoadConfigs(&cfgs)
	if err != nil {
		tele.Fatalf("failed to load env variables into config struct: %v", err)
	}

	return cfgs
}

func mytest() {
	seeds := []string{"localhost:9092"}
	// One client can both produce and consume!
	// Consuming can either be direct (no consumer group), or through a group. Below, we use a group.
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup("my-group-identifier"),
		kgo.ConsumeTopics("foo"),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		panic(err)
	}
	defer cl.Close()

	ctx := context.Background()

	// 1.) Producing a message
	// All record production goes through Produce, and the callback can be used
	// to allow for synchronous or asynchronous production.
	// var wg sync.WaitGroup
	// wg.Add(1)
	// fmt.Println("produce async start")
	record := &kgo.Record{Topic: "foo", Value: []byte("bar")}
	// cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
	// 	defer wg.Done()
	// 	if err != nil {
	// 		fmt.Printf("record had a produce error: %v\n", err)
	// 	}

	// })
	// fmt.Println("produce async wait")
	// wg.Wait()
	fmt.Println("produce async fin")
	fmt.Println("produce sync start")
	// Alternatively, ProduceSync exists to synchronously produce a batch of records.
	if err := cl.ProduceSync(ctx, record).FirstErr(); err != nil {
		fmt.Printf("record had a produce error while synchronously producing: %v\n", err)
	}
	fmt.Println("produce sync fin")
	// 2.) Consuming messages from a topic
	for {
		fetches := cl.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			// All errors are retried internally when fetching, but non-retriable errors are
			// returned from polls so that users can notice and take action.
			panic(fmt.Sprint(errs))
		}

		// We can iterate through a record iterator...
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			fmt.Println(string(record.Value), "from an iterator!")
		}

		// or a callback function.
		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			for _, record := range p.Records {
				fmt.Println(string(record.Value), "from range inside a callback!")
			}

			// We can even use a second callback!
			p.EachRecord(func(record *kgo.Record) {
				fmt.Println(string(record.Value), "from a second callback!")
			})
		})
	}
}
