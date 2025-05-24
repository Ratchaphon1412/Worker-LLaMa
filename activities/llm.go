package activities

import (
	"context"
	"io"

	"github.com/Ratchaphon1412/worker-llama/configs"
	"github.com/Ratchaphon1412/worker-llama/worker/drivers/http"
)

type LLMParam struct {
	SystemPrompt string
	Prompt       string
}

func LLM(ctx context.Context, conf configs.Config, prompt LLMParam) ([]byte, error) {

	apiClient := http.API{
		BaseURL: conf.LLMAPI,
	}

	jsonPayload := http.JSONBody(map[string]any{
		"model": conf.LLMModel,
		"messages": []map[string]any{
			{
				"role":    "system",
				"content": prompt.SystemPrompt,
			},
			{
				"role":    "user",
				"content": prompt.Prompt,
			},
		},
		"max_tokens": 150,
		"stream":     false,
	})

	resp, err := apiClient.Post(http.Params{
		Headers: map[string]string{
			"Authorization": "Bearer " + conf.LLMAPIKey,
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

	return body, nil
}
