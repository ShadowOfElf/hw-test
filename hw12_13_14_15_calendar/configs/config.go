package configs

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/logger"
	"github.com/spf13/viper"
)

type AppTypes string

const (
	MainApp      AppTypes = "main"
	SchedulerApp AppTypes = "scheduler"
	SenderApp    AppTypes = "sender"
)

type Config struct {
	Logger     LoggerConf
	StorageDB  bool
	Storage    StorageConf
	HTTP       HTTPConf
	GRPC       GRPCConf
	AppType    AppTypes
	Rabbit     RabbitMQConf
	UpdateTime time.Duration
}

type RabbitMQConf struct {
	User     string
	Password string
	Host     string
	Port     string
}

type GRPCConf struct {
	Addr string
}

type HTTPConf struct {
	Addr string
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

	appType := viper.GetString("app.type")
	if !appTypeValidator(appType) {
		fmt.Println("app type automatic set to main")
		appType = "main"
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
			fmt.Println("host or port DB incorrect, using default")
			storage.Address, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:5432")
		}
	}

	httpHost := viper.GetString("http.host")
	httpPort := viper.GetString("http.port")

	addr := net.JoinHostPort(httpHost, httpPort)
	_, err = net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println("host or port HTTP incorrect, using default")
		addr = "127.0.0.1:8070"
	}

	grpcHost := viper.GetString("grpc.host")
	grpcPort := viper.GetString("grpc.port")

	addrGRPC := net.JoinHostPort(grpcHost, grpcPort)
	_, err = net.ResolveTCPAddr("tcp", addrGRPC)
	if err != nil {
		fmt.Println("host or port GRPC incorrect, using default")
		addr = "127.0.0.1:8070"
	}

	updTime := 1
	if updTimeStr := viper.GetString("app.upd"); updTimeStr != "" {
		if newUpdTime, err := strconv.Atoi(updTimeStr); err == nil {
			updTime = newUpdTime
		}
	}

	return Config{
		Logger:     LoggerConf{Level: logger.LogLevel(logLevel)},
		StorageDB:  storageDB,
		Storage:    storage,
		HTTP:       HTTPConf{Addr: addr},
		GRPC:       GRPCConf{Addr: addrGRPC},
		AppType:    AppTypes(appType),
		Rabbit:     getRabbitConf(),
		UpdateTime: time.Duration(updTime) * time.Minute,
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

func appTypeValidator(appType string) bool {
	allowType := map[string]bool{
		string(MainApp):      true,
		string(SchedulerApp): true,
		string(SenderApp):    true,
	}
	return allowType[appType]
}

func getRabbitConf() RabbitMQConf {
	user := viper.GetString("rabbit.user")
	password := viper.GetString("rabbit.password")
	host := viper.GetString("rabbit.host")
	port := viper.GetString("rabbit.port")

	if host == "" || port == "" {
		fmt.Println("host or port is empty, use default")
		host = "127.0.0.1"
		port = "5672"
	}

	return RabbitMQConf{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	}
}
