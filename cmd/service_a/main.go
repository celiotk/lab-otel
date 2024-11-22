package main

import (
	"fmt"
	"net/http"

	web "github.com/celiotk/lab-otel/internal/infra"
	"github.com/celiotk/lab-otel/internal/infra/provider"
	webserver "github.com/celiotk/lab-otel/internal/infra/web"
	"github.com/celiotk/lab-otel/internal/usecase"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	webServerPort := viper.GetString("WEB_SERVER_PORT")

	tempFromBProvider := provider.NewServiceBProvider(viper.GetString("SERVICE_B_ADDRESS"))
	uc := usecase.NewTemperatureFromServiceBUsecase(tempFromBProvider)
	serviceAHandler := web.NewServiceAHandler(*uc)

	ws := webserver.NewWebServer(webServerPort)
	ws.AddHandler("/temperature", serviceAHandler.GetWeather, http.MethodPost)
	fmt.Println("Starting web server on port", webServerPort)
	ws.Start()

}
