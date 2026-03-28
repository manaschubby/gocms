package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/manaschubby/gocms/internal/config"
	"github.com/manaschubby/gocms/internal/db"
	"github.com/manaschubby/gocms/internal/modules/cms"
	"github.com/manaschubby/gocms/internal/modules/maintenance"
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
	maint := maintenance.Init(db)

	// Http Server Start Code
	server := echo.New()
	server.Use(middleware.RequestLogger())
	server.Use(middleware.Recover())

	// CMS Routes
	server.GET("/accounts", cms.Handlers.Account.GetAllAccounts)
	server.GET("/content_types", cms.Handlers.ContentType.GetContentType)
	server.POST("/content_types", cms.Handlers.ContentType.CreateContentType)
	server.DELETE("/content_types", cms.Handlers.ContentType.DeleteContentType)
	server.POST("/entries", cms.Handlers.Entry.AddEntry)
	server.GET("/entries", cms.Handlers.Entry.GetEntry)
	server.PUT("/entries", cms.Handlers.Entry.UpdateEntry)

	// Maintenance Routes
	m := server.Group("/maintenance")

	// Categories
	m.GET("/categories", maint.Handlers.Category.GetCategories)
	m.GET("/categories/:id", maint.Handlers.Category.GetCategory)
	m.POST("/categories", maint.Handlers.Category.CreateCategory)
	m.PUT("/categories/:id", maint.Handlers.Category.UpdateCategory)
	m.DELETE("/categories/:id", maint.Handlers.Category.DeleteCategory)

	// Subcategories
	m.GET("/subcategories", maint.Handlers.Subcategory.GetSubcategories) // ?categoryId=
	m.GET("/subcategories/:id", maint.Handlers.Subcategory.GetSubcategory)
	m.POST("/subcategories", maint.Handlers.Subcategory.CreateSubcategory)
	m.PUT("/subcategories/:id", maint.Handlers.Subcategory.UpdateSubcategory)
	m.DELETE("/subcategories/:id", maint.Handlers.Subcategory.DeleteSubcategory)

	// Details
	m.GET("/details", maint.Handlers.Detail.GetDetails) // ?subcategoryId=
	m.GET("/details/:id", maint.Handlers.Detail.GetDetail)
	m.POST("/details", maint.Handlers.Detail.CreateDetail)
	m.PUT("/details/:id", maint.Handlers.Detail.UpdateDetail)
	m.DELETE("/details/:id", maint.Handlers.Detail.DeleteDetail)

	// Workers
	m.GET("/workers", maint.Handlers.Worker.GetWorkers) // ?subcategoryId=
	m.GET("/workers/:id", maint.Handlers.Worker.GetWorker)
	m.POST("/workers", maint.Handlers.Worker.CreateWorker)
	m.PUT("/workers/:id", maint.Handlers.Worker.UpdateWorker)
	m.DELETE("/workers/:id", maint.Handlers.Worker.DeleteWorker)

	// Requests
	m.GET("/requests", maint.Handlers.Request.GetRequests)
	m.GET("/requests/:id", maint.Handlers.Request.GetRequest)
	m.POST("/requests", maint.Handlers.Request.CreateRequest)
	m.POST("/requests/:id/assign", maint.Handlers.Request.AssignWorker)
	m.POST("/requests/:id/resolve", maint.Handlers.Request.Resolve)
	m.POST("/requests/:id/reject", maint.Handlers.Request.Reject)

	// Config
	m.GET("/config", maint.Handlers.Config.GetConfig)
	m.PUT("/config", maint.Handlers.Config.UpdateConfig)

	server.Start(":7467")
}