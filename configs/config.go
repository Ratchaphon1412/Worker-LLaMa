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

	MinioEndpoint      string `env:"MINIO_ENDPOINT"`
	MinioUserAccessKey string `env:"MINIO_USER_ACCESS_KEY"`
	MinioUserSecretKey string `env:"MINIO_USER_SECRET_KEY"`
	MinioDefaultBucket string `env:"MINIO_DEFAULT_BUCKETS"`
	MinioSSLEnabled    bool   `env:"MINIO_SSL_ENABLED"`
	MinioPublicURL     string `env:"MINIO_PUBLIC_URL"`

	DB_HOST     string `env:"DB_HOST"`
	DB_PORT     string `env:"DB_PORT"`
	DB_USER     string `env:"DB_USER"`
	DB_PASSWORD string `env:"DB_PASSWORD"`
	DB_NAME     string `env:"DB_NAME"`
	DB_TIMEZONE string `env:"DB_TIMEZONE"`
	DB_SSL_MODE string `env:"DB_SSL_MODE"`

	// Redis Configuration
	REDIS_ADDR string `env:"REDIS_ADDR"`

	REDIS_USERNAME  string `env:"REDIS_USERNAME"`
	REDIS_PASSWORD  string `env:"REDIS_PASSWORD"`
	REDIS_DATABASE  int    `env:"REDIS_DATABASE"`
	REDIS_POOLFIFO  bool   `env:"REDIS_POOLFIFO"`
	REDIS_POOL_SIZE int    `env:"REDIS_POOL_SIZE"`
}
