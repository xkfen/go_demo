package mongodb

import (
	"gopkg.in/mgo.v2"
	"time"
	"io/ioutil"
	"fmt"
	"strings"
	"go_demo/util"
)

func NewDbConfig() *mgo.DialInfo {
	uname := "root"
	pwd := "root"
	pwdB, pwdErr := ioutil.ReadFile("/usr/local/.db/mongo.pas")
	unameB, unameErr := ioutil.ReadFile("/usr/local/.db/mongo.uname")
	if unameErr != nil { fmt.Println("读取mongo用户名文件出错:" + unameErr.Error()) }
	if pwdErr != nil { fmt.Println("读取mongo用户名文件出错:" + pwdErr.Error()) }
	if unameErr == nil && pwdErr == nil {
		uname = strings.TrimSpace(string(unameB))
		pwd = strings.TrimSpace(string(pwdB))
	}
	appConf := util.GetAppConfig()
	return &mgo.DialInfo{
		Addrs: []string{appConf.MongoHost},
		// 数据库的指定需要仔细参悟
		Database: "admin",
		Username:  uname,
		Password:  pwd,
		Direct:    false,
		Timeout:   time.Second * 1,
		PoolLimit: 100, // Session.SetPoolLimit
	}
}