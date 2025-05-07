package configs

type Config struct {
	TemporalHostPort  string `env:"TEMPORAL_HOST_PORT"`
	TemporalNamespace string `env:"TEMPORAL_NAMESPACE"`
	TemporalTaskQueue string `env:"TEMPORAL_TASK_QUEUE"`

	LLMAPI    string `env:"LLM_API"`
	LLMAPIKey string `env:"LLM_API_KEY"`
	LLMModel  string `env:"LLM_MODEL"`

	TTSAPI         string `env:"TTS_API"`
	TTSAPIKey      string `env:"TTS_API_KEY"`
	TTSModel       string `env:"TTS_MODEL"`
	TTSSaveToLocal string `env:"TTS_SAVE_TO_LOCAL"`

	GoogleAPIKEYCustomSearch   string `env:"GOOGLE_API_KEY_CUSTOM_SEARCH"`
	GoogleCustomSearchEngineID string `env:"GOOGLE_CUSTOM_SEARCH_ENGINE_ID"`
	GoogleCustomSearchURL      string `env:"GOOGLE_CUSTOM_SEARCH_URL"`
	GoogleMaxResults           int    `env:"GOOGLE_MAX_RESULTS"`
}
