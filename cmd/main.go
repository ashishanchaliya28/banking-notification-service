package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/banking-superapp/notification-service/config"
	"github.com/banking-superapp/notification-service/handler"
	"github.com/banking-superapp/notification-service/repository"
	"github.com/banking-superapp/notification-service/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	cfg := config.Load()
	mongoClient, err := repository.NewMongoClient(cfg.MongoAtlasURI)
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	db := mongoClient.Database("banking_notifications")
	if err := repository.CreateIndexes(db); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	notifRepo := repository.NewNotificationRepo(db)
	deviceRepo := repository.NewDeviceTokenRepo(db)
	prefRepo := repository.NewPreferencesRepo(db)
	notifSvc := service.NewNotificationService(notifRepo, deviceRepo, prefRepo)
	notifHandler := handler.NewNotificationHandler(notifSvc)

	app := fiber.New(fiber.Config{AppName: cfg.ServiceName, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second})
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": cfg.ServiceName})
	})

	v1 := app.Group("/v1")
	notifs := v1.Group("/notifications")
	notifs.Post("/send", notifHandler.Send)
	notifs.Get("/", notifHandler.GetNotifications)
	notifs.Put("/preferences", notifHandler.UpdatePreferences)
	notifs.Post("/device/register", notifHandler.RegisterDevice)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Printf("Starting %s on port %s", cfg.ServiceName, cfg.Port)
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = app.ShutdownWithContext(ctx)
}
