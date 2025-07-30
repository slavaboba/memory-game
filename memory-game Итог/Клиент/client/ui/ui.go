package ui

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"memory-master/client/internal"
	"memory-master/shared"
)

var myApp fyne.App
var achivments [5]bool
var PlNames [10]string
var PlScores [10]int
var MaxScore int

func ShowMainWindow() {
	myApp = app.New()
	window := myApp.NewWindow("Memory Master")
	window.SetContent(createMainMenu(window))
	window.Resize(fyne.NewSize(800, 600))
	window.Show()
	myApp.Run()
}

func createMainMenu(window fyne.Window) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Memory Master", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	startButton := widget.NewButton("Новая игра", func() {
		window.SetContent(createLevelSelection(window))
	})

	scoresButton := widget.NewButton("Рейтинг игроков", func() {
		showScores(window)
	})

	achievementsButton := widget.NewButton("Достижения", func() {
		showAchievements(window)
	})

	exitButton := widget.NewButton("Выход", func() {
		exit()
	})

	return container.NewVBox(
		title,
		startButton,
		scoresButton,
		achievementsButton,
		exitButton,
	)
}

func createLevelSelection(window fyne.Window) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Выберите уровень сложности", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	easyButton := widget.NewButton("Лёгкий (4×4)", func() {
		startGame(window, shared.Easy)
	})

	mediumButton := widget.NewButton("Средний (6×6)", func() {
		startGame(window, shared.Medium)
	})

	hardButton := widget.NewButton("Сложный (8×8)", func() {
		startGame(window, shared.Hard)
	})

	expertButton := widget.NewButton("Эксперт (10×10)", func() {
		startGame(window, shared.Expert)
	})

	backButton := widget.NewButton("Назад", func() {
		window.SetContent(createMainMenu(window))
	})

	return container.NewVBox(
		title,
		easyButton,
		mediumButton,
		hardButton,
		expertButton,
		backButton,
	)
}

func startGame(window fyne.Window, level shared.GameLevel) {
	// Здесь будет реализация начала игры
	game := shared.GameResponse{
		Level:   level,
		Rows:    getRowsForLevel(level),
		Columns: getColumnsForLevel(level),
		Cards:   generateCardsForLevel(level),
	}

	gameUI := NewGameUI(window, game)
	gameUI.Render(level)
}

func showScores(window fyne.Window) {
	Players, _, _ := internal.Writing()

	for i := 0; i < len(Players); i++ {
		PlNames[i] = Players[i].Name
		PlScores[i], _ = strconv.Atoi(Players[i].Text)
	}
	scores := []shared.Score{
		{PlayerName: PlNames[0], Score: PlScores[0]},
		{PlayerName: PlNames[1], Score: PlScores[1]},
		{PlayerName: PlNames[2], Score: PlScores[2]},
		{PlayerName: PlNames[3], Score: PlScores[3]},
		{PlayerName: PlNames[4], Score: PlScores[4]},
		{PlayerName: PlNames[5], Score: PlScores[5]},
		{PlayerName: PlNames[6], Score: PlScores[6]},
		{PlayerName: PlNames[7], Score: PlScores[7]},
		{PlayerName: PlNames[8], Score: PlScores[8]},
		{PlayerName: PlNames[9], Score: PlScores[9]},
	}

	list := widget.NewList(
		func() int { return len(scores) },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(fmt.Sprintf("%s: %d",
				scores[i].PlayerName,
				scores[i].Score))
		},
	)

	backButton := widget.NewButton("Назад", func() {
		window.SetContent(createMainMenu(window))
	})

	window.SetContent(container.NewBorder(
		widget.NewLabel("Топ игроков"),
		backButton,
		nil,
		nil,
		list,
	))
}

func showAchievements(window fyne.Window) {
	// Здесь будет реализация отображения достижений
	achievements := []shared.Achievement{
		{Name: "Новичок", Description: "Завершить лёгкий уровень", Unlocked: achivments[0]},
		{Name: "Знаток", Description: "Пройти средний уровень за 3 минуты", Unlocked: achivments[1]},
		{Name: "Мастер памяти", Description: "Пройти сложный уровень без ошибок.", Unlocked: achivments[2]},
		{Name: "Легенда", Description: "Победить на экспертном уровне.", Unlocked: achivments[3]},
	}

	list := widget.NewList(
		func() int { return len(achievements) },
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel(""),
				widget.NewLabel(""),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			box := o.(*fyne.Container)
			title := box.Objects[0].(*widget.Label)
			desc := box.Objects[1].(*widget.Label)

			title.SetText(achievements[i].Name)
			if achievements[i].Unlocked {
				title.SetText(achievements[i].Name + " ✓")
			}
			desc.SetText(achievements[i].Description)
		},
	)

	backButton := widget.NewButton("Назад", func() {
		window.SetContent(createMainMenu(window))
	})

	window.SetContent(container.NewBorder(
		widget.NewLabel("Достижения"),
		backButton,
		nil,
		nil,
		list,
	))
}

func exit() {
	if MaxScore != 0 {
		internal.Writer("Player", fmt.Sprint(MaxScore))
	}
	myApp.Quit()
}

// Вспомогательные функции
func getRowsForLevel(level shared.GameLevel) int {
	switch level {
	case shared.Easy:
		return 4
	case shared.Medium:
		return 6
	case shared.Hard:
		return 8
	case shared.Expert:
		return 10
	}
	return 4
}

func getColumnsForLevel(level shared.GameLevel) int {
	return getRowsForLevel(level)
}

func generateCardsForLevel(level shared.GameLevel) []shared.Card {
	pairs := getRowsForLevel(level) * getColumnsForLevel(level) / 2
	var cards []shared.Card

	// Генерируем пары цифр от 1 до количества пар
	for i := 0; i < pairs; i++ {
		value := fmt.Sprintf("%d", i+1) // Цифры начиная с 1
		cards = append(cards, shared.Card{ID: i * 2, Value: value})
		cards = append(cards, shared.Card{ID: i*2 + 1, Value: value})
	}

	// Перемешиваем карточки в случайном порядке
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	return cards
}

func getLevelName(level shared.GameLevel) string {
	switch level {
	case shared.Easy:
		return "Лёгкий"
	case shared.Medium:
		return "Средний"
	case shared.Hard:
		return "Сложный"
	case shared.Expert:
		return "Эксперт"
	}
	return ""
}
