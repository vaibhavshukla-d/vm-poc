package configmanager

import (
	"vm/pkg/cinterface"

	"gorm.io/gorm"
)

type Config struct {
	App    *ApplicationConfigModal
	Db     *gorm.DB
	Logger cinterface.Logger
}
type ApplicationConfigModal struct {
	Application ApplicationConfig `mapstructure:"application_config"`
}

type ApplicationConfig struct {
	Application struct {
		Name                    string `mapstructure:"name"`
		Profile                 string `mapstructure:"profile"`
		Port                    string `mapstructure:"port"`
		ImageManagerServiceName string `mapstructure:"image_manager_service_name"`
		InfraMonitorServiceName string `mapstructure:"infra_monitor_service_name"`
		VmMonitorServiceName    string `mapstructure:"vm_monitor_service_name"`
		ValidateClientRequest   bool   `mapstructure:"validate_client_request"`
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
}
