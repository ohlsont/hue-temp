package huetemp

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/amimof/huego"
	"github.com/boltdb/bolt"
)

const (
	hueUser             = "hue-temp"
	authedUserName      = "TMgPM8kRFIvVUn3QDX1ItcHc-WwA9pmvEsxsdIpK"
	openWeatherApiKey   = "64751223d350c8cf76cfcb467091a231"
	tempsensorTypeID    = "ZLLTemperature"
)

type Service struct {
	bridge *huego.Bridge
	db     *bolt.DB
	lastPublished time.Time
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
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		temp, err := s.getCurrentStatus()
		if err != nil {
			return fmt.Errorf("run: %w", err)
		}
		if s.lastPublished.Equal(temp.Time) {
			continue
		}
		s.lastPublished = temp.Time
		log.Printf("current temperature: %.2f : %v", temp.TempCelsius, temp.Time)
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
	temp, lastUpdated, err := s.getIndoorTemp()
	if err != nil {
		return nil, fmt.Errorf("get current status: %w", err)
	}
	outdoorTemp, err := getCurrentTemperatureCelsius(openWeatherApiKey, 57.702757, 11.958240)
	if err != nil {
		return nil, fmt.Errorf("get current status: %w", err)
	}
	return &TemperaturePoint{
		TempCelsius:               temp,
		OutdoorTemperatureCelsius: outdoorTemp,
		Time:                      lastUpdated,
	}, nil
}

func (s *Service) getIndoorTemp() (float64, time.Time, error) {
	sensors, err := s.bridge.GetSensors()
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("get indoor temp: %w", err)
	}

	var temp interface{}
	var timestamp interface{}
	var ok bool
	for i := range sensors {
		if sensors[i].Type != tempsensorTypeID {
			continue
		}
		temp, ok = sensors[i].State["temperature"]
		if !ok {
			continue
		}
		timestamp, ok = sensors[i].State["lastupdated"]
	}
	if !ok {
		return 0, time.Time{}, errors.New("temperature not found")
	}
	temperatureCelsius := (temp.(float64)) / 100
	lastUpdated, err := time.Parse("2006-01-02T15:04:05", timestamp.(string))
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("get indoor temp: %w", err)
	}
	return temperatureCelsius, lastUpdated, nil
}
