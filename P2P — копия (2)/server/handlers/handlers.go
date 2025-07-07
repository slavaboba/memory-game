package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"memory-master/shared"
	"net/http"
	"sync"
	"time"
)

var (
	activeGames = make(map[string]*shared.GameResponse)
	gameMutex   = sync.Mutex{}
)

func StartGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req shared.GameRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	game := createGame(req.Level)
	gameID := generateGameID()

	gameMutex.Lock()
	activeGames[gameID] = &game
	gameMutex.Unlock()

	// Save to database
	_, err = db.Exec("INSERT INTO games (id, player_name, level, start_time) VALUES (?, ?, ?, ?)",
		gameID, req.PlayerName, int(req.Level), time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		GameID string
		Game   shared.GameResponse
	}{
		GameID: gameID,
		Game:   game,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func MoveHandler(w http.ResponseWriter, r *http.Request) {
	// Implement move logic
}

func EndGameHandler(w http.ResponseWriter, r *http.Request) {
	// Implement game end logic
}

func GetScoresHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT player_name, score, level, time, moves FROM scores ORDER BY score DESC LIMIT 10")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var scores []shared.Score
	for rows.Next() {
		var s shared.Score
		err := rows.Scan(&s.PlayerName, &s.Score, &s.Level, &s.Time, &s.Moves)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		scores = append(scores, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scores)
}

func GetAchievementsHandler(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Query().Get("player")

	var rows *sql.Rows
	var err error

	if player != "" {
		rows, err = db.Query("SELECT name, description, unlocked FROM achievements WHERE player_name = ?", player)
	} else {
		rows, err = db.Query("SELECT name, description FROM achievements")
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var achievements []shared.Achievement
	for rows.Next() {
		var a shared.Achievement
		if player != "" {
			err := rows.Scan(&a.Name, &a.Description, &a.Unlocked)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err := rows.Scan(&a.Name, &a.Description)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		achievements = append(achievements, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(achievements)
}

func createGame(level shared.GameLevel) shared.GameResponse {
	var rows, cols int
	var cards []shared.Card

	switch level {
	case shared.Easy:
		rows, cols = 4, 4
	case shared.Medium:
		rows, cols = 6, 6
	case shared.Hard:
		rows, cols = 8, 8
	case shared.Expert:
		rows, cols = 10, 10
	}

	// Generate pairs
	pairs := rows * cols / 2
	for i := 0; i < pairs; i++ {
		// In a real app, you would have actual images or emojis
		value := string(rune(65 + i%26)) // A-Z
		cards = append(cards, shared.Card{ID: i * 2, Value: value})
		cards = append(cards, shared.Card{ID: i*2 + 1, Value: value})
	}

	// Shuffle cards
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })

	return shared.GameResponse{
		Cards:   cards,
		Level:   level,
		Rows:    rows,
		Columns: cols,
	}
}

func generateGameID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
