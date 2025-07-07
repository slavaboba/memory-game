package main

import (
	"log"
	"memory-master/server/handlers"
	"net/http"
)

func main() {
	// Initialize database
	handlers.InitDB()

	// Setup routes
	http.HandleFunc("/api/game/start", handlers.StartGameHandler)
	http.HandleFunc("/api/game/move", handlers.MoveHandler)
	http.HandleFunc("/api/game/end", handlers.EndGameHandler)
	http.HandleFunc("/api/scores", handlers.GetScoresHandler)
	http.HandleFunc("/api/achievements", handlers.GetAchievementsHandler)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
