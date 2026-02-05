package repo

import (
	"context"
	"database/sql"
	"go-servie/dbmodel"

	"github.com/google/uuid"
)

func (authrepo *Repo) SignUpRepo(ctx context.Context, email string, password string, user_device string) (dbmodel.SignupRow, error) {
	return authrepo.db.Signup(ctx, dbmodel.SignupParams{
		Email:      email,
		Password:   sql.NullString{String: password, Valid: true},
		UserDevice: sql.NullString{String: user_device, Valid: true},
	})
}

func (authrepo *Repo) GoogleAuthRepo(ctx context.Context, email, google_id, picture, user_device string) (dbmodel.GoogleAuthRow, error) {
	return authrepo.db.GoogleAuth(ctx, dbmodel.GoogleAuthParams{
		Email:      email,
		GoogleID:   sql.NullString{String: google_id, Valid: true},
		Picture:    sql.NullString{String: picture, Valid: true},
		UserDevice: sql.NullString{String: user_device, Valid: true},
	})
}

func (authrepo *Repo) GetUserById(ctx context.Context, id uuid.UUID) (dbmodel.GetUserByIdRow, error) {
	return authrepo.db.GetUserById(ctx, id)
}

func (authrepo *Repo) GetUserByEmail(ctx context.Context, email string) (dbmodel.GetUserByEmailRow, error) {
	return authrepo.db.GetUserByEmail(ctx, email)
}
func (authrepo *Repo) GetUserByEmailLogin(ctx context.Context, email string) (dbmodel.User, error) {
	return authrepo.db.GetUserByEmailLogin(ctx, email)
}
