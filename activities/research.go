package activities

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/Ratchaphon1412/worker-llama/configs"
	"github.com/Ratchaphon1412/worker-llama/worker/drivers/http"
)

type ResearchResult struct {
	Items []struct {
		Kind             string         `json:"kind"`
		Title            string         `json:"title"`
		HtmlTitle        string         `json:"htmlTitle"`
		Link             string         `json:"link"`
		DisplayLink      string         `json:"displayLink"`
		Snippet          string         `json:"snippet"`
		HtmlSnippet      string         `json:"htmlSnippet"`
		FormattedUrl     string         `json:"formattedUrl"`
		HtmlFormattedUrl string         `json:"htmlFormattedUrl"`
		Pagemap          map[string]any `json:"pagemap"`
		CacheId          string         `json:"cacheId"`
	} `json:"items"`
}

func Research(ctx context.Context, conf configs.Config, topic string) (ResearchResult, error) {

	apiClient := http.API{
		BaseURL: conf.GoogleCustomSearchURL,
	}

	resp, err := apiClient.Get(http.Params{
		Query: map[string]string{
			"q":   topic,
			"key": conf.GoogleAPIKEYCustomSearch,
			"cx":  conf.GoogleCustomSearchEngineID,
			"num": strconv.Itoa(conf.GoogleMaxResults),
		},
	})
	if err != nil {
		return ResearchResult{}, err
	}
	var result ResearchResult

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ResearchResult{}, err
	}

	return result, nil

}
