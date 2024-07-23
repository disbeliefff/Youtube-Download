package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"youtubedownload/handlers"
	"youtubedownload/tasks"
	"youtubedownload/worker"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	rdb    *redis.Client
	logger *zap.Logger
)

func main() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("Cannot start logger, error is %s", err))
	}
	defer func() {
		_ = logger.Sync()
	}()

	redisURL := os.Getenv("REDIS_ADDR")
	if redisURL == "" {
		logger.Fatal("REDIS_ADDR is not set")
	}

	var redisAddr, redisPassword string
	if strings.HasPrefix(redisURL, "redis://") {
		parsedURL, err := url.Parse(redisURL)
		if err != nil {
			logger.Fatal("Invalid REDIS_ADDR", zap.Error(err))
		}
		redisAddr = parsedURL.Host
		if parsedURL.User != nil {
			redisPassword, _ = parsedURL.User.Password()
		}
	} else {
		redisAddr = redisURL
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPassword,
	})
	defer asynqClient.Close()

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		c.Set("asynqClient", asynqClient)
		c.Next()
		logger.Info("Request",
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.String()),
			zap.Int("status", c.Writer.Status()),
			zap.String("client_ip", c.ClientIP()))
	})

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	router.GET("/", handlers.ShowIndexPage)
	router.POST("/download", handlers.DownloadHandler)
	router.GET("/progress", handlers.GetDownloadProgress)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: redisPassword,
		},
		asynq.Config{
			Concurrency: 10,
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeYouTubeDownload, worker.HandleYouTubeDownloadTask)

	go func() {
		if err := srv.Run(mux); err != nil {
			logger.Fatal("failed to run Asynq server", zap.Error(err))
		}
	}()

	go func() {
		if err := router.Run(":" + port); err != nil {
			logger.Fatal("failed to run server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	srv.Shutdown()
}
