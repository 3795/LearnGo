package conf

const DriverName = "mysql"

type DBConfig struct {
	Host 	string
	Port 	int
	User 	string
	Pwd		string
	Database 	string
	IsRunning 	bool
}

var DBMasterList = []DBConfig {
	{
		Host:      "127.0.0.1",
		Port:      3306,
		User:      "root",
		Pwd:       "root",
		Database:  "lottery",
		IsRunning: true,
	},
}

var DBMaster DBConfig = DBMasterList[0]