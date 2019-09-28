package core

const (
	COOKIE_AUTH_KEY = "_ak"

	USERNAME_KEY = "_username"

	PASSWORD_KEY = "_password"

	DEFAULT_SERVER_PORT = 6010

	TABLE_PREFIX = "tank30_"

	VERSION = "3.0.5"
)

type Config interface {
	Installed() bool

	ServerPort() int

	// 数据库地址
	MysqlUrl() string

	// 文件存储路径
	MatterPath() string

	FinishInstall(mysqlPort int, mysqlHost string, mysqlSchema string, mysqlUsername string, mysqlPassword string)
}
