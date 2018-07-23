package mysql
import (
	"io/ioutil"
	"strings"
	"gcoresys/common"
	"fmt"
)

type DbConfig struct {
	Username string
	Password string
	Host string
	Port string
	DbName string
	MaxIdleConns int
	MaxOpenConns int
}

// 只能通过这种方式获取配置对象
func NewDbConfig() *DbConfig {
	pwd := "1234"
	data, err := ioutil.ReadFile("/usr/local/.db/mysql.pas")
	if err != nil {
		fmt.Println("读取mysql密码文件出错:" + err.Error())
	}else{
		pwd = string(data)
		pwd = strings.TrimSpace(pwd)
	}
	appConf := common.GetAppConfig()
	conf := &DbConfig{
		Username: appConf.MysqlUname,
		Password: pwd,
		Host: appConf.MysqlHost,

		Port: appConf.MysqlPort,
		MaxIdleConns: 10,
		MaxOpenConns: 100,
	}
	return conf
}

// 复制一份配置
func (dbConfig *DbConfig) Clone() (*DbConfig) {
	return &DbConfig {
		Username: dbConfig.Username,
		Password: dbConfig.Password,
		Host: dbConfig.Host,
		Port: dbConfig.Port,
		DbName: dbConfig.DbName,
		MaxIdleConns: dbConfig.MaxIdleConns,
		MaxOpenConns: dbConfig.MaxOpenConns,
	}
}
