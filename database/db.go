package database

import (
	"database/sql"
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

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

func choiceLowest(list []float64) int {
	n := len(list)
	cum := make([]float64, n)
	cum[0] = 100 - list[0]
	for i := 1; i < n; i++ {
		cum[i] = cum[i-1] + 100 - list[i]
	}
	rand.Seed(time.Now().UnixNano())
	max := cum[n-1]
	r := rand.Float64() * max
	for i, c := range cum {
		if c > r {
			return i + 1
		}
	}
	return -1
}

func GetRatesList() (ratesList []float64) {
	db := openDB()
	defer closeDB(db)
	rows, err := db.Query("select  rate from questions order by id")
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var rate float64
		if err := rows.Scan(&rate); err != nil {
			log.Println(err)
			return nil
		}
		ratesList = append(ratesList, rate)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil
	}
	return
}

func GetQuestionByRate() (id int, chapter, section string, rate float64) {
	ratesList := GetRatesList()
	if ratesList == nil {
		return
	}
	id = choiceLowest(ratesList)
	chapter, section, rate = GetQuestionById(id)
	return
}

func GetQuestionById(id int) (chapter, section string, rate float64) {
	db := openDB()
	defer closeDB(db)
	row := db.QueryRow("select chapter, section, rate from questions where id = $1", id)
	err := row.Scan(&chapter, &section, &rate)
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
	rate = (rate*count + float64(score)) / (count + 1)
	_, err = db.Exec("update questions set rate = $1, count = $2 where id = $3", rate, count+1, id)
	if err != nil {
		log.Println(err)
		return errors.New("スコアの保存に失敗しました。")
	}
	return nil
}

func SaveScore(id, time, score int) error {
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
