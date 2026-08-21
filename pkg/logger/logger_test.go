package logger_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/elliot14A/fincher/pkg/logger"
)

func TestLogger_Levels(t *testing.T) {
	var buf bytes.Buffer
	logger.Init("development", &buf)

	logger.Debug("test debug", "key", "val")
	logger.Info("test info", "count", 42)
	logger.Warn("test warn", "issue", "minor")
	logger.Error("test error", "err", "fail")

	output := buf.String()
	if !strings.Contains(output, "test debug") {
		t.Errorf("expected output to contain debug log, got: %s", output)
	}
	if !strings.Contains(output, "test info") {
		t.Errorf("expected output to contain info log, got: %s", output)
	}
	if !strings.Contains(output, "test warn") {
		t.Errorf("expected output to contain warn log, got: %s", output)
	}
	if !strings.Contains(output, "test error") {
		t.Errorf("expected output to contain error log, got: %s", output)
	}

	buf.Reset()
	logger.DebugContext(context.Background(), "context debug")
	if !strings.Contains(buf.String(), "context debug") {
		t.Errorf("expected output to contain context debug, got: %s", buf.String())
	}
}

func TestLogger_ProductionJSON(t *testing.T) {
	var buf bytes.Buffer
	logger.Init("production", &buf)

	logger.Info("json test", "service", "fincher")
	output := buf.String()
	if !strings.Contains(output, `"level":"INFO"`) || !strings.Contains(output, `"msg":"json test"`) {
		t.Errorf("expected json log output, got: %s", output)
	}
}
