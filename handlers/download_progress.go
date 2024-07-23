package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func GetDownloadProgress(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.String(http.StatusBadRequest, "URL is required")
		return
	}

	ctx := c.Request.Context()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	progress, err := rdb.Get(ctx, "progress:"+url).Result()
	if err == redis.Nil {
		c.String(http.StatusNotFound, "Progress not found")
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, "Error retrieving progress")
		return
	}

	c.String(http.StatusOK, progress)
}
