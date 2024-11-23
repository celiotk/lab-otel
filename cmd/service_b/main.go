package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/celiotk/lab-otel/configs"
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

	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	shutdown, err := otel_provider.InitProvider("service_b", viper.GetString("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	otelTracer := otel.Tracer("microservice-tracer")

	locationProvider := provider.NewViaCepProvider(otelTracer)
	temperatureProvider := provider.NewWeatherApiProvider(configs.WeatherAPIKey, otelTracer)
	uc := usecase.NewTemperatureByCepUsecase(temperatureProvider, locationProvider, otelTracer)
	weatherHandler := web.NewTemperatureByCepHandler(*uc, otelTracer)

	ws := webserver.NewWebServer(configs.WebServerPort)
	ws.AddHandler("/temperature/{cep}", weatherHandler.GetTemperature, http.MethodGet)
	fmt.Println("Starting web server on port", configs.WebServerPort)
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
