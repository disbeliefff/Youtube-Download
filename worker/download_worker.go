package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"youtubedownload/tasks"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

const (
	DownloadDirectory = "downloads"
)

func HandleYouTubeDownloadTask(ctx context.Context, t *asynq.Task) error {
	var p tasks.DownloadTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// Обновляем прогресс
	rdb.Set(ctx, "progress:"+p.URL, "Downloading...", 0)

	// Скачивание видео
	err := downloadVideo(p.URL, p.Format)
	if err != nil {
		rdb.Set(ctx, "progress:"+p.URL, "Failed", 0)
		return fmt.Errorf("failed to download video: %w", err)
	}

	rdb.Set(ctx, "progress:"+p.URL, "Completed", 0)
	return nil
}

func downloadVideo(url, format string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download video: %w", err)
	}
	defer resp.Body.Close()

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "video/") {
		return fmt.Errorf("URL does not point to a video")
	}

	filename := fmt.Sprintf("%s_%s.mp4", urlToFilename(url), format)
	out, err := os.Create(filepath.Join(DownloadDirectory, filename))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

func urlToFilename(url string) string {
	return filepath.Base(url)
}
