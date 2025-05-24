package activities

import (
	"context"
	"encoding/json"

	"github.com/Ratchaphon1412/worker-llama/configs"
	"github.com/Ratchaphon1412/worker-llama/worker/drivers/database"
)

func PublisherToChat(ctx context.Context, conf configs.Config, channel string, answer Answer) error {

	database.ConnectRedis(&conf)
	// Publish the media url to the Redis channel

	answerJson, err := json.Marshal(answer)
	if err != nil {
		return err
	}

	if err := database.Redis.Rd.Publish(ctx, channel, answerJson).Err(); err != nil {
		return err
	}
	defer database.Redis.Rd.Close()

	return nil
}
