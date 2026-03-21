package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/manaschubby/gocms/internal/config"
	"github.com/manaschubby/gocms/internal/db"
	"github.com/manaschubby/gocms/internal/modules/cms"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Panicf("ENV File failed to load. Err: %v Quitting...", err)
	}

	db, err := db.Connect(*cfg)
	if err != nil {
		log.Fatalf("DB failed to load. Err: %v Quitting...", err)
	}
	log.Println("Successfully connected to DB: " + db.DriverName())

	cms := cms.Init(cfg, db)

	// Http Server Start Code
	server := echo.New()
	server.Use(middleware.RequestLogger())
	server.Use(middleware.Recover())

	// Route Registration
	server.GET("/accounts", cms.Handlers.Account.GetAllAccounts)
	server.GET("/content_types", cms.Handlers.ContentType.GetContentType)
	server.POST("/content_types", cms.Handlers.ContentType.CreateContentType)
	server.DELETE("/content_types", cms.Handlers.ContentType.DeleteContentType)

	server.Start(":7467")
}
