package config

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Config holds the application configuration and business rule defaults.
type Config struct {
	Port        int           `kong:"default=8080,env='FINCHER_PORT',help='HTTP server port'"`
	Environment string        `kong:"default='development',env='FINCHER_ENV',help='Runtime environment (development, production)'" validate:"oneof=development production staging"`
	StepTimeout time.Duration `kong:"default='30s',env='FINCHER_STEP_TIMEOUT',help='Workflow execution timeout'"`

	TursoURL   string `kong:"default='fincher.db',env='FINCHER_TURSO_URL',help='Turso/libSQL database connection URL or local file'"`
	TursoToken string `kong:"env='FINCHER_TURSO_TOKEN',help='Turso authentication token'"`

	MCPURL string `kong:"default='http://127.0.0.1:8000/mcp',env='FINCHER_MCP_URL',help='ClickHouse MCP HTTP endpoint'"`

	GeminiAPIKey string `kong:"env='FINCHER_GEMINI_API_KEY',help='Google Gemini API key'"`
	FlashModel   string `kong:"default='gemini-2.5-flash',env='FINCHER_FLASH_MODEL',help='Model for query agents and notifications'"`
	ProModel     string `kong:"default='gemini-2.5-pro',env='FINCHER_PRO_MODEL',help='Model for assessment and decision nodes'"`

	MaxSyncDriftMs               float64 `kong:"default='120.0',env='FINCHER_MAX_SYNC_DRIFT_MS',help='Tolerance threshold for audio sync drift in ms'" validate:"gte=0"`
	VendorDefectThreshold        float64 `kong:"default='0.30',env='FINCHER_VENDOR_DEFECT_THRESHOLD',help='Threshold for vendor historical failure rate'" validate:"gte=0,lte=1"`
	ImminentLaunchThresholdHours float64 `kong:"default='72.0',env='FINCHER_IMMINENT_LAUNCH_HOURS',help='Threshold for imminent launch in hours'" validate:"gte=0"`
	DailyModelInvocationCap      int     `kong:"default=200,env='FINCHER_DAILY_MODEL_CAP',help='Hard daily limit on model calls to protect credit'" validate:"gte=1"`
}

// Validate verifies configuration constraints.
func (c *Config) Validate() error {
	return validate.Struct(c)
}
