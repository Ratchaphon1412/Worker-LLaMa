package activities

import (
	"context"

	"github.com/Ratchaphon1412/worker-llama/configs"
	"github.com/Ratchaphon1412/worker-llama/worker/drivers/http"
)

func WebScrap(ctx context.Context, conf configs.Config, urlItem ResearchResult) ([]string, error) {
	var scrapResult []string
	for _, item := range urlItem.Items {
		if item.Kind == "customsearch#result" {
			text, err := http.GetTextFromURL(item.Link)
			if err == nil {
				cleanText := http.NormalizeSpace(text)
				summary := http.SummarizeText(cleanText, 5) // keep only 3 leading sentences
				scrapResult = append(scrapResult, summary)
			}

		}
	}
	return scrapResult, nil
}
