package conf

type SystemUserIds struct {
	AddContactUserId string
}

type MysqlConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	DbName   string
}

type Config struct {
	HttpHostPort             string
	TcpHostPort              string
	MysqlConfig              *MysqlConfig
	PrintLog                 bool
	StaticFileSystemPathRoot string
	StaticHttpRootPath       string
	SystemUserIds            SystemUserIds
}

var gConfig *Config

func GetConfig() *Config {
	return gConfig
}

var configDev *Config = &Config{
	HttpHostPort: "0.0.0.0:8080",
	MysqlConfig: &MysqlConfig{
		Username: "root",
		Password: "123456",
		Host:     "127.0.0.1",
		Port:     3306,
		DbName:   "lite_im",
	},
	StaticFileSystemPathRoot: `C:\Users\nerd\catlog\code\LiteIM\server\tmp\%s`,
	StaticHttpRootPath:       `http://127.0.0.1:8080/static/download/%s`,
	SystemUserIds: SystemUserIds{
		AddContactUserId: "AddContactUserId",
	},
}

func InitConfig(env string) {
	switch env {
	case "dev":
		gConfig = configDev
	default:
		panic("not support env")
	}
}
