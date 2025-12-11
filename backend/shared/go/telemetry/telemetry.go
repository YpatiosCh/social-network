package tele

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel/exporters/prometheus"
)

type telemetry struct {
}

func initOpenTelemetrySDK(ctx context.Context) func() {
	otelShutdown, err := SetupOTelSDK(ctx)
	if err != nil {
		log.Fatal("open telemetry sdk failed, ERROR:", err.Error())
	}
	fmt.Println("open telemetry ready")

	return func() {
		err := otelShutdown(context.Background())
		if err != nil {
			log.Println("otel shutdown ungracefully! ERROR: " + err.Error())
		} else {
			log.Println("otel shutdown gracefully")
		}
	}
}

func NewPrometheus() *prometheus.Exporter {
	x, _ := prometheus.New(nil)
	return x
}
