package configmanager

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"path/filepath"
	"sync"
	"time"
)

var (
	ErrConfigUnmarshal = errors.New("failed to unmarshal configuration")
	ErrConfigRead      = errors.New("failed to read configuration")
)

// ConfigLocker provides thread-safe access to configuration
type ConfigLocker struct {
	mutex sync.RWMutex
}

// ConfigViper holds the viper instance and provides safe access methods
type ConfigViper struct {
	viper *viper.Viper
	lock  sync.RWMutex
}

// UnmarshallViper unmarshalls configuration with additional error handling and retry logic
func (cl *ConfigLocker) UnmarshallViper(v *viper.Viper, configStruct interface{}, configName string) error {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	// Add retry logic for transient failures
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		err := v.Unmarshal(configStruct)
		if err == nil {
			return nil
		}
		if i == maxRetries-1 {
			return fmt.Errorf("%w: %s: %v", ErrConfigUnmarshal, configName, err)
		}
		time.Sleep(time.Millisecond * 100 * time.Duration(i+1))
	}
	return nil
}

// CreateViperInstance creates a new viper instance with improved error handling and validation
func CreateViperInstance(fileName, fileType, profile, serviceName, folderPath string) (*ConfigViper, error) {
	if fileName == "" || fileType == "" || profile == "" || serviceName == "" {
		return nil, errors.New("all parameters must be non-empty")
	}

	v := viper.New()
	v.SetConfigType(fileType)
	v.SetConfigName(fileName)

	// Add current directory first for faster lookup
	v.AddConfigPath(".")

	baseFolder := folderPath
	paths := buildConfigPaths(baseFolder, profile, serviceName)
	for _, path := range paths {
		v.AddConfigPath(path)
	}

	// Enable environment variable binding with prefix
	v.SetEnvPrefix(serviceName)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%w: %s: %v", ErrConfigRead, fileName, err)
	}

	return &ConfigViper{
		viper: v,
		lock:  sync.RWMutex{},
	}, nil
}

// Helper functions

func buildConfigPaths(baseFolder, profile, serviceName string) []string {
	paths := make([]string, 0, 5)
	current := filepath.Join(baseFolder, profile, serviceName)

	// Build paths up to 3 levels up
	for i := 0; i < 3; i++ {
		paths = append(paths, current)
		current = filepath.Join("..", current)
	}
	return paths
}

// ConfigViper methods

func (c *ConfigViper) GetViper() *viper.Viper {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.viper
}

func (c *ConfigViper) OnConfigChange(run func(in fsnotify.Event)) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.viper.OnConfigChange(run)
}

func (c *ConfigViper) WatchConfig() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.viper.WatchConfig()
}