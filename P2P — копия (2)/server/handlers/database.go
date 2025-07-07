package handlers

import (
	"database/sql"
	"log"
	"sync"

	"memory-master/shared"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

func InitDB() {
	once.Do(func() {
		var err error
		db, err = sql.Open("sqlite3", "./memory.db")
		if err != nil {
			log.Fatal(err)
		}

		createTables()
		seedData()
	})
}

func createTables() {
	// Games table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS games (
			id TEXT PRIMARY KEY,
			player_name TEXT,
			level INTEGER,
			start_time DATETIME,
			end_time DATETIME,
			completed BOOLEAN
		)`)
	if err != nil {
		log.Fatal(err)
	}

	// Scores table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS scores (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_id TEXT,
			player_name TEXT,
			score INTEGER,
			level INTEGER,
			time INTEGER,
			moves INTEGER,
			FOREIGN KEY(game_id) REFERENCES games(id)
		)`)
	if err != nil {
		log.Fatal(err)
	}

	// Achievements table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS achievements (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			player_name TEXT,
			name TEXT,
			description TEXT,
			unlocked BOOLEAN,
			unlocked_at DATETIME
		)`)
	if err != nil {
		log.Fatal(err)
	}
}

func seedData() {
	// Predefined achievements
	achievements := []shared.Achievement{
		{Name: "Новичок", Description: "Завершить лёгкий уровень"},
		{Name: "Знаток", Description: "Пройти средний уровень за 3 минуты"},
		{Name: "Мастер памяти", Description: "Пройти сложный уровень без ошибок"},
		{Name: "Легенда", Description: "Победить на экспертном уровне"},
		{Name: "Скорострел", Description: "Открыть 3 пары подряд за 5 секунд"},
	}

	for _, a := range achievements {
		_, err := db.Exec("INSERT OR IGNORE INTO achievements (name, description) VALUES (?, ?)", a.Name, a.Description)
		if err != nil {
			log.Println("Error seeding achievements:", err)
		}
	}
}
