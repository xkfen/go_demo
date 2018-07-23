package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 创建数据库 information_schema
func CreateDB(dbConfig *DbConfig) {
	if !checkConfig(dbConfig) {
		panic("配置不正确")
	}
	cStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, "information_schema")
	openedDb, err := gorm.Open("mysql", cStr)
	if err != nil {
		fmt.Println(cStr)
		panic("连接数据库出错:" + err.Error())
	}

	createDbSQL := "CREATE DATABASE IF NOT EXISTS " + dbConfig.DbName + " DEFAULT CHARSET utf8 COLLATE utf8_general_ci;"

	err = openedDb.Exec(createDbSQL).Error
	if err != nil {
		fmt.Println("创建失败：" + err.Error() + " sql:" + createDbSQL)
		return
	}
	fmt.Println(dbConfig.DbName + "数据库创建成功")
}

// 删除数据库
func DropDB(dbConfig *DbConfig) {
	if !checkConfig(dbConfig) {
		panic("配置不正确")
	}

	openedDb, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, "information_schema"))
	if err != nil {
		panic("连接数据库出错:" + err.Error())
	}

	dropDbSQL := "DROP DATABASE IF EXISTS " + dbConfig.DbName + ";"

	err = openedDb.Exec(dropDbSQL).Error
	if err != nil {
		fmt.Println("删除失败：" + err.Error() + " sql:" + dropDbSQL)
		return
	}
	fmt.Println(dbConfig.DbName + "数据库删除成功")
}

func checkConfig(config *DbConfig) bool {
	if len(config.DbName) == 0 {
		fmt.Println("db_name 不能为空")
		return false
	}
	if len(config.Username) == 0 {
		fmt.Println("username 不能为空")
		return false
	}
	if len(config.Password) == 0 {
		fmt.Println("password 不能为空")
		return false
	}
	if len(config.Host) == 0 {
		fmt.Println("host 不能为空")
		return false
	}
	if len(config.Port) == 0 {
		fmt.Println("port 不能为空")
		return false
	}
	return true
}
