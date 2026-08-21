package config

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Config holds the application configuration and business rule defaults.
type Config struct {
	// Server settings
	Port        int           `kong:"default=8080,env='FINCHER_PORT',help='HTTP server port'" validate:"required,gte=1,lte=65535"`
	Environment string        `kong:"default='development',env='FINCHER_ENV',help='Runtime environment (development, production)'" validate:"required,oneof=development production test"`
	StepTimeout time.Duration `kong:"default='30s',env='FINCHER_STEP_TIMEOUT',help='Workflow execution timeout'" validate:"required"`

	// Database (Turso / libSQL)
	TursoURL   string `kong:"default='fincher.db',env='FINCHER_TURSO_URL',help='Turso/libSQL database connection URL or local file'" validate:"required"`
	TursoToken string `kong:"env='FINCHER_TURSO_TOKEN',help='Turso authentication token'"`

	// ClickHouse MCP Endpoint
	MCPURL string `kong:"default='http://127.0.0.1:8000/mcp',env='FINCHER_MCP_URL',help='ClickHouse MCP HTTP endpoint'" validate:"required,url"`

	// Google AI Credentials & Models
	GeminiAPIKey string `kong:"env='FINCHER_GEMINI_API_KEY',help='Google Gemini API key'"`
	FlashModel   string `kong:"default='gemini-2.5-flash',env='FINCHER_FLASH_MODEL',help='Model for query agents and notifications'" validate:"required"`
	ProModel     string `kong:"default='gemini-2.5-pro',env='FINCHER_PRO_MODEL',help='Model for assessment and decision nodes'" validate:"required"`

	// Operational Thresholds
	MaxSyncDriftMs               float64 `kong:"default=120.0,env='FINCHER_MAX_SYNC_DRIFT_MS',help='Tolerance threshold for audio sync drift in ms'" validate:"required,gt=0"`
	VendorDefectThreshold        float64 `kong:"default=0.30,env='FINCHER_VENDOR_DEFECT_THRESHOLD',help='Threshold for vendor historical failure rate'" validate:"required,gte=0,lte=1"`
	ImminentLaunchThresholdHours float64 `kong:"default=72.0,env='FINCHER_IMMINENT_LAUNCH_HOURS',help='Threshold for imminent launch in hours'" validate:"required,gt=0"`
	DailyModelInvocationCap      int     `kong:"default=200,env='FINCHER_DAILY_MODEL_CAP',help='Hard daily limit on model calls to protect credit'" validate:"required,gt=0"`
}

// Validate checks the configuration struct using go-playground/validator.
func (c *Config) Validate() error {
	return validate.Struct(c)
}
