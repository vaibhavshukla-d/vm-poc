package config

// import (
// 	"vm/pkg/configmanager"
// 	"fmt"
// 	"github.com/fsnotify/fsnotify"
// 	"sync/atomic"
// )



// type Holder struct {
// 	AppCfg         atomic.Value // holds *configmanager.ApplicationConfigModal
// }

// var (
// 	configHolder = &Holder{}
// 	configLock   = &configmanager.ConfigLocker{}
// )

// func Load() error {

// 	// Initialize application config
// 	appConfigViper, err := configmanager.CreateViperInstance("config", "json", "dev", "app", ".")
// 	if err != nil {
// 		return fmt.Errorf("failed to create app config: %w", err)
// 	}

// 	appCfg := &configmanager.ApplicationConfigModal{}
// 	if err := configLock.UnmarshallViper(appConfigViper.GetViper(), appCfg, "config"); err != nil {
// 		return fmt.Errorf("failed to unmarshal app config: %w", err)
// 	}
// 	configHolder.AppCfg.Store(appCfg)

// 	appConfigViper.OnConfigChange(func(e fsnotify.Event) {
// 		if err := configLock.UnmarshallViper(appConfigViper.GetViper(), appCfg, "config"); err != nil {
// 			panic("failed to unmarshal app config : " + err.Error())
// 		}
// 		configHolder.AppCfg.Store(appCfg)
// 	})
// 	appConfigViper.WatchConfig()



// 	return nil
// }

// // Getter methods

// func GetAppConfig() *configmanager.ApplicationConfigModal {
// 	return configHolder.AppCfg.Load().(*configmanager.ApplicationConfigModal)
// }

// func GetConfig() *configmanager.Config {
// 	err := Load()
// 	if err != nil {
// 		panic("failed to load config: " + err.Error())
// 	}
// 	config := configmanager.Config{
// 		App:         configHolder.AppCfg.Load().(*configmanager.ApplicationConfigModal),
// 	}
// 	return &config
// }