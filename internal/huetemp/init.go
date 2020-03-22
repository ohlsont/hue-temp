package huetemp

import (
	"errors"
	"fmt"
	"time"

	"github.com/amimof/huego"
	"github.com/boltdb/bolt"
)

const (
	hueUser             = "hue-temp"
	authedUserName      = "TMgPM8kRFIvVUn3QDX1ItcHc-WwA9pmvEsxsdIpK"
	temperatureSensorID = 6
	openWeatherApiKey   = "64751223d350c8cf76cfcb467091a231"
)

type Service struct {
	bridge *huego.Bridge
	db     *bolt.DB
}

type TemperaturePoint struct {
	TempCelsius               float64   `json:"temp_celsius"`
	OutdoorTemperatureCelsius float64   `json:"outdoor_temperature_celsius"`
	Time                      time.Time `json:"time"`
}

func Init(db *bolt.DB) (*Service, error) {
	bridge, err := huego.Discover()
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	bridge.Login(authedUserName)
	return &Service{
		bridge: bridge,
		db:     db,
	}, nil
}

func (s *Service) Run() error {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		temp, err := s.getCurrentStatus()
		if err != nil {
			return fmt.Errorf("run: %w", err)
		}
		err = s.InsertDataPoint(temp)
		if err != nil {
			return fmt.Errorf("run: %w", err)
		}
		pts, err := s.getDataPoints()
		if err != nil {
			return fmt.Errorf("run: %w", err)
		}
		err = Draw(pts)
		if err != nil {
			return fmt.Errorf("run: %w", err)
		}
	}
	return nil
}

func (s *Service) getCurrentStatus() (*TemperaturePoint, error) {
	temp, err := s.getIndoorTemp()
	if err != nil {
		return nil, fmt.Errorf("run: %w", err)
	}
	outdoorTemp, err := getCurrentTemperatureCelsius(openWeatherApiKey, 57.702757, 11.958240)
	if err != nil {
		return nil, fmt.Errorf("run: %w", err)
	}
	return &TemperaturePoint{
		TempCelsius:               temp,
		OutdoorTemperatureCelsius: outdoorTemp,
		Time:                      time.Now(),
	}, nil
}

func (s *Service) getIndoorTemp() (float64, error) {
	tempSensor, err := s.bridge.GetSensor(temperatureSensorID)
	if err != nil {
		return 0, fmt.Errorf("connect: %w", err)
	}
	temp, ok := tempSensor.State["temperature"]
	if !ok {
		return 0, errors.New("temperature not found")
	}
	temperatureCelsius := (temp.(float64)) / 100
	return temperatureCelsius, nil
}
