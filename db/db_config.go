package db

import (
	"io/ioutil"
	"fmt"
	"strings"
	"go_demo/util"
	"go_demo/mysql"
	"flag"
)
const (
	DefaultReadUname = "db_read"
	DefaultReadPwd   = "db_read"
)
// 每次使用时clone一下
var curMysqlConf *mysql.DbConfig

// 获取只读配置（返回clone出去，否则要加锁）
func getMysqlReadConf() *mysql.DbConfig {
	if curMysqlConf != nil {
		return curMysqlConf.Clone()
	}
	curMysqlConf = mysql.NewDbConfig()
	if !util.IsTestEnv() {
		curMysqlConf.Username = DefaultReadUname
		curMysqlConf.Password = DefaultReadPwd
	}
	//curMysqlConf.Host = "172.16.1.90"
	// 如果在k8s中配置了，则读取配置
	if unameData, uErr := ioutil.ReadFile("/usr/local/.db/mysql_read.uname"); uErr != nil {
		fmt.Println("读取mysql用户名报错:" + uErr.Error())
	} else {
		curMysqlConf.Username = strings.TrimSpace(string(unameData))
	}
	if pwdData, pErr := ioutil.ReadFile("/usr/local/.db/mysql_read.pas"); pErr != nil {
		fmt.Println("读取mysql密码报错:" + pErr.Error())
	} else {
		curMysqlConf.Password = strings.TrimSpace(string(pwdData))
	}
	//fmt.Println("数据库链接:" + curMysqlConf.Host)
	return curMysqlConf.Clone()
}



// 根据运行环境获取数据库名称
func getDbName(dbNamePre string) string {
	if util.IsTestEnv() {
		return dbNamePre + "_test"
	} else {
		envF := flag.Lookup("env")
		if envF != nil && envF.Value != nil {
			switch envF.Value.String() {
			case "test":
				return dbNamePre + "_test"
			case "prod":
				return dbNamePre + "_prod"
			default:
				return dbNamePre + "_dev"
			}
		}
		return dbNamePre + "_prod"
	}
}