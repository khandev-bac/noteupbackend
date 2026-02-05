package handler

import (
	"encoding/json"
	"go-servie/utils"
	"log"
	"net/http"
)

func (authhander *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var req *utils.SignUpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJsonError(w, "Invalid request", http.StatusBadRequest, err)
		return
	}
	user, err := authhander.service.SignupService(r.Context(), req.Email, req.Password, r.Header.Get("User-Agent"))
	if err != nil {
		log.Println("error in signup-service: ", err)
		utils.WriteJsonError(w, err.Error(), http.StatusInternalServerError, err)
		return
	}
	utils.WriteJsonESuccess(w, "Successfully authenticated", http.StatusCreated, user)
}
func (authhander *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req *utils.SignUpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJsonError(w, "Invalid request", http.StatusBadRequest, err)
		return
	}
	user, err := authhander.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		utils.WriteJsonError(w, err.Error(), http.StatusInternalServerError, err)
		return
	}
	utils.WriteJsonESuccess(w, "Login Successfully", http.StatusOK, user)
}

func (authhandler *Handler) Test(w http.ResponseWriter, r *http.Request) {
	utils.WriteJsonESuccess(w, "Health ok", http.StatusOK, nil)
}

func (authhandler *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req *utils.RefreshTokenBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJsonError(w, "Invalid request", http.StatusBadRequest, err)
		return
	}
	payload, err := utils.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		utils.WriteJsonError(w, "invalid token or verification failed", http.StatusUnauthorized, err)
		return
	}
	tokens, err := utils.TokenGeneration(payload.UserId)
	if err != nil {
		utils.WriteJsonError(w, "verification failed", http.StatusUnauthorized, err)
		return
	}
	utils.WriteJsonESuccess(w, "Verification successfull", http.StatusOK, utils.Token{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (authhandler *Handler) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	var req *utils.IdTokenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJsonError(w, "Invalid request", http.StatusBadRequest, err)
		return
	}
	authclient, err := authhandler.firebaseApp.Auth(r.Context())
	if err != nil {
		log.Println("firebase err in authhandler: ", err)
		utils.WriteJsonError(w, "something went wrong", http.StatusInternalServerError, err)
		return
	}
	userdata, err := authclient.VerifyIDToken(r.Context(), req.IdToken)
	if err != nil {
		utils.WriteJsonError(w, "google auth failed", http.StatusInternalServerError, err)
		return
	}
	google_id := userdata.UID
	email := userdata.Claims["email"].(string)
	picture := userdata.Claims["picture"].(string)
	user_device := r.Header.Get("User-Agent")
	user, err := authhandler.service.GoogleAuth(r.Context(), email, google_id, picture, user_device)
	if err != nil {
		log.Println("failed to save in db: ", err)
		utils.WriteJsonError(w, "something went wrong", http.StatusInternalServerError, err)
		return
	}
	utils.WriteJsonESuccess(w, "Authentication successfull", http.StatusOK, user)
}
