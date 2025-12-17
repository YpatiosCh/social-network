package entry

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"social-network/services/media/internal/application"
	"social-network/services/media/internal/client"
	"social-network/services/media/internal/configs"
	"social-network/services/media/internal/db/dbservice"
	"social-network/services/media/internal/handler"
	"social-network/services/media/internal/validator"
	pb "social-network/shared/gen-go/media"
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/gorpc"

	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func Run() error {
	cfgs := configs.Config{
		Server: configs.Server{
			Port: os.Getenv("SERVICE_PORT"),
		},
		DB: configs.Db{
			URL:                os.Getenv("DATABASE_URL"),
			StaleFilesInterval: 1 * time.Hour,
		},
		FileService: configs.FileService{
			Buckets: configs.Buckets{
				Originals: "uploads-originals",
				Variants:  "uploads-variants",
			},
			VariantWorkerInterval: 30 * time.Second,
			FileConstraints: configs.FileConstraints{
				MaxImageUpload: 5 << 20, // 5MB
				MaxWidth:       4096,
				MaxHeight:      4096,
				AllowedMIMEs: map[string]bool{
					"image/jpeg": true,
					"image/png":  true,
					"image/gif":  true,
					"image/webp": true,
				},
			},
			Endpoint:  os.Getenv("MINIO_ENDPOINT"),
			AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
			Secret:    os.Getenv("MINIO_SECRET_KEY"),
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := connectToDb(ctx, cfgs.DB.URL)
	if err != nil {
		return fmt.Errorf("failed to connect db: %v", err)
	}
	defer pool.Close()

	log.Println("Connected to media database")

	fileServiceClient, err := NewMinIOConn(cfgs.FileService)
	if err != nil {
		return err
	}
	querier := dbservice.NewQuerier(pool)
	app := application.NewMediaService(
		pool,
		&client.Clients{
			MinIOClient: fileServiceClient,
			Validator: &validator.ImageValidator{
				Config: cfgs.FileService.FileConstraints,
			},
		},
		querier,
		cfgs,
	)
	w := dbservice.NewWorker(querier)

	app.StartVariantWorker(ctx, cfgs.FileService.VariantWorkerInterval)
	w.StartStaleFilesWorker(ctx, cfgs.DB.StaleFilesInterval)

	service := &handler.MediaHandler{
		Application: app,
		Configs:     cfgs.Server,
	}

	log.Println("Running gRpc service...")

	grpc, err := RunGRPCServer(service)
	if err != nil {
		return err
	}

	// wait here for process termination signal to initiate graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")
	cancel()
	grpc.GracefulStop()
	log.Println("Server stopped")
	return nil
}

func connectToDb(ctx context.Context, connStr string) (pool *pgxpool.Pool, err error) {
	for i := range 10 {
		pool, err = pgxpool.New(ctx, connStr)
		if err == nil {
			break
		}
		log.Printf("DB not ready yet (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	return pool, err
}

// RunGRPCServer starts the gRPC server and blocks
func RunGRPCServer(s *handler.MediaHandler) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", s.Configs.Port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", s.Configs.Port, err)
	}

	customUnaryInterceptor, err := gorpc.UnaryServerInterceptorWithContextKeys([]gorpc.StringableKey{ct.UserId, ct.ReqID, ct.TraceId}...)
	if err != nil {
		return nil, err
	}
	customStreamInterceptor, err := gorpc.StreamServerInterceptorWithContextKeys([]gorpc.StringableKey{ct.UserId, ct.ReqID, ct.TraceId}...)
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(customUnaryInterceptor),
		grpc.StreamInterceptor(customStreamInterceptor),
	)

	pb.RegisterMediaServiceServer(grpcServer, s)

	log.Printf("gRPC server listening on %s", s.Configs.Port)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()
	return grpcServer, nil
}
