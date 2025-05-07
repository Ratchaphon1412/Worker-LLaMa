package main

import (
	"log"

	"github.com/Ratchaphon1412/worker-llama/activities"
	"github.com/Ratchaphon1412/worker-llama/configs"
	"github.com/Ratchaphon1412/worker-llama/workflow"
	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var cfg configs.Config

func main() {

	// Load environment variables
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse env vars: %v", err)
	}

	// Create a worker that listens on the llm task queue
	c, err := client.Dial(client.Options{
		HostPort:  cfg.TemporalHostPort,
		Namespace: cfg.TemporalNamespace,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	w := worker.New(c, cfg.TemporalTaskQueue, worker.Options{})

	w.RegisterWorkflow(workflow.AIWorkflow)
	w.RegisterActivity(activities.Research)
	w.RegisterActivity(activities.WebScrap)
	w.RegisterActivity(activities.LLM)
	w.RegisterActivity(activities.TTS)
	w.RegisterActivity(activities.Storage)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
