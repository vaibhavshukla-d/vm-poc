package db

import (
	"fmt"
	"sync"
	"time"

	"vm/pkg/cinterface"
	configmanager "vm/pkg/config-manager"
	"vm/pkg/constants"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database interface {
	InitDB(config configmanager.ApplicationConfig) (*gorm.DB, error)
	GetReader() *gorm.DB
	Ping(db *gorm.DB) error
}

type DatabaseImpl struct {
	Logger       cinterface.Logger
	Db           *gorm.DB
	lastPingTime time.Time
	initOnce     sync.Once
}

func NewDatabase(logger cinterface.Logger) Database {
	return &DatabaseImpl{
		Logger: logger,
	}
}

func (rd *DatabaseImpl) InitDB(config configmanager.ApplicationConfig) (*gorm.DB, error) {
	var err error
	dsn := createDSN(config.Database)
	rd.Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := rd.Db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(config.Database.MaxOpenConnection)
	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConnection)
	sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(config.Database.MaxConnectionLifeTime))

	rd.lastPingTime = time.Now().UTC()
	return rd.Db, nil
}

func (rd *DatabaseImpl) GetReader() *gorm.DB {
	if rd.Db == nil || (rd.lastPingTime.Add(time.Minute*5).Before(time.Now().UTC()) && rd.Ping(rd.Db) != nil) {
		rd.Logger.Info(constants.MySql, constants.Startup, "Re-initializing database connection", nil)
	}
	return rd.Db
}

func (rd *DatabaseImpl) Ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	pingErr := sqlDB.Ping()
	if pingErr == nil {
		rd.lastPingTime = time.Now().UTC()
	}
	return pingErr
}

func createDSN(config configmanager.Database) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.Username, config.Password, config.Host, config.Port, config.DBName,
	)
}
