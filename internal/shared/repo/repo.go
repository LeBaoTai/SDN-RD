package repo

import "gorm.io/gorm"

type Repo struct {
	dbConn *gorm.DB
}

func NewRepo(conn *gorm.DB) *Repo {
	return &Repo{
		dbConn: conn,
	}
}

func (r *Repo) GetAllDevice() {
}
