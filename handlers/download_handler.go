package handlers

import (
	"net/http"
	"strings"
	"youtubedownload/tasks"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

func ShowIndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// DownloadHandler обрабатывает запрос на скачивание видео
func DownloadHandler(c *gin.Context) {
	url := c.PostForm("url")
	format := c.PostForm("format")

	if url == "" || format == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL and format are required"})
		return
	}

	task, err := tasks.NewDownloadTask(url, format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create download task"})
		return
	}

	// Используем клиент Asynq для постановки задачи в очередь
	client := c.MustGet("asynqClient").(*asynq.Client)
	if _, err := client.Enqueue(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Download started", "url": url})
}

// ProgressHandler отслеживает прогресс скачивания
func ProgressHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL is required"})
		return
	}

	progress, err := getDownloadProgress(c, url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get progress"})
		return
	}

	c.String(http.StatusOK, progress)
}

// getDownloadProgress получает прогресс скачивания из Redis
func getDownloadProgress(c *gin.Context, url string) (string, error) {
	rdb := c.MustGet("redisClient").(*redis.Client)

	progress, err := rdb.Get(c.Request.Context(), "progress:"+url).Result()
	if err != nil {
		if strings.Contains(err.Error(), "redis: nil") {
			return "Not started", nil
		}
		return "", err
	}
	return progress, nil
}
