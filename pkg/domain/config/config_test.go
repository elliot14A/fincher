package config_test

import (
	"testing"
	"time"

	"github.com/elliot14A/fincher/pkg/domain/config"
)

func TestConfig_Validation(t *testing.T) {
	validCfg := config.Config{
		Port:                         8080,
		Environment:                  "development",
		StepTimeout:                  30 * time.Second,
		TursoURL:                     "fincher.db",
		TursoToken:                   "mock-token",
		MCPURL:                       "http://127.0.0.1:8000/mcp",
		GeminiAPIKey:                 "mock-key",
		FlashModel:                   "gemini-2.5-flash",
		ProModel:                     "gemini-2.5-pro",
		MaxSyncDriftMs:               120.0,
		VendorDefectThreshold:        0.30,
		ImminentLaunchThresholdHours: 72.0,
		DailyModelInvocationCap:      200,
	}

	if err := validCfg.Validate(); err != nil {
		t.Fatalf("expected valid config, got: %v", err)
	}

	// Invalid Port
	invalid := validCfg
	invalid.Port = 0
	if err := invalid.Validate(); err == nil {
		t.Error("expected error for Port=0")
	}

	// Invalid Environment
	invalid = validCfg
	invalid.Environment = "staging"
	if err := invalid.Validate(); err == nil {
		t.Error("expected error for invalid Environment")
	}

	// Invalid MCP URL
	invalid = validCfg
	invalid.MCPURL = "not-a-valid-url"
	if err := invalid.Validate(); err == nil {
		t.Error("expected error for invalid MCP URL")
	}

	// Invalid Vendor Defect Threshold (> 1.0)
	invalid = validCfg
	invalid.VendorDefectThreshold = 1.5
	if err := invalid.Validate(); err == nil {
		t.Error("expected error for VendorDefectThreshold > 1.0")
	}

	// Invalid Daily Invocation Cap
	invalid = validCfg
	invalid.DailyModelInvocationCap = 0
	if err := invalid.Validate(); err == nil {
		t.Error("expected error for DailyModelInvocationCap <= 0")
	}
}
