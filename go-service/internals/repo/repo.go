package repo

import "go-servie/dbmodel"

type Repo struct {
	db *dbmodel.Queries
}

func NewRepository(db *dbmodel.Queries) *Repo {
	return &Repo{
		db: db,
	}
}
