package handler

import (
	"encoding/json"
	"fmt"
	middlewareV1 "go-servie/internals/middleware"
	"go-servie/utils"
	"log"
	"net/http"
	"time"
)

const SecondsPerCoin = 30

func (authhandler *Handler) UserInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, err := middlewareV1.ExtractUser(r)
	if err != nil {
		utils.WriteJsonError(w, "Unauthorized", http.StatusUnauthorized, err)
		return
	}
	cacheKey := fmt.Sprintf("user:info:v1:%d", userId.UserId)
	cached, err := authhandler.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var result utils.User
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			log.Println("cahedValueOfUser --> ", result)
			utils.WriteJsonESuccess(w, "fetched user info", http.StatusOK, result)
			return
		}
	} else {
		log.Println("no cached  val: ", err)
	}
	user, err := authhandler.service.GetUserById(r.Context(), userId.UserId)
	if err != nil {
		utils.WriteJsonError(w, "failed to fetch user info", http.StatusInternalServerError, err)
		return
	}
	if data, err := json.Marshal(user); err == nil {
		err = authhandler.redis.SetEx(ctx, cacheKey, data, 5*time.Minute).Err()
		log.Println("setx: ", err)
	}
	utils.WriteJsonESuccess(w, "fetched user info", http.StatusOK, user)
}

func (h *Handler) GetUserCoinsHandler(w http.ResponseWriter, r *http.Request) {

	user, err := middlewareV1.ExtractUser(r)
	if err != nil {
		utils.WriteJsonError(w, "Unauthorized", http.StatusUnauthorized, err)
		return
	}

	ctx := r.Context()

	balance, err := h.service.GetUserCoinBalanceService(ctx, user.UserId)
	if err != nil {
		utils.WriteJsonError(w, "failed to fetch user coins", http.StatusInternalServerError, err)
		return
	}

	utils.WriteJsonESuccess(w, "success", http.StatusOK, map[string]any{
		"balance":          balance,
		"seconds_per_coin": SecondsPerCoin,
		"max_seconds":      balance * SecondsPerCoin,
	})
}
