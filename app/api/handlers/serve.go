package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/file-share/foundation/database"
	"github.com/FlameInTheDark/file-share/foundation/logs"
	"github.com/FlameInTheDark/file-share/foundation/s3"
)

func Run(logger *zap.Logger) error {
	conf, err := getConfig()
	if err != nil {
		logger.Error("error getting config", zap.Error(err))
		return err
	}

	dbConf := database.Config{
		Host:       conf.Database.Host,
		Port:       conf.Database.Port,
		Database:   conf.Database.Database,
		Username:   conf.Database.Username,
		Password:   conf.Database.Password,
		DisableTLS: conf.Database.DisableTLS,
		Logger:     logs.NewDBLogger(logger),
	}

	db, err := database.NewConnection(dbConf)
	if err != nil {
		logger.Error("error connecting to database", zap.Error(err))
		return err
	}
	defer func() {
		err := db.Close()
		if err != nil {
			logger.Error("error close database connection", zap.Error(err))
		}
	}()

	fs, err := s3.NewMinioClient(
		conf.Minio.Endpoint,
		conf.Minio.AccessTokenID,
		conf.Minio.SecretAccessKey,
		conf.Minio.Region,
		conf.Minio.UseSSL,
		logger,
	)
	if err != nil {
		logger.Error("error creating storage client", zap.Error(err))
		return err
	}
	defer fs.Close()

	handler := NewHandler(db, fs)

	// Start storage notification listener
	go fs.HandleUploadNotification(handler.file.Uploaded)

	httpLogger := logs.NewEchoLogger(logger)

	// Routes
	e := echo.New()
	v1 := e.Group("/api/v1", httpLogger.Middleware)
	v1.GET("/file/:id", handler.Download)
	v1.GET("/file/:id/statistics", handler.Statistics)
	v1.POST("/file", handler.Upload)

	// Health check
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, time.Now().String())
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.Http.Port)))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	return nil
}
