package configs

import (
	"fmt"
	"net"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Logger    LoggerConf
	StorageDB bool
	Storage   StorageConf
}

type LoggerConf struct {
	Level logger.LogLevel
}

type StorageConf struct {
	Address  net.Addr
	User     string
	Password string
	DBName   string
	SslMode  string
}

func NewConfig(configFile string) Config {
	var err error
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Configuration is not loaded, default values will be used")
	}

	logLevel := viper.GetString("logger.level")
	logValid := logLevelValidator(logLevel)
	if logLevel == "" || !logValid {
		fmt.Println("log automatic set to default")
		logLevel = "INFO"
	}
	storageDB := viper.GetString("storage.type") == "DB"
	storage := StorageConf{}

	if storageDB {
		host := viper.GetString("storage.host")
		port := viper.GetString("storage.port")
		address := net.JoinHostPort(host, port)
		user := viper.GetString("storage.user")
		password := viper.GetString("storage.password")
		dbName := viper.GetString("storage.dbname")
		sslmode := viper.GetString("storage.sslmode")
		if user == "" || password == "" || dbName == "" || sslmode == "" {
			fmt.Println("One of credits DB not exist, using default")
			user = "postgres"
			password = "postgres"
			dbName = "events_db"
			sslmode = "disable"
		}
		storage.User = user
		storage.Password = password
		storage.SslMode = sslmode
		storage.DBName = dbName
		storage.Address, err = net.ResolveTCPAddr("tcp", address)
		if err != nil {
			fmt.Println("host or port incorrect, using default")
			storage.Address, _ = net.ResolveIPAddr("ip", "127.0.0.1:5432")
		}
	}

	// TODO:  Тут надо в конфиге проверять есть ли нстройки бд и если нет то сторедж ин мемори

	return Config{
		Logger:    LoggerConf{Level: logger.LogLevel(logLevel)},
		StorageDB: storageDB,
		Storage:   storage,
	}
}

func logLevelValidator(level string) bool {
	allowLevel := map[string]bool{
		string(logger.DebugLevel): true,
		string(logger.WarnLevel):  true,
		string(logger.InfoLevel):  true,
		string(logger.ErrorLevel): true,
	}
	return allowLevel[level]
}
