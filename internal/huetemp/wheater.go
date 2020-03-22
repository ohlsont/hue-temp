package huetemp

import (
	"fmt"
	owm "github.com/briandowns/openweathermap"
)

func getCurrentTemperatureCelsius(openWeatherApiKey string, lat, lon float64) (float64, error)  {
	w, err := owm.NewCurrent("C", "se", openWeatherApiKey)
	if err != nil {
		return 0,fmt.Errorf("get weather: %w", err)
	}
	err = w.CurrentByCoordinates(&owm.Coordinates{
		Latitude:  lat,
		Longitude: lon,
	})
	if err != nil {
		return 0, fmt.Errorf("get weather: %w", err)
	}
	return w.Main.Temp, nil
}
