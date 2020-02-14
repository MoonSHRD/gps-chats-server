package context

import (
	"github.com/MoonSHRD/sonis/internal/database"
	"github.com/MoonSHRD/sonis/internal/utils"
	"github.com/MoonSHRD/sonis/internal/webserver"
	"github.com/sirupsen/logrus"
)

const (
	DefaultWebserverPort = 23419
)

type Context struct {
	DB        *database.Database
	Webserver *webserver.Webserver
	Logger    *logrus.Logger
}

func New(cfg utils.Config) (*Context, error) {
	logger := logrus.New()
	db, err := database.New(cfg)
	if err != nil {
		logger.Error("Failed to initialize database. Reason: " + err.Error())
		return nil, err
	}
	if !utils.IsWebserverPortValid(&cfg) {
		cfg.WebserverPort = DefaultWebserverPort
	}
	ws, err := webserver.New(&cfg, db)
	if err != nil {
		logger.Error("Failed to initialize webserver. Reason: " + err.Error())
		return nil, err
	}
	return &Context{
		DB:        db,
		Webserver: ws,
	}, nil
}

func (ctx *Context) Destroy() {
	ctx.DB.CloseConnection()
	ctx.Webserver.Destroy()
}
