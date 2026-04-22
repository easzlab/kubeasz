package main

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate [up|down]")
		os.Exit(1)
	}

	dbHost := getEnv("KSK8S_DB_HOST", "localhost")
	dbPort := getEnv("KSK8S_DB_PORT", "3306")
	dbUser := getEnv("KSK8S_DB_USER", "ksk8s")
	dbPass := getEnv("KSK8S_DB_PASSWORD", "ksk8s_pass")
	dbName := getEnv("KSK8S_DB_NAME", "ksk8s")

	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "migrate init failed: %v\n", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			fmt.Fprintf(os.Stderr, "migrate up failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations applied successfully")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			fmt.Fprintf(os.Stderr, "migrate down failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations rolled back successfully")
	default:
		fmt.Println("Usage: migrate [up|down]")
		os.Exit(1)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
