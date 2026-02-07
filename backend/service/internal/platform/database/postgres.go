package database

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

func NewPostgres() (*sql.DB, error) {
    dsn := "postgres://tennisleague:tennisleague@localhost:5432/tennisleague?sslmode=disable"
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