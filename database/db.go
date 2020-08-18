package database

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"

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

func GetScore(time int) (score int) {
	db := openDB()
	defer closeDB(db)
	row := db.QueryRow("select score from scores where time = $1", time)
	if err := row.Scan(&score); err != nil {
		log.Println(err)
	}
	return
}

func UpdateScore(time, score int) error {
	db := openDB()
	defer closeDB(db)
	_, err := db.Exec("update scores set score = $1 where time = $2", score, time)
	if err != nil {
		log.Println(err)
		return errors.New("スコアの保存に失敗しました。")
	}
	return nil
}

func UpdateRate(id, score int) error {
	db := openDB()
	defer closeDB(db)
	var rate, count float64
	row := db.QueryRow("select rate, count from questions where id = $1", id)
	err := row.Scan(&rate, &count)
	if err != nil {
		log.Println(err)
		return errors.New("スコアの保存に失敗しました。")
	}
	rate = (rate * count + float64(score)) / (count + 1)
	_, err = db.Exec("update questions set rate = $1, count = $2 where id = $3", rate, count + 1, id)
	if err != nil {
		log.Println(err)
		return errors.New("スコアの保存に失敗しました。")
	}
	return nil
}

func SaveScore(id, time, score int) error {
	if GetScore(time) != 0 {
		return errors.New(strconv.Itoa(time) + "時のスコアはすでに登録されています。")
	}
	if err := UpdateScore(time, score); err != nil {
		return err
	}
	if err := UpdateRate(id, score); err != nil {
		return err
	}
	return nil
}

func ResetScores() {
	db := openDB()
	defer closeDB(db)
	if _, err := db.Exec("update scores set score = 0"); err != nil {
		log.Println(err)
	}
}
