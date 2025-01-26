package systems

import (
	ui "github.com/prestonchoate/space-shmup/Systems/UI"
)

type UIManager struct {
	TransitionReady bool
	TransitionState GameState
	screenList      map[GameState]ui.Screens
}

type UIUpdate struct {
	health     int
	score      int
	enemyCount int
	state      GameState
}

func CreateUIManager() *UIManager {
	screens := make(map[GameState]ui.Screens)

	screens[Start] = &ui.MainMenuScreen{
		ScreenState: make(map[string]any),
	}

	screens[Playing] = &ui.PlayingScreen{
		ScreenState: make(map[string]any),
	}

	screens[Paused] = &ui.PausedScreen{
		ScreenState: make(map[string]any),
	}

	screens[GameOver] = &ui.GameOverScreen{
		ScreenState: make(map[string]any),
	}

	return &UIManager{
		screenList: screens,
	}
}

func (u *UIManager) HandleGameStateRender(state GameState) {
	screen, exists := u.screenList[state]
	if exists {
		screen.Draw()
	}
}

func (u *UIManager) Update(update UIUpdate) {
	screenUpdate := map[string]any{
		"health":     update.health,
		"score":      update.score,
		"enemyCount": update.enemyCount,
	}

	screen, exists := u.screenList[update.state]
	if exists {
		screen.Update(screenUpdate)
		screenState := screen.GetScreenState()
		switch update.state {
		case Start:
			startButtonPressed, exists := screenState["startButtonPressed"].(bool)
			if exists && startButtonPressed {
				u.TransitionReady = true
				u.TransitionState = Playing
				screenState["startButtonPressed"] = false
				screen.Update(screenState)
			}
			break
		case Playing:
			break
		case GameOver:
			restartButtonPressed, exists := screenState["restartButtonPressed"].(bool)
			if exists && restartButtonPressed {
				u.TransitionReady = true
				u.TransitionState = Restart
				screenState["restartButtonPressed"] = false
				screen.Update(screenState)
			}
			break
		default:
			break
		}
	}

}
