package systems

import (
	entities "github.com/prestonchoate/space-shmup/Entities"
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

func CreateUIManager(p *entities.Player, em *entities.EnemyManager) *UIManager {
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
			}
			break
		case Playing:
			break
		default:
			break
		}
	}

}
