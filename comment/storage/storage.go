package storage

import (
	"fmt"
	"os"
	"sync"

	// driver for postgres
	"github.com/go-pg/pg"
)

var (
	// initOnce protects the following
	initOnce sync.Once
	db       *pg.DB
)

// NewDBHandler returns a db handler.
func NewDBHandler() *pg.DB {
	initOnce.Do(func() {
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_NAME")
		// psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
		db = pg.Connect(&pg.Options{
			Addr:     host + ":" + port,
			User:     user,
			Password: password,
			Database: dbName,
		})
		_, err := db.Exec("SELECT 1")
		if err != nil {
			fmt.Println("PostgreSQL is down")
		} else {
			fmt.Println("Successfully connected!")
		}
	})
	return db
}
