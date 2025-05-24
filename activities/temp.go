package activities

import (
	"context"

	"github.com/Ratchaphon1412/worker-llama/configs"
	"github.com/Ratchaphon1412/worker-llama/worker/drivers/utils"
)

func ClearTemp(ctx context.Context, conf configs.Config, filename string) error {
	err := utils.DeleteFile(conf.TTSSaveToLocal + filename)
	if err != nil {
		return err
	}

	return nil
}
