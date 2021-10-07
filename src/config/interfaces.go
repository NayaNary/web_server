package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	User       string
	Password   string
	DbName     string
	SslMode    string
	NameDriver string
	TimeWrite  uint64
}

type WebConfig struct {
	Host string
	Port string
}

type Config struct {
	Db  DbConfig
	Web WebConfig
}

func New() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	return &Config{
		Db:  dbConfig(),
		Web: webCongig(),
	}
}

func dbConfig() (data DbConfig) {
	data.DbName = getEnv("DB_NAME","test_task")
	data.NameDriver = getEnv("DRIVER_NAME","postgres")
	data.Password = getEnv("PASSWORD", "123456789")
	data.User = getEnv("DB_USER", "postgres")
	data.SslMode = getEnv("SSL_MODE", "disable")
	timeWrite := getEnv("TIME_WRITE","30")
	var err error
	data.TimeWrite, err = strconv.ParseUint(timeWrite, 10, 0) 
	if err != nil{
		log.Println("Time write not entered:",err.Error())
	}
	return
}

func webCongig() (data WebConfig) {
	data.Host = getEnv("HOST","localhost")
	data.Port = getEnv("PORT","5000")
	return
}

func getEnv(key string, defaultVal string) string {
    if value, exists := os.LookupEnv(key); exists {
	return value
    }

    return defaultVal
}