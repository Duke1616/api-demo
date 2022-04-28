package conf

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/infraboard/mcube/logger/zap"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var global *Config

func C() *Config {
	if global == nil {
		panic("config required")
	}
	return global
}

func SetGlobalConfig(conf *Config) {
	global = conf
}

func NewDefaultConfig() *Config {
	return &Config{
		App:   newDefaultApp(),
		MySQL: newDefaultMysql(),
		Log:   newDefaultLog(),
	}
}

type Config struct {
	App   *app
	MySQL *mysql
	Log   *log
}

func newDefaultApp() *app {
	return &app{
		Name: "restful-api",
		Host: "127.0.0.1",
		Port: "8030",
		Key:  "123456",
	}
}

type app struct {
	Name string `toml:"name"`
	Host string `toml:"host"`
	Port string `toml:"port"`
	Key  string `toml:"key"`
}

func (a *app) Addr() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

func newDefaultMysql() *mysql {
	return &mysql{
		Host:        "127.0.0.1",
		Port:        "3306",
		Username:    "root",
		Password:    "123456",
		Database:    "go_course",
		MaxOpenConn: 100,
		MaxIdleConn: 20,
		MaxLifeTime: 600,
		MaxIdleTime: 300,
	}
}

type mysql struct {
	Host        string `toml:"host"`
	Port        string `toml:"port"`
	Username    string `toml:"username"`
	Password    string `toml:"password"`
	Database    string `toml:"database"`
	MaxOpenConn int    `toml:"max_open_conn"`
	MaxIdleConn int    `toml:"max_idle_conn"`
	MaxLifeTime int    `toml:"max_life_time"`
	MaxIdleTime int    `toml:"max_idle_time"`
	lock        sync.Mutex
}

var (
	db *sql.DB
)

func (m *mysql) getDBConn() (*sql.DB, error) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true", m.Username, m.Password, m.Host, m.Port, m.Database)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql<%s> error, %s", dsn, err.Error())
	}
	db.SetMaxOpenConns(m.MaxOpenConn)
	db.SetMaxIdleConns(m.MaxIdleConn)
	db.SetConnMaxLifetime(time.Second * time.Duration(m.MaxLifeTime))
	db.SetConnMaxIdleTime(time.Second * time.Duration(m.MaxIdleTime))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql<%s> error, %s", dsn, err.Error())
	}
	return db, nil
}

func (m *mysql) GetDB() (*sql.DB, error) {
	// 加载全局数据量单例
	m.lock.Lock()
	defer m.lock.Unlock()
	if db == nil {
		conn, err := m.getDBConn()
		if err != nil {
			return nil, err
		}
		db = conn
	}
	return db, nil
}

func newDefaultLog() *log {
	return &log{
		Level:  zap.DebugLevel.String(),
		Format: TextFormat,
		To:     ToStdout,
	}
}

type log struct {
	Level   string    `toml:"level" env:"LOG_LEVEL"`
	PathDir string    `toml:"path_dir" env:"LOG_PATH_DIR"`
	Format  LogFormat `toml:"format" env:"LOG_FORMAT"`
	To      LogTo     `toml:"to" env:"LOG_TO"`
}
