package game

import (
	"etic"
)

const Title = "Собери все пары!"

type ScreenState int

const (
	StateScreenSelect ScreenState = iota
	StateScreenPlaying
	StateScreenScore
)

func (s ScreenState) String() string { return []string{"StateSelect", "StatePlaying", "StateScore"}[s] }

func NewMainScreen(id etic.ElementID) *etic.Element {
	var (
		mainScreenState *etic.Property[ScreenState]
		screenSelect    *ScreenSelect
		screenGame      *ScreenGame
		screenScore     *ScreenScore
		gamesData       *GamesData = NewGamesData()
	)
	mainScene := etic.NewElement(id)

	mainSceneLogic := func(b *etic.Button) {
		switch mainScreenState.Get() {
		case StateScreenSelect:
			gamesData.SetCurrentDim(ParseDim(b.Text.Get()))
			screenGame.New(NewGameData(gamesData.CurrentDim(), gamesData.ID()))
			mainScreenState.Set(StateScreenPlaying)
		case StateScreenPlaying:
		case StateScreenScore:
			switch b.Text.Get() {
			case btnResetGame:
				screenGame.New(NewGameData(gamesData.CurrentDim(), gamesData.ID()))
				mainScreenState.Set(StateScreenPlaying)
			case btnNewGame:
				gamesData.SetCurrentDim(gamesData.Next())
				screenGame.New(NewGameData(gamesData.CurrentDim(), gamesData.ID()))
				mainScreenState.Set(StateScreenPlaying)
			case btnGoMenu:
				mainScreenState.Set(StateScreenSelect)
			}
		}
	}

	gameLogic := func(b *etic.Button) {
		game := screenGame.game
		for i, v := range screenGame.board.Children() {
			btn := v.(*etic.Button)
			if btn == b {
				switch game.state.Get() {
				case GameStart:
					game.state.Set(GamePlay)
					game.data.Start()
					game.Move(i)
				case GamePlay:
					game.Move(i)
				case GamePause:
				case GameWin:
				}
			}
		}
	}

	screenSelect = NewScreenSelect("screen.select", *gamesData, mainSceneLogic)
	screenGame = NewScreenGame("screen.game", gameLogic)
	screenScore = NewScreenScore("screen.score", mainSceneLogic)

	mainScreenState = etic.NewProperty(StateScreenSelect)
	mainScreenState.OnChange.Connect(func(gs ScreenState) {
		switch gs {
		case StateScreenSelect:
			screenSelect.Show()
			screenGame.Hide()
			screenScore.Hide()
		case StateScreenPlaying:
			screenGame.Show()
			screenSelect.Hide()
			screenScore.Hide()
		case StateScreenScore:
			screenScore.Show()
			screenGame.Hide()
			screenSelect.Hide()
		}
	})

	screenGame.game.state.OnChange.Connect(func(gs GameState) {
		switch gs {
		case GameStart:
		case GamePlay:
		case GamePause:
		case GameWin:
			result := screenGame.game.Finish()
			gamesData.SaveGame(result)
			screenScore.SetLastGameScore(result)
			screenScore.SetTop10(gamesData.Top10())
			idx, _ := gamesData.NextUnlocked()
			screenSelect.levels.Children()[idx].(*etic.Button).State.Set(etic.ElementNormal)
			mainScreenState.Set(StateScreenScore)
		}
	})

	mainScene.Layout().Set(etic.NewStackLayout(0.01))
	mainScene.Add(screenSelect)
	mainScene.Add(screenGame)
	mainScene.Add(screenScore)
	return mainScene
}
