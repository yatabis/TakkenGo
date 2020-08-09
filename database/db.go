package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func openDB() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Printf("failed to open db: %e\n", err)
	}
	return db
}

func closeDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("failed to close db: %e\n", err)
	}
}

func GetQuestions() (chapter, section string) {
	db := openDB()
	defer closeDB(db)
	row := db.QueryRow("select chapter, section from questions where id = (select (max(id) * random())::int from questions)")
	err := row.Scan(&chapter, &section)
	if err != nil {
		fmt.Println(err)
	}
	return
}
