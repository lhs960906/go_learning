package toml

import (
	"github.com/BurntSushi/toml"
)

type Host string
type Port int
type DatabaseName string
type User string
type Password string

type mysql struct {
	host         Host
	port         Port
	databaseName DatabaseName
	user         User
	password     Password
}

// 默认的 MySQL 配置
func newDefaultMySQL() *mysql {
	return &mysql{
		host:         "127.0.0.1",
		port:         3306,
		databaseName: "test",
		user:         "lhs",
		password:     "123456",
	}
}

func (m *mysql) GetMySQLHost() Host {
	return m.host
}

func (m *mysql) GetMySQLPort() Port {
	return m.port
}

func (m *mysql) GetMySQLDatabaseName() DatabaseName {
	return m.databaseName
}

func (m *mysql) GetMySQLUser() User {
	return m.user
}

func (m *mysql) GetMySQLPassword() Password {
	return m.password
}

type config struct {
	mysql
}

// 默认的 config 配置
func newDefaultConfig() *config {
	return &config{
		mysql: *newDefaultMySQL(),
	}
}

func (cfg *config) GetDB() *mysql {
	return &cfg.mysql
}

func Load(filepath string) {
	config := newDefaultConfig()
	toml.DecodeFile("example.toml", config)
}
