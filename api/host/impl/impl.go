package impl

import (
	"database/sql"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/Duke1616/api-demo/conf"
)

var Service *impl = &impl{}

type impl struct {
	log logger.Logger
	db  *sql.DB
}

func (i *impl) Init() error {
	i.log = zap.L().Named("Host")

	db, err := conf.C().MySQL.GetDB()
	if err != nil {
		return err
	}
	i.db = db
	return nil
}
