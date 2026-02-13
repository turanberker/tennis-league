package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/turanberker/tennis-league-service/internal/platform"
)

func NewPostgres() (*sql.DB, error) {

	config := platform.LoadPostgresConfig()
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
