package main

import (
	"fmt"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Logger    LoggerConf
	StorageDB bool
}

type LoggerConf struct {
	Level logger.LogLevel
}

func NewConfig() Config {
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

	// TODO:  Тут надо в конфиге проверять есть ли нстройки бд и если нет то сторедж ин мемори

	return Config{
		Logger: LoggerConf{Level: logger.LogLevel(logLevel)},
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
