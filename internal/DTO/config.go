package dto

type Config struct {
	Database struct {
		Host     string `config:"host"`
		Port     int    `config:"port"`
		Username string `config:"username"`
		Password string `config:"password"`
		Database string `config:"database"`
	} `config:"database"`

	HttpServer struct {
		Address string `config:"address"`
	} `config:"http-server"`

	LogLevel string `config:"log-level"`
}
