package config_test

import (
	"testing"
	"time"

	"github.com/elliot14A/fincher/internal/config"
)

func TestConfig_Validation(t *testing.T) {
	validCfg := config.Config{
		Port:                         8080,
		Environment:                  "development",
		StepTimeout:                  30 * time.Second,
		TursoURL:                     "fincher.db",
		TursoToken:                   "test-token",
		MCPURL:                       "http://127.0.0.1:8000/mcp",
		GeminiAPIKey:                 "test-key",
		FlashModel:                   "gemini-2.5-flash",
		ProModel:                     "gemini-2.5-pro",
		MaxSyncDriftMs:               120.0,
		VendorDefectThreshold:        0.30,
		ImminentLaunchThresholdHours: 72.0,
		DailyModelInvocationCap:      200,
	}

	if err := validCfg.Validate(); err != nil {
		t.Fatalf("expected valid config to pass validation, got: %v", err)
	}

	invalidPort := validCfg
	invalidPort.Environment = "invalid-env"
	if err := invalidPort.Validate(); err == nil {
		t.Errorf("expected invalid environment to fail validation")
	}

	invalidVendor := validCfg
	invalidVendor.VendorDefectThreshold = 1.5
	if err := invalidVendor.Validate(); err == nil {
		t.Errorf("expected invalid vendor defect threshold to fail validation")
	}

	invalidCap := validCfg
	invalidCap.DailyModelInvocationCap = 0
	if err := invalidCap.Validate(); err == nil {
		t.Errorf("expected invalid daily cap to fail validation")
	}
}
