package masters

import (
	"context"

	"github.com/elliot14A/fincher/internal/turso"
	"github.com/elliot14A/fincher/internal/turso/ent"
	entmaster "github.com/elliot14A/fincher/internal/turso/ent/master"
	domainerrors "github.com/elliot14A/fincher/pkg/domain/errors"
)

// Delete removes a master by ID and reconciles the title's current_master_version if needed.
func Delete(ctx context.Context, client *ent.Client, id string) domainerrors.Result[bool] {
	tx, err := client.Tx(ctx)
	if err != nil {
		return domainerrors.Err[bool](turso.NewError("masters.Delete", domainerrors.CodeInternal, "failed to start transaction", err))
	}
	defer tx.Rollback()

	m, err := tx.Master.Get(ctx, id)
	if err != nil {
		return domainerrors.Err[bool](turso.MapEntError("masters.Delete", "master", id, err))
	}

	title, err := tx.Title.Get(ctx, m.TitleID)
	if err == nil && title.CurrentMasterVersion == m.Version {
		latestRemaining, err := tx.Master.Query().
			Where(entmaster.TitleIDEQ(m.TitleID), entmaster.IDNEQ(id)).
			Order(ent.Desc(entmaster.FieldCreatedAt)).
			First(ctx)
		newVersion := ""
		if err == nil && latestRemaining != nil {
			newVersion = latestRemaining.Version
		}
		if _, updateErr := tx.Title.UpdateOneID(m.TitleID).SetCurrentMasterVersion(newVersion).Save(ctx); updateErr != nil {
			return domainerrors.Err[bool](turso.MapEntError("masters.Delete", "title", m.TitleID, updateErr))
		}
	}

	if err := tx.Master.DeleteOneID(id).Exec(ctx); err != nil {
		return domainerrors.Err[bool](turso.MapEntError("masters.Delete", "master", id, err))
	}

	if err := tx.Commit(); err != nil {
		return domainerrors.Err[bool](turso.NewError("masters.Delete", domainerrors.CodeInternal, "failed to commit transaction", err))
	}

	return domainerrors.Ok(true)
}
