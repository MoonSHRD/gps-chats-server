package context

import (
	"github.com/MoonSHRD/sonis/internal/database"
	"github.com/MoonSHRD/sonis/internal/webserver"
	"github.com/sirupsen/logrus"
)

type Context struct {
	DB        *database.Database
	Webserver *webserver.Webserver
	Logger    *logrus.Logger
}

func New() (*Context, error) {
	logger := logrus.New()
	db, err := database.New()
	if err != nil {
		logger.Error("Failed to initialize database. Reason: " + err.Error())
		return nil, err
	}
	ws, err := webserver.New(db)
	if err != nil {
		logger.Error("Failed to initialize webserver. Reason: " + err.Error())
		return nil, err
	}
	return &Context{
		DB:        db,
		Webserver: ws,
	}, nil
}
