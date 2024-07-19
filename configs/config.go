package configs

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	dbStringConnection = "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4"
)

type configYAML struct {
	// App config
	AppName  string `yaml:"appname"`
	HttpPort string `yaml:"httpport"`

	DB dbYAML `yaml:"db"`
}

type dbConfigYAML struct {
	Host    string `yaml:"host"`
	Maxidle int    `yaml:"maxidle"`
	Maxopen int    `yaml:"maxopen"`
	Name    string `yaml:"name"`
	Pass    string `yaml:"pass"`
	Port    string `yaml:"port"`
	User    string `yaml:"user"`
}

type dbYAML struct {
	Master dbConfigYAML `yaml:"master"`
}

type Config struct {
	// app
	AppName  string
	HttpPort string

	CronCheckLoanStatusSchedule string

	DBMaster *sqlDatabase
}

type sqlDatabase struct {
	ConnectionString string
	MaxIdle          int
	MaxOpen          int
	User             string
	Password         string
	Host             string
	Port             string
	Name             string
}

var appConfig *Config

func Init(serviceName string) {
	cfg := &configYAML{}
	logrus.Info("loading config from local config")
	if err := loadConfigFromLocalFile(cfg); err != nil {
		panic(fmt.Sprintf("failed reading local config, err: %v", err))
	}

	// App config
	appConfig = &Config{}
	appConfig.initCommonConfig(cfg)
	appConfig.initSqlDBConfig(cfg)
}

func Get() *Config {
	return appConfig
}

func loadConfigFromLocalFile(c *configYAML) error {
	yamlFile, err := os.ReadFile("./configs/app.yaml")
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(yamlFile, c); err != nil {
		return err
	}

	if c.HttpPort == "" {
		c.HttpPort = "8080"
	}
	return nil
}

func (c *Config) initCommonConfig(cfg *configYAML) {
	c.AppName = cfg.AppName

	c.HttpPort = cfg.HttpPort
	if c.HttpPort == "" {
		c.HttpPort = "8080"
	}
}

func (c *Config) initSqlDBConfig(cfg *configYAML) {
	appConfig.DBMaster = &sqlDatabase{
		User:     cfg.DB.Master.User,
		Password: cfg.DB.Master.Pass,
		Host:     cfg.DB.Master.Host,
		Port:     cfg.DB.Master.Port,
		Name:     cfg.DB.Master.Name,
		ConnectionString: fmt.Sprintf(
			dbStringConnection,
			cfg.DB.Master.User,
			cfg.DB.Master.Pass,
			cfg.DB.Master.Host,
			cfg.DB.Master.Port,
			cfg.DB.Master.Name,
		),
		MaxIdle: cfg.DB.Master.Maxidle,
		MaxOpen: cfg.DB.Master.Maxopen,
	}
}
