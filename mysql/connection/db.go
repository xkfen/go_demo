package connection

import (
	"fmt"
	"go_demo/mysql"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gcoresys/common/logger"
	"sync"
	"time"
	"gcoresys/common/util"
)

var conns map[string]*gorm.DB = make(map[string]*gorm.DB)
var mutex sync.Mutex

func GetDb(dbConfig *mysql.DbConfig) *gorm.DB {
	mutex.Lock()
	defer mutex.Unlock()
	tmpDb := conns[dbConfig.DbName]
	if tmpDb == nil {
		logger.Info("初始化数据库连接：", "db_host", dbConfig.Host, "db_name", dbConfig.DbName, "user", dbConfig.Username)
		openedDb, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DbName))
		if err != nil {
			panic("数据库连接出错：" + err.Error())
		}
		// 在上面赋值会有坑，不信试试，呵呵
		tmpDb = openedDb
		tmpDb.DB().SetMaxIdleConns(dbConfig.MaxIdleConns)
		tmpDb.DB().SetMaxOpenConns(dbConfig.MaxOpenConns)
		// 避免久了不使用，导致连接被mysql断掉的问题
		tmpDb.DB().SetConnMaxLifetime(time.Hour * 2)
		// 如果不是生产数据库则打开详细日志
		//if !strings.Contains(dbConfig.DbName, "prod") {
		if util.Substr(dbConfig.DbName, len(dbConfig.DbName)-4, 4) != "prod" {
			tmpDb.LogMode(true)
		}
		conns[dbConfig.DbName] = tmpDb
	}
	return tmpDb
}

// 清空指定数据库下的所有数据，只在测试情况下执行
func ClearAllData(dbConfig *mysql.DbConfig) {
	if strings.Contains(dbConfig.DbName, "test") {
		tmpDb := conns[dbConfig.DbName]
		if tmpDb == nil {
			fmt.Println("尚未初始化数据库")
			return
		}
		if rs, err := tmpDb.Raw("show tables;").Rows(); err == nil {
			var tName string
			for rs.Next() {
				rs.Scan(&tName)
				if tName != "" {
					tmpDb.Exec(fmt.Sprintf("delete from %s", tName))
				}
			}
		}
	} else {
		logger.Warn("非法操作！在非测试环境下调用了清空所有数据的方法")
	}
}

func CloseDb(dbConfig *mysql.DbConfig) {
	mutex.Lock()
	defer mutex.Unlock()
	tmpDb := conns[dbConfig.DbName]
	if tmpDb != nil {
		tmpDb.Close()
		conns[dbConfig.DbName] = nil
	}
}
