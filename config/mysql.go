package config

type MySQL struct {
	Host     string `json:"host"`
	User     string `json:"userapi"`
	Pwd      string `json:"pwd"`
	Database string `json:"database"`
	Port     string `json:"port"`
}

// GetMySQLConfig 初始化MySQL配置
func GetMySQLConfig() *MySQL {
	return &MySQL{
		Host:     "localhost",
		User:     "root",
		Pwd:      "123456",
		Database: "k8s",
		Port:     "3306",
	}
}
