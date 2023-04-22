package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/k4zb3k/project/pkg/logger"
	"sync"
)

type Config struct {
	IsDebug      *bool              `yaml:"is_debug" env-required:"true"`
	Listen       ListenConfig       `yaml:"listen"`
	DatabaseConn DatabaseConnConfig `yaml:"database_conn"`
	CacheConn    CacheConnConfig    `yaml:"cache_conn"`
	BrokerConn   BrokerConnConfig   `yaml:"broker_conn"`
	JwtConfig    JWTConfig          `yaml:"jwt_config"`
}

type ListenConfig struct {
	Type   string `yaml:"type" env-default:"port"`
	BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
	Port   string `yaml:"port" env-default:"8080"`
}

type DatabaseConnConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Dbname   string `json:"dbname"`
	User     string `json:"user"`
	Password string `json:"password"`
	Sslmode  string `json:"sslmode"`
}

type CacheConnConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

type BrokerConnConfig struct {
	Url string `json:"url"`
}

type JWTConfig struct {
	AccessSecret  string `json:"access_secret"`
	RefreshSecret string `json:"refresh_secret"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		logger.Info.Println("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("/home/k4zb3k/Desktop/project/config/config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info.Println(help)
			logger.Error.Fatalln(err)
		}
	})
	return instance
}
