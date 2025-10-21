package dependency

import (
	"context"
	"os"
	"strconv"
	"vm/pkg/cinterface"
	configmanager "vm/pkg/config-manager"
	"vm/pkg/constants"
	"vm/pkg/db"
	"vm/pkg/logger"
)

type Dependency struct {
	Ctx context.Context
	*ClientDependency
	Logger   cinterface.Logger
	Database db.Database
	Config   *configmanager.Config
}

// getEnv returns the environment variable or default value if not set
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt reads an env variable and converts it to int
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func Setup(ctx context.Context) (*Dependency, error) {
	// Load DB config from environment
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnvInt("DB_PORT", 3306)
	dbName := getEnv("DB_NAME", "vmdb")
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASS", "test@12")
	maxIdle := getEnvInt("DB_MAX_IDLE", 10)
	maxOpen := getEnvInt("DB_MAX_OPEN", 100)
	connLife := getEnvInt("DB_CONN_LIFETIME", 60)

	// Load App config from environment
	appPort := getEnv("APP_PORT", "8080")
	imageManagerServiceName := getEnv("IMAGE_MANAGER_SERVICE_NAME", "image-manager:8081")
	infraMonitorServiceName := getEnv("INFRA_MONITOR_SERVICE_NAME", "infra-monitor:8082")
	vmMonitorServiceName := getEnv("VM_MONITOR_SERVICE_NAME", "vm-monitor:8083")

	// Build configuration
	cfg := &configmanager.Config{
		App: &configmanager.ApplicationConfigModal{
			Application: struct {
				Application struct {
					Name                    string `mapstructure:"name"`
					Profile                 string `mapstructure:"profile"`
					Port                    string `mapstructure:"port"`
					ImageManagerServiceName string `mapstructure:"image_manager_service_name"`
					InfraMonitorServiceName string `mapstructure:"infra_monitor_service_name"`
					VmMonitorServiceName    string `mapstructure:"vm_monitor_service_name"`
				} `mapstructure:"app"`
				Database struct {
					Host                  string `mapstructure:"host"`
					Port                  int    `mapstructure:"port"`
					DBName                string `mapstructure:"dbName"`
					Username              string `mapstructure:"username"`
					Password              string `mapstructure:"password"`
					MaxIdleConnection     int    `mapstructure:"maxIdleConnection"`
					MaxOpenConnection     int    `mapstructure:"maxOpenConnection"`
					MaxConnectionLifeTime int    `mapstructure:"MaxConnectionLifeTime"`
				} `mapstructure:"database"`
				Log struct {
					Level         string `mapstructure:"Level"`
					FilePath      string `mapstructure:"FilePath"`
					FileName      string `mapstructure:"FileName"`
					Encoding      string `mapstructure:"Encoding"`
					EnableConsole bool   `mapstructure:"EnableConsole"`
					EnableFile    bool   `mapstructure:"EnableFile"`
				} `mapstructure:"log"`
			}{
				Application: struct {
					Name                    string `mapstructure:"name"`
					Profile                 string `mapstructure:"profile"`
					Port                    string `mapstructure:"port"`
					ImageManagerServiceName string `mapstructure:"image_manager_service_name"`
					InfraMonitorServiceName string `mapstructure:"infra_monitor_service_name"`
					VmMonitorServiceName    string `mapstructure:"vm_monitor_service_name"`
				}{
					Name:                    "vm",
					Profile:                 "dev",
					Port:                    appPort,
					ImageManagerServiceName: imageManagerServiceName,
					InfraMonitorServiceName: infraMonitorServiceName,
					VmMonitorServiceName:    vmMonitorServiceName,
				},
				Database: struct {
					Host                  string `mapstructure:"host"`
					Port                  int    `mapstructure:"port"`
					DBName                string `mapstructure:"dbName"`
					Username              string `mapstructure:"username"`
					Password              string `mapstructure:"password"`
					MaxIdleConnection     int    `mapstructure:"maxIdleConnection"`
					MaxOpenConnection     int    `mapstructure:"maxOpenConnection"`
					MaxConnectionLifeTime int    `mapstructure:"MaxConnectionLifeTime"`
				}{
					Host:                  dbHost,
					Port:                  dbPort,
					DBName:                dbName,
					Username:              dbUser,
					Password:              dbPass,
					MaxIdleConnection:     maxIdle,
					MaxOpenConnection:     maxOpen,
					MaxConnectionLifeTime: connLife,
				},
				Log: struct {
					Level         string `mapstructure:"Level"`
					FilePath      string `mapstructure:"FilePath"`
					FileName      string `mapstructure:"FileName"`
					Encoding      string `mapstructure:"Encoding"`
					EnableConsole bool   `mapstructure:"EnableConsole"`
					EnableFile    bool   `mapstructure:"EnableFile"`
				}{
					Level:         "debug",
					EnableConsole: true,
				},
			},
		},
	}

	// Initialize logger
	log := logger.NewLogger(cfg)
	log.Info(constants.General, constants.Startup, "Logger initialized", nil)

	// Initialize client dependencies
	clientDeps, err := SetupClientDependencies(cfg, log)
	if err != nil {
		return nil, err
	}

	// Initialize database
	database := db.NewDatabase(log)
	_, err = database.InitDB(*cfg.App)
	if err != nil {
		log.Error(constants.MySql, constants.Startup, "Failed to connect to database", map[constants.ExtraKey]interface{}{"error": err})
		return nil, err
	}
	log.Info(constants.MySql, constants.Startup, "Database connection established", nil)

	// Run migrations
	migrator := db.NewMigrateAllTables(database, log)
	if err := migrator.MigrateAllTables(); err != nil {
		log.Error(constants.MySql, constants.Migration, "Failed to run migrations", map[constants.ExtraKey]interface{}{"error": err})
		os.Exit(1)
	}
	log.Info(constants.MySql, constants.Migration, "Migrations completed", nil)

	return &Dependency{
		Ctx:              ctx,
		ClientDependency: clientDeps,
		Logger:           log,
		Database:         database,
		Config:           cfg,
	}, nil
}
