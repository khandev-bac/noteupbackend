package service

import (
	"context"
	"database/sql"
	"errors"
	"go-servie/utils"
	"log"

	"github.com/google/uuid"
)

func (authservice *Service) SignupService(ctx context.Context, email, password string, userDevice string) (*utils.AuthResponse, error) {
	_, err := authservice.repo.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("User already found")
	}
	// if errors.Is(err, sql.ErrNoRows) {
	// 	return nil, err
	// }
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user, err := authservice.repo.SignUpRepo(ctx, email, hashed, userDevice)
	if err != nil {
		return nil, err
	}
	err = authservice.repo.CreateUserCoinsRepo(ctx, user.ID, 2)
	if err != nil {
		log.Println("failed to createusercoin")
	}
	tokens, err := utils.TokenGeneration(user.ID)
	if err != nil {
		return nil, err
	}
	return &utils.AuthResponse{
		Tokens: utils.Token{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
		User_ID:     user.ID,
		Email:       user.Email,
		Plan:        user.Plan,
		User_Device: user.UserDevice.String,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (authservice *Service) Login(ctx context.Context, email, password string) (*utils.AuthResponse, error) {
	existuser, err := authservice.repo.GetUserByEmailLogin(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user Not found")
		}
		return nil, err // real DB error
	}
	if !existuser.IsActive {
		return nil, errors.New("account disabled")
	}
	if existuser.Password.Valid {
		if !utils.ComparePassword(existuser.Password.String, password) {
			return nil, errors.New("invalid credentials")
		}
	} else {
		return nil, errors.New("invalid credentials")
	}
	tokens, err := utils.TokenGeneration(existuser.ID)
	if err != nil {
		return nil, err
	}
	return &utils.AuthResponse{
		Tokens: utils.Token{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
		User_ID:     existuser.ID,
		Email:       existuser.Email,
		Plan:        existuser.Plan,
		User_Device: existuser.UserDevice.String,
		CreatedAt:   existuser.CreatedAt,
	}, nil
}

func (authservice *Service) GoogleAuth(ctx context.Context, email, google_id, picture, user_device string) (*utils.AuthResponse, error) {
	user, err := authservice.repo.GoogleAuthRepo(ctx, email, google_id, picture, user_device)
	if err != nil {
		log.Println("error in google auth service: ", err)
		return nil, err
	}
	tokens, err := utils.TokenGeneration(user.ID)
	if err != nil {
		log.Println("failed to create user token: ", err)
		return nil, err
	}
	err = authservice.repo.CreateUserCoinsRepo(ctx, user.ID, 2)
	log.Println("Successfully added coins to user: ", user.ID)
	if err != nil {
		log.Println("failed to createusercoin")
	}
	return &utils.AuthResponse{
		Tokens: utils.Token{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
		User_ID:     user.ID,
		Email:       user.Email,
		Plan:        user.Plan,
		User_Device: user.UserDevice.String,
		CreatedAt:   user.CreatedAt,
	}, nil
}
func (authservice *Service) GetUserById(ctx context.Context, id uuid.UUID) (*utils.User, error) {
	user, err := authservice.repo.GetUserById(ctx, id)
	return &utils.User{
		ID:         user.ID,
		Email:      user.Email,
		GoogleID:   &user.GoogleID.String,
		Picture:    &user.Picture.String,
		IsActive:   user.IsActive,
		Plan:       user.Plan,
		UserDevice: &user.UserDevice.String,
		CreatedAt:  user.CreatedAt,
	}, err
}
