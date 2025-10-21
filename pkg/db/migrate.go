package db

import (
	"fmt"

	"vm/internal/modals"
	"vm/pkg/cinterface"
	"vm/pkg/constants"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Migrate interface {
	MigrateAllTables() error
}

type MigrateImpl struct {
	db     Database
	Logger cinterface.Logger
}

func NewMigrateAllTables(db Database, logger cinterface.Logger) Migrate {
	return &MigrateImpl{
		db:     db,
		Logger: logger,
	}
}

func (t *MigrateImpl) MigrateAllTables() error {
	migrationID := uuid.New().String()
	db := t.db.GetReader()

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: migrationID,
			Migrate: func(tx *gorm.DB) error {
				entities := []interface{}{
					&modals.VMRequest{},
					&modals.VMDeployInstance{},
				}

				for _, entity := range entities {
					if err := tx.AutoMigrate(entity); err != nil {
						t.Logger.Error(constants.MySql, constants.Migration, "Migration failed for entity", map[constants.ExtraKey]interface{}{
							"Entity":      fmt.Sprintf("%T", entity),
							"MigrationID": migrationID,
							"Error":       err,
						})
						return tx.Rollback().Error
					}
					t.Logger.Info(constants.MySql, constants.Migration, "Migration successful for entity", map[constants.ExtraKey]interface{}{
						"Entity":      fmt.Sprintf("%T", entity),
						"MigrationID": migrationID,
					})
				}

				t.Logger.Info(constants.MySql, constants.Migration, "All migrations applied successfully", map[constants.ExtraKey]interface{}{
					"MigrationID": migrationID,
				})
				return nil
			},
		},
	})

	if err := m.Migrate(); err != nil {
		t.Logger.Error(constants.MySql, constants.Migration, "Migration process failed", map[constants.ExtraKey]interface{}{
			"MigrationID": migrationID,
			"Error":       err,
		})
		return fmt.Errorf("migration process failed for Migration ID %s: %w", migrationID, err)
	}

	t.Logger.Info(constants.MySql, constants.Migration, "Migration completed successfully", map[constants.ExtraKey]interface{}{
		"MigrationID": migrationID,
	})
	return nil
}
