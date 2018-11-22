package setting

import (
	. "easyweb/utils/tools"
	"github.com/go-ini/ini"
)

var (
	Config *ini.File
	// server
	Port           int
	ExpireDuration int
	// database
	DatabaseType     string
	DatabaseHost     string
	DatabasePort     int
	DatabaseUsername string
	DatabasePassword string
	Database         string
	MaxOpenConn      int
	MaxIdleConn      int
	// gin
	Mode string
)

func init() {
	var err error
	Config, err = ini.Load("conf/server.ini")
	CheckErr(err)
	// 初始化配置参数
	loadServerSetting()
	loadDatabaseSetting()
	loadGinSetting()
}

func loadServerSetting() {
	server, err := Config.GetSection("server")
	CheckErr(err)
	Port = server.Key("port").MustInt(8080)
	ExpireDuration = server.Key("session_expire_duration").MustInt(30)
}

func loadDatabaseSetting() {
	database, err := Config.GetSection("database")
	CheckErr(err)
	DatabaseType = database.Key("type").MustString("mysql")
	DatabaseHost = database.Key("host").MustString("localhost")
	DatabasePort = database.Key("port").MustInt(3306)
	DatabaseUsername = database.Key("username").String()
	DatabasePassword = database.Key("password").String()
	Database = database.Key("database").String()
	MaxOpenConn = database.Key("max_open_conn").MustInt(100)
	MaxIdleConn = database.Key("max_idle_conn").MustInt(10)
}

func loadGinSetting() {
	gin, err := Config.GetSection("server")
	CheckErr(err)
	Mode = gin.Key("mode").MustString("release")
}
