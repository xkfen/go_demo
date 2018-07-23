package mysql

import (
	"testing"
)

func TestDB(t *testing.T) {
	b := BaseModel{ID: 1}
	if b.ID != 1 {
		t.Fatal("id应该等于1")
	}
	conf := NewDbConfig()
	conf.DbName = "sdfsdfewf123"
	CreateDB(conf)
	DropDB(conf)
}
