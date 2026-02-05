package repo

import (
	"context"

	"go-servie/dbmodel"

	"github.com/google/uuid"
)

/* =========================
   COIN PACKS REPO
========================= */

func (r *Repo) GetActiveCoinPacksRepo(ctx context.Context) ([]dbmodel.GetActiveCoinPacksRow, error) {
	return r.db.GetActiveCoinPacks(ctx)
}

func (r *Repo) GetCoinPackByIdRepo(ctx context.Context, packID uuid.UUID) (dbmodel.CoinPack, error) {
	return r.db.GetCoinPackById(ctx, packID)
}
