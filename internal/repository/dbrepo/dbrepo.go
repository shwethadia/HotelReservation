package dbrepo

import (
	"database/sql"

	"github.com/shwethadia/HotelReservation/internal/config"
	"github.com/shwethadia/HotelReservation/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type postgresDBTestRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {

	return &postgresDBRepo{

		App: a,
		DB:  conn,
	}

}

func NewPostgresTestRepo(a *config.AppConfig) repository.DatabaseRepo {

	return &postgresDBTestRepo{

		App: a,
	}

}
