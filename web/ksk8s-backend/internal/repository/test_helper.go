package repository

import (
	"testing"

	"github.com/easzlab/ksk8s/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB creates an in-memory SQLite database for testing.
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// Auto-migrate all models
	models := []interface{}{
		&model.User{},
		&model.Cluster{},
		&model.ClusterNode{},
		&model.ClusterParam{},
		&model.ClusterParamList{},
		&model.HostsGlobalVars{},
		&model.ClusterVersion{},
		&model.Template{},
		&model.Task{},
		&model.Log{},
		&model.RunningTask{},
		&model.Audit{},
		&model.UserClusterBinding{},
	}
	for _, m := range models {
		if err := db.AutoMigrate(m); err != nil {
			t.Fatalf("failed to migrate %T: %v", m, err)
		}
	}

	// Set global DB for repository functions
	DB = db
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})

	return db
}
