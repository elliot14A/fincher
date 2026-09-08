package config

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

const (
	DefaultTimeScale             = time.Second
	DefaultMaxSyncDriftMs        = 120.0
	DefaultVendorDefectThreshold = 0.05
	DefaultPremiereUrgentHours   = 72.0
	DefaultTurnaroundHours       = 12.0
	DefaultMasterReconformHours  = 48.0
	MaxRedeliveryAttempts        = 3
)

type Config struct {
	Port        int           `kong:"default=8080,env='FINCHER_PORT',help='HTTP server port'"`
	Environment string        `kong:"default='development',env='FINCHER_ENV',help='Runtime environment (development, production)'" validate:"oneof=development production staging"`
	StepTimeout time.Duration `kong:"default='30s',env='FINCHER_STEP_TIMEOUT',help='Workflow execution timeout'"`

	TursoURL   string `kong:"default='fincher.db',env='FINCHER_TURSO_URL',help='Turso/libSQL database connection URL or local file'"`
	TursoToken string `kong:"env='FINCHER_TURSO_TOKEN',help='Turso authentication token'"`

	ClickHouseDSN string `kong:"default='127.0.0.1:9000',env='FINCHER_CLICKHOUSE_DSN',help='ClickHouse native TCP endpoint'"`
	MCPURL        string `kong:"default='http://127.0.0.1:8000/mcp',env='FINCHER_MCP_URL',help='ClickHouse MCP HTTP endpoint'"`

	GeminiAPIKey  string            `kong:"env='FINCHER_GEMINI_API_KEY',help='Google Gemini API key'"`
	GeminiModel   string            `kong:"default='gemini-2.5-flash',env='FINCHER_GEMINI_MODEL',help='Gemini model for all agent workflows'"`
	GeminiOptions map[string]string `kong:"env='FINCHER_GEMINI_OPTIONS',help='Gemini client options as key=value pairs (e.g. location=asia-south1)'"`

	DailyModelInvocationCap     int `kong:"default=200,env='FINCHER_DAILY_MODEL_CAP',help='Hard daily limit on model calls to protect credit'" validate:"gte=1"`
	MaxConcurrentAIWorkflowRuns int `kong:"default=3,env='FINCHER_MAX_CONCURRENT_AI_WORKFLOW_RUNS',help='Concurrency semaphore for active AI workflow runs'" validate:"gte=1"`
}

func (c *Config) Validate() error {
	return validate.Struct(c)
}
