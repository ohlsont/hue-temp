package huetemp

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"strconv"
	"time"
)

const (
	bucketNameOld = "hue-temp"
	bucketName = "hue-temp-v2"
)
func (s *Service) getDataPoints() ([]*TemperaturePoint, error) {
	dataPoints := []*TemperaturePoint{}
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.ForEach(func(k, v []byte) error {
			te := &TemperaturePoint{}
			err := json.Unmarshal(v, te)
			if err != nil {
				return fmt.Errorf("view bucket: %w", err)
			}
			dataPoints = append(dataPoints, te)
			return nil
		})
		if err != nil {
			return fmt.Errorf("view bucket: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("get data points: %w", err)
	}
	return dataPoints, nil
}

func (s *Service) InsertDataPoint(dataPoint *TemperaturePoint) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("update db: %w", err)
		}
		data, err := json.Marshal(dataPoint)
		if err != nil {
			return fmt.Errorf("could not insert weight: %v", err)
		}
		err = b.Put([]byte(dataPoint.Time.Format(time.RFC3339)), data)
		if err != nil {
			return fmt.Errorf("could not insert weight: %v", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("insert data poitn")
	}
	return nil
}

func (s *Service) migrateData() error {
	dataPoints := []*TemperaturePoint{}
	log.Println("get datapoints")
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketNameOld))
		err := b.ForEach(func(k, v []byte) error {
			t, err := time.Parse(time.RFC3339, string(k))
			if err != nil {
				return fmt.Errorf("view bucket: %w", err)
			}
			temp, err := strconv.ParseFloat(string(v),10)
			if err != nil {
				return fmt.Errorf("view bucket: %w", err)
			}
			dataPoints = append(dataPoints, &TemperaturePoint{
				TempCelsius:               temp,
				OutdoorTemperatureCelsius: 0,
				Time:                      t,
			})
			return nil
		})
		if err != nil {
			return fmt.Errorf("view bucket: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("view bucket: %w", err)
	}
	log.Println("loop datapoints")
	for i, d := range dataPoints {
		log.Println("loop datapoints", i)
		err := s.InsertDataPoint(d)
		if err != nil {
			return fmt.Errorf("view bucket: %w", err)
		}
	}
	return nil
}
