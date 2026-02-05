package middlewareV1

import (
	"context"
	"errors"
	"go-servie/utils"
	"log"
	"net/http"
	"strings"
)

type ctxkey string

const UserKey ctxkey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
			utils.WriteJsonError(w, "Unthorized", http.StatusUnauthorized, nil)
			return
		}
		token := strings.Split(authHeader, " ")[1]
		verifyAccessToken, err := utils.VerifyAccessToken(token)
		if err != nil {
			utils.WriteJsonError(w, "Unthorized", http.StatusUnauthorized, nil)
			log.Println("middleware error: ", err)
			return
		}
		ctx := context.WithValue(r.Context(), UserKey, verifyAccessToken)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
func ExtractUser(r *http.Request) (*utils.Payload, error) {
	user, ok := r.Context().Value(UserKey).(*utils.Payload)
	if !ok || user == nil {
		return nil, errors.New("unauthorized")
	}
	return user, nil
}
