package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/hibiken/asynq"
)

const (
	TypeYouTubeDownload = "youtube:download"
)

type YoutubeDownloadPayload struct {
	URL    string `json:"url"`
	Format string `json:"format"`
}

func NewDownloadTask(url, format string) (*asynq.Task, error) {
	payload, err := json.Marshal(YoutubeDownloadPayload{URL: url, Format: format})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeYouTubeDownload, payload), nil
}

func HandleYoutubeDownload(ctx context.Context, t *asynq.Task) error {
	var p YoutubeDownloadPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}
	output := "video." + p.Format
	cmd := exec.Command("yt-dlp", "-f", p.Format, "-o", output, p.URL)
	err := cmd.Run()
	if err != nil {
		fmt.Errorf("yt-dlp failed: %v", err)
	}
	return nil
}
