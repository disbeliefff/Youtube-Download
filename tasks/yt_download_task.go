package tasks

import (
	"encoding/json"
	"github.com/hibiken/asynq"
)

const (
	TypeYouTubeDownload = "youtube:download"
)

type DownloadTaskPayload struct {
	URL    string `json:"url"`
	Format string `json:"format"`
}

func NewDownloadTask(url, format string) (*asynq.Task, error) {
	payload := DownloadTaskPayload{
		URL:    url,
		Format: format,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeYouTubeDownload, data), nil
}
