package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	web "github.com/celiotk/lab-otel/internal/infra"
	"github.com/celiotk/lab-otel/internal/infra/provider"
	webserver "github.com/celiotk/lab-otel/internal/infra/web"
	"github.com/celiotk/lab-otel/internal/otel_provider"
	"github.com/celiotk/lab-otel/internal/usecase"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	viper.AutomaticEnv()
	webServerPort := viper.GetString("WEB_SERVER_PORT")

	shutdown, err := otel_provider.InitProvider("service_a", viper.GetString("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	otelTracer := otel.Tracer("microservice-tracer")

	tempFromBProvider := provider.NewServiceBProvider(viper.GetString("SERVICE_B_ADDRESS"), otelTracer)
	uc := usecase.NewTemperatureFromServiceBUsecase(tempFromBProvider, otelTracer)
	serviceAHandler := web.NewServiceAHandler(*uc, otelTracer)

	ws := webserver.NewWebServer(webServerPort)
	ws.AddHandler("/temperature", serviceAHandler.PostTemperature, http.MethodPost)
	fmt.Println("Starting web server on port", webServerPort)
	go func() {
		if err := ws.Start(); err != nil {
			panic(err)
		}
	}()

	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down due to other reason...")
	}

	// Create a timeout context for the graceful shutdown
	ctx2, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := ws.Stop(ctx2); err != nil {
		log.Fatal("Failed to stop server: %w", err)
	}
}
