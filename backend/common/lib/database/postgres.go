package database

import (
	"database/sql"
	"fmt"

	"tennis-league/common/lib/config"

	_ "github.com/lib/pq"
)

func NewPostgres() (*sql.DB, error) {

	config := config.LoadPostgresConfig()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Dbname,
		config.SSLMode,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Postgres connected")
	return db, nil
}
