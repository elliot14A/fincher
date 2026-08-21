package masters

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
	"github.com/elliot14A/fincher/pkg/domain/models"
)

// Create inserts a new master version and updates the title's current_master_version.
func Create(ctx context.Context, client *ent.Client, m *models.Master) domainerrors.Result[*models.Master] {
	if err := m.Validate(); err != nil {
		return domainerrors.Err[*models.Master](turso.NewError("masters.Create", domainerrors.CodeInvalidInput, "invalid master data", err))
	}

	// Run within transaction to keep Master insertion and Title current_master_version in sync
	tx, err := client.Tx(ctx)
	if err != nil {
		return domainerrors.Err[*models.Master](turso.NewError("masters.Create", domainerrors.CodeInternal, "failed to begin transaction", err))
	}

	builder := tx.Master.Create().
		SetID(m.ID).
		SetTitleID(m.TitleID).
		SetVersion(m.Version)

	if m.SupersedesVersion != "" {
		builder.SetSupersedesVersion(m.SupersedesVersion)
	}

	created, err := builder.Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		return domainerrors.Err[*models.Master](turso.MapEntError("masters.Create", "master", m.ID, err))
	}

	// Update parent title current_master_version
	_, err = tx.Title.UpdateOneID(m.TitleID).
		SetCurrentMasterVersion(m.Version).
		Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		return domainerrors.Err[*models.Master](turso.MapEntError("masters.Create", "title", m.TitleID, err))
	}

	if err := tx.Commit(); err != nil {
		return domainerrors.Err[*models.Master](turso.NewError("masters.Create", domainerrors.CodeInternal, "failed to commit transaction", err))
	}

	return domainerrors.Ok(toDomain(created))
}
