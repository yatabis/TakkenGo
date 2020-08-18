package database

import (
	"database/sql"
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

func GetQuestions() (id int, chapter, section string) {
	db := openDB()
	defer closeDB(db)
	row := db.QueryRow("select id, chapter, section from questions where id = (select (max(id) * random())::int from questions)")
	err := row.Scan(&id, &chapter, &section)
	if err != nil {
		log.Println(err)
	}
	return
}

func GetQuestionsById(id int) (chapter, section string) {
	db := openDB()
	defer closeDB(db)
	row := db.QueryRow("select chapter, section from questions where id = $1", id)
	err := row.Scan(&chapter, &section)
	if err != nil {
		log.Println(err)
	}
	return
}
