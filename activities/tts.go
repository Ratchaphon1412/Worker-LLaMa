package activities

import (
	"context"
	"io"

	"github.com/Ratchaphon1412/worker-llama/configs"
	"github.com/Ratchaphon1412/worker-llama/worker/drivers/http"
	"github.com/Ratchaphon1412/worker-llama/worker/drivers/utils"
)

func TTS(ctx context.Context, conf configs.Config, text string, workflowID string) ([]byte, error) {

	apiClient := http.API{
		BaseURL: conf.TTSAPI,
	}

	jsonPayload := http.JSONBody(map[string]any{
		"input":           text,
		"voice":           conf.TTSModel,
		"response_format": "mp3",
		"speed":           1.1,
	})

	resp, err := apiClient.Post(http.Params{
		Headers: map[string]string{
			"Authorization": "Bearer " + conf.TTSAPIKey,
		},
	}, jsonPayload)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	filePath := conf.TTSSaveToLocal + "/" + workflowID + ".mp3"
	file, err := utils.SaveToFile(filePath, body)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return body, nil
}
