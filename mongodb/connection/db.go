package connection
import (
	"gopkg.in/mgo.v2"
	"sync"
	"strings"
)

var session *mgo.Session
var mutex sync.Mutex

// 获取数据库连接
// 外边传进来的Database要是其要连的db，这里会先尝试去连admin（设置readWriteAnyDatabase就在admin），如果不成功则可能是只针对对应的数据库做了授权，因此再尝试去连对应数据库，如果对应数据库也失败了，那么就是真的用户名和密码错误了
func GetDb(dbConfig *mgo.DialInfo) *mgo.Session {
	mutex.Lock()
	defer mutex.Unlock()
	if session == nil {
		//logger.Info("初始化mongo db session")
		tmpDbName := dbConfig.Database
		// 先用admin尝试
		dbConfig.Database = "admin"
		var err error
		if session, err = mgo.DialWithInfo(dbConfig); err != nil {
			// 出错的话有可能是管理员只给该用户开了该数据库的权限
			dbConfig.Database = tmpDbName
			if session, err = mgo.DialWithInfo(dbConfig); err != nil {
				panic("mongodb连接报错:" + err.Error())
			}
		}
		session.SetMode(mgo.Strong, true)
	}
	return session.Clone()
}

// 清空某个数据库下的所有数据
func ClearAllData(dbConfig *mgo.DialInfo) {
	if strings.Contains(dbConfig.Database, "test") {
		// 获取连接
		tmpS := GetDb(dbConfig)
		tmpDb := tmpS.DB(dbConfig.Database)
		cName, _ := tmpDb.CollectionNames()
		for _, cn := range cName {
			tmpDb.C(cn).DropCollection()
		}
	} else {
		//logger.Warn("非法操作！在非测试环境下调用了清空所有数据的方法")
	}
}

// 关闭连接
func CloseDb(dbConfig *mgo.DialInfo) {
	mutex.Lock()
	defer mutex.Unlock()
	if session != nil {
		session.Close()
		session = nil
	}
}