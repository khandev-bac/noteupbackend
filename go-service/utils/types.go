package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenOUT struct {
	UserId uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}
type Payload struct {
	UserId uuid.UUID `json:"user_id"`
}
type ErrorMessage struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statuscode"`
	Error      error  `json:"error"`
}
type SuccessMessage struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statuscode"`
	Data       any    `json:"data"`
}
type AuthResponse struct {
	Tokens      Token     `json:"tokens"`
	User_ID     uuid.UUID `json:"user_id"`
	Email       string    `json:"email"`
	Plan        string    `json:"plan"`
	User_Device string    `json:"user_device"`
	CreatedAt   time.Time `json:"created_at"`
}

type SignUpReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RefreshTokenBody struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type CreateNote struct {
	AudioUrl string `json:"audio_url" validate:"required"`
}

type IdTokenReq struct {
	IdToken string `json:"idToken"`
}
type User struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	GoogleID   *string   `json:"google_id"`
	Picture    *string   `json:"picture"`
	IsActive   bool      `json:"is_active"`
	Plan       string    `json:"plan"`
	UserDevice *string   `json:"user_device"`
	CreatedAt  time.Time `json:"created_at"`
}

type NoteRes struct {
	ID                   uuid.UUID   `json:"id"`
	UserID               uuid.UUID   `json:"user_id"`
	AudioUrl             string      `json:"audio_url"`
	AudioDurationSeconds int32       `json:"audio_duration_seconds"`
	AudioFileSizeMb      int32       `json:"audio_file_size_mb"`
	Transcript           string      `json:"transcript"`
	Title                string      `json:"title"`
	WordCount            int32       `json:"word_count"`
	Status               string      `json:"status"`
	SearchVector         interface{} `json:"search_vector"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
}

type NoteUpdate struct {
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
}

type NoteEvent struct {
	Event     string `json:"event"`
	NoteId    string `json:"note_id"`
	AudioUrl  string `json:"audio_url"`
	CreatedAt string `json:"created_at"`
}

type NoteEventResponse struct {
	Event      string `json:"event"`
	NoteId     string `json:"note_id"`
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
	CreatedAt  string `json:"created_at"`
}

type Coins struct {
	ID      uuid.UUID `json:"id"`
	Value   int       `json:"coin_value"`
	Price   int       `json:"coin_price"`
	Popular bool      `json:"popular"`
}
