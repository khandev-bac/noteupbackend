package handler

import (
	"net/http"

	"go-servie/utils"
)

func (h *Handler) GetCoinPacksHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	packs, err := h.service.GetActiveCoinPacksService(ctx)
	if err != nil {
		utils.WriteJsonError(w, "failed to fetch coin packs", http.StatusInternalServerError, err)
		return
	}

	var result []map[string]any
	for _, p := range packs {
		result = append(result, map[string]any{
			"id":         p.ID,
			"coin_value": p.CoinValue,
			"price":      p.CoinPrice, // paise (frontend converts)
			"popular":    p.Popular.Bool,
		})
	}

	utils.WriteJsonESuccess(w, "success", http.StatusOK, result)
}
