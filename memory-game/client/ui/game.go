package ui

import (
	"fmt"
	"time"

	"memory-master/shared"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type GameUI struct {
	window     fyne.Window
	game       shared.GameResponse
	cards      []*widget.Button
	flipped    []int
	matched    map[int]bool
	moves      int
	startTime  time.Time
	movesLabel *widget.Label
}

func NewGameUI(window fyne.Window, game shared.GameResponse) *GameUI {
	return &GameUI{
		window:    window,
		game:      game,
		matched:   make(map[int]bool),
		startTime: time.Now(),
		moves:     0,
	}
}

func (g *GameUI) Render(level shared.GameLevel) {
	g.movesLabel = widget.NewLabel(fmt.Sprintf("Ходы: %d", g.moves))

	grid := container.NewGridWithColumns(g.game.Columns)
	g.cards = make([]*widget.Button, len(g.game.Cards))

	for i := range g.game.Cards {
		index := i // Capture for closure
		btn := widget.NewButton("", func() {
			g.handleCardClick(index, level)
		})
		btn.Resize(fyne.NewSize(80, 80))
		g.cards[index] = btn
		grid.Add(btn)
	}

	backButton := widget.NewButton("Назад", func() {
		g.window.SetContent(createMainMenu(g.window))
	})

	content := container.NewVBox(
		g.movesLabel,
		grid,
		backButton,
	)

	g.window.SetContent(content)
}

func (g *GameUI) handleCardClick(index int, level shared.GameLevel) {
	if len(g.flipped) >= 2 || g.matched[index] {
		return
	}

	g.flipped = append(g.flipped, index)
	g.cards[index].SetText(g.game.Cards[index].Value)

	if len(g.flipped) == 2 {
		g.moves++
		g.movesLabel.SetText(fmt.Sprintf("Ходы: %d", g.moves))

		// Check for match
		if g.game.Cards[g.flipped[0]].Value == g.game.Cards[g.flipped[1]].Value {
			g.matched[g.flipped[0]] = true
			g.matched[g.flipped[1]] = true
			g.flipped = g.flipped[:0]
			// Check for win
			if len(g.matched) == len(g.game.Cards) {
				g.gameWon(level)
			}
		} else {
			// No match, flip back after delay
			go func() {
				time.Sleep(1 * time.Second)
				g.cards[g.flipped[0]].SetText("")
				g.cards[g.flipped[1]].SetText("")
				g.flipped = g.flipped[:0]
			}()
		}
	}
}

func (g *GameUI) gameWon(level shared.GameLevel) {
	duration := time.Since(g.startTime)
	score := g.calculateScore(int(duration.Seconds()))

	winContent := container.NewVBox(
		widget.NewLabel("Вы выиграли!"),
		widget.NewLabel(fmt.Sprintf("Ходы: %d", g.moves)),
		widget.NewLabel(fmt.Sprintf("Время: %.1f секунд", duration.Seconds())),
		widget.NewLabel(fmt.Sprintf("Очки: %d", score)),
	)
	if score > game.MaxScore {
		game.MaxScore = score
	}
	dialog.ShowCustom(
		"Победа!",
		"OK",
		winContent,
		g.window,
	)
	if level == shared.Easy {
		game.Аchives[0] = true
	}
	if level == shared.Medium && duration.Seconds() <= 180 {
		game.Аchives[1] = true
	}
	if level == shared.Hard && g.moves == 32 {
		game.Аchives[2] = true
	}
	if level == shared.Expert {
		game.Аchives[3] = true
	}
}

func (g *GameUI) calculateScore(timeSeconds int) int {
	// Базовые очки за уровень
	levelScores := map[shared.GameLevel]int{
		shared.Easy:   100,
		shared.Medium: 300,
		shared.Hard:   500,
		shared.Expert: 1000,
	}

	// Бонус за скорость (чем быстрее, тем больше очков)
	timeBonus := 0
	if timeSeconds < 60 {
		timeBonus = (60 - timeSeconds) * 10
	}

	// Штраф за лишние ходы
	perfectMoves := len(g.game.Cards) / 2
	movePenalty := 0
	if g.moves > perfectMoves {
		movePenalty = (g.moves - perfectMoves) * 5
	}

	return levelScores[g.game.Level] + timeBonus - movePenalty
}
