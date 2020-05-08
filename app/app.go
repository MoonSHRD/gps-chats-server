package app

import (
	"fmt"

	_ "github.com/lib/pq"

	"github.com/MoonSHRD/sonis/app/migrations"
	"github.com/MoonSHRD/sonis/config"
	"github.com/google/logger"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	migrate "github.com/rubenv/sql-migrate"
)

type App struct {
	Config *config.Config
	DBConn *sqlx.DB
}

func NewApp(config config.Config) (*App, error) {
	app := &App{}
	dbConn, err := initDBConn(&config)
	if err != nil {
		return nil, err
	}
	app.Config = &config
	app.DBConn = dbConn
	return app, nil
}

func initDBConn(config *config.Config) (*sqlx.DB, error) {
	var err error
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s", config.PostgreSQL.User, config.PostgreSQL.Password, config.PostgreSQL.DatabaseName, config.PostgreSQL.Host, config.PostgreSQL.Port, "disable")
	dbConnection, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = initDBMigrations(dbConnection)
	if err != nil {
		logger.Errorf("Failed to process the migrations. Reason: %s", err.Error())
		return nil, err
	}
	return dbConnection, nil
}

func initDBMigrations(conn *sqlx.DB) error {
	_, err := migrate.Exec(conn.DB, "postgres", migrations.MigrationsList, migrate.Up)
	return err
}

func (app *App) Run(e *echo.Echo) {
	port := app.Config.HTTP.Port
	addr := app.Config.HTTP.Address
	logger.Infof("HTTP server starts listening at %s:%d", addr, port)
	logger.Fatal(e.Start(fmt.Sprintf("%s:%d", addr, port)))
}
