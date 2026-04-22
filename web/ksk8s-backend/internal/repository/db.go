package repository

import (
	"fmt"
	"os"

	"github.com/easzlab/ksk8s/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	// SQLite mode for dev/QA
	if os.Getenv("KSK8S_USE_SQLITE") == "1" {
		dbPath := getEnv("KSK8S_SQLITE_PATH", "/tmp/ksk8s.db")
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to open sqlite: %w", err)
		}
		DB = db
		return db, nil
	}

	host := getEnv("KSK8S_DB_HOST", "localhost")
	port := getEnv("KSK8S_DB_PORT", "3306")
	user := getEnv("KSK8S_DB_USER", "ksk8s")
	pass := getEnv("KSK8S_DB_PASSWORD", "ksk8s_pass")
	dbname := getEnv("KSK8S_DB_NAME", "ksk8s")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	return db, nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(
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
		&model.UserClusterBinding{},
		&model.SystemSetting{},
		&model.Audit{},
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
