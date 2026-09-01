package titles_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"google.golang.org/adk/v2/model"

	"github.com/elliot14A/fincher/internal/api/titles"
	"github.com/elliot14A/fincher/internal/clickhouse"
	"github.com/elliot14A/fincher/internal/scheduler"
	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/packages"
	tursotitles "github.com/elliot14A/fincher/internal/turso/titles"
	tursovendors "github.com/elliot14A/fincher/internal/turso/vendors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

func TestArmTitleDeadline_Matrix(t *testing.T) {
	conn, err := clickhouse.Open("127.0.0.1:9000")
	if err != nil {
		t.Skip("skipping integration: clickhouse connection failed:", err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = clickhouse.AutoMigrate(ctx, conn)

	tests := []struct {
		name          string
		initialStatus models.TitleStatus
		packageStates []models.PackageStatus
		expectBreach  bool
	}{
		{
			name:          "PROCESSING with defective/pending package -> emits breach event",
			initialStatus: models.StatusProcessing,
			packageStates: []models.PackageStatus{models.PackageStatusValid, models.PackageStatusInvalidated},
			expectBreach:  true,
		},
		{
			name:          "PROCESSING with all VALID packages -> ready, skips breach",
			initialStatus: models.StatusProcessing,
			packageStates: []models.PackageStatus{models.PackageStatusValid, models.PackageStatusValid},
			expectBreach:  false,
		},
		{
			name:          "SHIPPED title -> dedup guard skips breach",
			initialStatus: models.StatusShipped,
			packageStates: []models.PackageStatus{models.PackageStatusInvalidated},
			expectBreach:  false,
		},
		{
			name:          "OVERDUE title -> dedup guard skips breach",
			initialStatus: models.StatusOverdue,
			packageStates: []models.PackageStatus{models.PackageStatusInvalidated},
			expectBreach:  false,
		},
		{
			name:          "PROCESSING with zero packages -> unready, emits breach event",
			initialStatus: models.StatusProcessing,
			packageStates: []models.PackageStatus{},
			expectBreach:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client, err := turso.Open(":memory:", "")
			if err != nil {
				t.Fatalf("failed to open memory db: %v", err)
			}
			defer client.Close()
			_ = turso.AutoMigrate(context.Background(), client)

			_ = tursovendors.Create(context.Background(), client, &models.Vendor{
				Base:            models.Base{ID: "vendor-test"},
				Name:            "Test Vendor",
				Components:      []string{"AUDIO"},
				Markets:         []string{"de-DE"},
				HourlyRateUSD:   100.0,
				TurnaroundHours: 12,
			})

			slug := "deadline-" + uuid.NewString()[:8]
			titleID := "title-" + slug
			premiereDate := time.Now().UTC().Add(-1 * time.Second)

			_ = tursotitles.Create(context.Background(), client, &models.Title{
				Base:                 models.Base{ID: titleID},
				Name:                 "Deadline Title",
				Slug:                 slug,
				Type:                 models.TitleTypeFeature,
				PremiereDate:         premiereDate,
				Territories:          1,
				CurrentMasterVersion: "V01",
				OverallStatus:        tc.initialStatus,
			})

			for i, pStatus := range tc.packageStates {
				_ = packages.Create(context.Background(), client, &models.Package{
					Base:                     models.Base{ID: uuid.NewString()},
					TitleID:                  titleID,
					Component:                models.ComponentAudio,
					Language:                 "de-DE",
					Version:                  "V01",
					DerivedFromMasterVersion: "V01",
					VendorID:                 "vendor-test",
					Status:                   pStatus,
					RedeliveryCount:          i,
				})
			}

			// Clean events for this subject
			_, _ = conn.ExecContext(context.Background(), "alter table fincher.events delete where subject = ?", slug)

			timeScale := time.Millisecond
			sched := scheduler.NewScheduler(timeScale)
			defer sched.Stop()

			titleObj := &models.Title{
				Base:                 models.Base{ID: titleID},
				Slug:                 slug,
				PremiereDate:         premiereDate,
				OverallStatus:        tc.initialStatus,
				CurrentMasterVersion: "V01",
			}

			titles.ArmTitleDeadline(client, conn, func() model.LLM { return nil }, sched, titleObj)

			// Allow timer and callback to complete
			time.Sleep(100 * time.Millisecond)

			var breachCount int
			_ = conn.QueryRowContext(context.Background(),
				"select count() from fincher.events where subject = ? and type = ?",
				slug, models.TypeTitleDeadlineReached,
			).Scan(&breachCount)

			if tc.expectBreach && breachCount == 0 {
				t.Errorf("expected TypeTitleDeadlineReached event to be emitted into ClickHouse, got count 0")
			}
			if !tc.expectBreach && breachCount > 0 {
				t.Errorf("expected NO breach event, but got count %d", breachCount)
			}
		})
	}
}

func TestArmTitleDeadline_NilGuards(t *testing.T) {
	// Must not panic on nil arguments
	titles.ArmTitleDeadline(nil, nil, nil, nil, nil)

	timeScale := time.Millisecond
	sched := scheduler.NewScheduler(timeScale)
	defer sched.Stop()

	titles.ArmTitleDeadline(nil, nil, nil, sched, nil)

	title := &models.Title{
		Base:         models.Base{ID: "title-nil-test"},
		Slug:         "nil-test",
		PremiereDate: time.Now().UTC().Add(time.Hour),
	}
	titles.ArmTitleDeadline(nil, nil, nil, sched, title)
}
