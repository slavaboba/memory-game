package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"memory-master/shared"
	"net/http"
)

const serverURL = "http://localhost:8080"

func StartGame(player string, level shared.GameLevel) (*shared.GameResponse, error) {
	req := shared.GameRequest{
		PlayerName: player,
		Level:      level,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(serverURL+"/api/game/start", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var gameResp struct {
		GameID string
		Game   shared.GameResponse
	}
	err = json.Unmarshal(body, &gameResp)
	if err != nil {
		return nil, err
	}

	return &gameResp.Game, nil
}

func GetScores() ([]shared.Score, error) {
	resp, err := http.Get(serverURL + "/api/scores")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var scores []shared.Score
	err = json.NewDecoder(resp.Body).Decode(&scores)
	return scores, err
}

func GetAchievements(player string) ([]shared.Achievement, error) {
	resp, err := http.Get(serverURL + "/api/achievements?player=" + player)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var achievements []shared.Achievement
	err = json.NewDecoder(resp.Body).Decode(&achievements)
	return achievements, err
}
