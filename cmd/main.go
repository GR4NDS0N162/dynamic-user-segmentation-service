package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/GR4NDS0N162/dynamic-user-segmentation-service/pkg/handler"
	"github.com/GR4NDS0N162/dynamic-user-segmentation-service/pkg/repository"
	"github.com/GR4NDS0N162/dynamic-user-segmentation-service/pkg/service"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitializeDatabase() *gorm.DB {
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("Failed to initialize database. %s", err.Error())
	}
	return db
}

func CloseDatabaseConnection(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Errorf("Error occured on retrieving sql database: %s", err.Error())
	}

	if err := sqlDB.Close(); err != nil {
		logrus.Errorf("Error occured on database connection close: %s", err.Error())
	}
}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	db := InitializeDatabase()
	app := handler.NewHandler(service.NewService(repository.NewRepository(db))).InitRoutes()

	go func() {
		addr := fmt.Sprintf(":%v", os.Getenv("WEB_PORT"))
		if err := app.Listen(addr); err != nil {
			logrus.Fatalf("Error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("Dynamic user segmentation service started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Dynamic user segmentation service shutting down")

	if err := app.ShutdownWithContext(context.Background()); err != nil {
		logrus.Errorf("Error occurred on server shutting down: %s", err.Error())
	}

	CloseDatabaseConnection(db)
}
