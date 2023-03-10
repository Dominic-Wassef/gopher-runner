package dbrepo

import (
	"database/sql"

	"github.com/dominic-wassef/gopher-runner/internal/config"
	"github.com/dominic-wassef/gopher-runner/internal/repository"
)

var app *config.AppConfig

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewPostgresRepo creates the repository
func NewPostgresRepo(Conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	app = a
	return &postgresDBRepo{
		App: a,
		DB:  Conn,
	}
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewTestingRepo creates a repo with a dummy database for testing
func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
