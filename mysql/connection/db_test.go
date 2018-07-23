package connection

import (
	"testing"
	"go_demo/mysql"
)

func TestGetDb(t *testing.T) {
	conf := mysql.NewDbConfig()
	conf.DbName = "hahahah121312321"
	mysql.CreateDB(conf)
	GetDb(conf)
	CloseDb(conf)
	mysql.DropDB(conf)
}