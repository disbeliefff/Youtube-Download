package handlers

import (
	"net/http"
	"youtubedownload/tasks"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

func ShowIndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func DownloadVideo(client *asynq.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.PostForm("url")
		format := c.PostForm("format")

		task, err := tasks.NewDownloadTask(url, format)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": "Failed to create task"})
			return
		}

		_, err = client.Enqueue(task)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": "Failed to enqueue task"})
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{"message": "Download started"})
	}
}
