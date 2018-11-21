package setting

import (
	. "easyweb/utils/tools"
	"github.com/go-ini/ini"
)

var (
	Config *ini.File

	Port int

	DatabaseType     string
	DatabaseHost     string
	DatabasePort     int
	DatabaseUsername string
	DatabasePassword string
)

func init() {
	var err error
	Config, err = ini.Load("conf/server.ini")
	CheckErr(err)
	Port = Config.Section("server").Key("port").MustInt(8080)

	loadDatabaseSetting()
}

func loadDatabaseSetting() {
	database, err := Config.GetSection("database")
	CheckErr(err)
	DatabaseType = database.Key("type").MustString("mysql")
	DatabaseHost = database.Key("host").MustString("localhost")
	DatabasePort = database.Key("port").MustInt(3306)
	DatabaseUsername = database.Key("username").String()
	DatabasePassword = database.Key("password").String()
}
