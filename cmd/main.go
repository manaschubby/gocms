package main

import (
	"log"

	"github.com/manaschubby/gocms/internal/config"
	"github.com/manaschubby/gocms/internal/db"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Panicf("ENV File failed to load. Err: %v Quitting...", err)
	}

	db, err := db.Connect(*cfg)
	log.Println(db, err)
}
