package main

import (
	"github.com/boltdb/bolt"
	"github.com/ohlsont/hue-temp/internal/huetemp"
	"log"
)

func main()  {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	s, err := huetemp.Init(db)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}

}
