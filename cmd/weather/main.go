package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/celiotk/lab-otel/configs"
	web "github.com/celiotk/lab-otel/internal/infra"
	"github.com/celiotk/lab-otel/internal/infra/provider"
	webserver "github.com/celiotk/lab-otel/internal/infra/web"
	"github.com/celiotk/lab-otel/internal/usecase"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	locationProvider := provider.NewViaCepProvider()
	temperatureProvider := provider.NewWeatherApiProvider(context.Background(), configs.WeatherAPIKey)
	uc := usecase.NewTemperatureFromCepUsecase(temperatureProvider, locationProvider)
	webOrderHandler := web.NewWeatherHandler(*uc)

	ws := webserver.NewWebServer(configs.WebServerPort)
	ws.AddHandler("/temperature/{cep}", webOrderHandler.GetWeather, http.MethodGet)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	ws.Start()

}
