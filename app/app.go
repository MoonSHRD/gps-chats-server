package app

import (
	"context"
	"fmt"

	migrate "github.com/xakep666/mongo-migrate"

	//_ "github.com/MoonSHRD/sonis/app/migrations"
	"github.com/MoonSHRD/logger"
	"github.com/MoonSHRD/sonis/config"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	Config       *config.Config
	MongoClient  *mongo.Client
	MainDatabase *mongo.Database
}

func NewApp(config config.Config) (*App, error) {
	app := &App{}
	mongoClient, mongoDB, err := initMongoConnection(&config)
	if err != nil {
		return nil, err
	}
	app.Config = &config
	app.MongoClient = mongoClient
	app.MainDatabase = mongoDB
	return app, nil
}

func initMongoConnection(config *config.Config) (*mongo.Client, *mongo.Database, error) {
	var mongoURI string
	if config.MongoDB.User == "" && config.MongoDB.Password == "" {
		mongoURI = fmt.Sprintf("mongodb://%s:%d", config.MongoDB.Host, config.MongoDB.Port)
	} else {
		mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%d", config.MongoDB.User, config.MongoDB.Password, config.MongoDB.Host, config.MongoDB.Port)
	}
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, nil, err
	}

	db := client.Database(config.MongoDB.DatabaseName)
	err = initDBMigrations(db)
	if err != nil {
		logger.Errorf("Failed to process the migrations. Reason: %s", err.Error())
		return nil, nil, err
	}
	return client, db, nil
}

func initDBMigrations(db *mongo.Database) error {
	migrate.SetDatabase(db)
	return migrate.Up(migrate.AllAvailable)
}

func (app *App) Run(e *echo.Echo) {
	port := app.Config.HTTP.Port
	addr := app.Config.HTTP.Address
	logger.Infof("HTTP server starts listening at %s:%d", addr, port)
	logger.Fatal(e.Start(fmt.Sprintf("%s:%d", addr, port)))
}
