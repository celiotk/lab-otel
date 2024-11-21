package entity

import "testing"

func TestTemperatureConversion(t *testing.T) {
	weather := NewWeather("São Paulo", 25)
	weather.ConvertToFahrenheitAndKelvin()

	if weather.TempFahrenheit != 77 {
		t.Errorf("Expected 77, got %v", weather.TempFahrenheit)
	}

	if weather.TempKelvin != 298 {
		t.Errorf("Expected 298, got %v", weather.TempKelvin)
	}
}
