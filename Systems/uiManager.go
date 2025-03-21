package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	assets "github.com/prestonchoate/space-shmup/Systems/Assets"
	systems_data "github.com/prestonchoate/space-shmup/Systems/Data"
	events "github.com/prestonchoate/space-shmup/Systems/Events"
	events_data "github.com/prestonchoate/space-shmup/Systems/Events/Data"
	ui "github.com/prestonchoate/space-shmup/Systems/UI"
	"github.com/prestonchoate/space-shmup/Systems/saveManager"
)

type UIManager struct {
	screenList map[systems_data.GameState]ui.Screens
}

type UIUpdate struct {
	health     int
	score      int
	enemyCount int
	state      systems_data.GameState
}

func CreateUIManager() *UIManager {
	screens := make(map[systems_data.GameState]ui.Screens)

	screens[systems_data.Start] = &ui.MainMenuScreen{
		ScreenState: make(map[string]any),
	}

	screens[systems_data.Playing] = &ui.PlayingScreen{
		ScreenState: make(map[string]any),
	}

	screens[systems_data.Paused] = &ui.PausedScreen{
		ScreenState: make(map[string]any),
	}

	screens[systems_data.GameOver] = &ui.GameOverScreen{
		ScreenState: make(map[string]any),
	}

	screens[systems_data.Settings] = &ui.SettingsScreen{
		ScreenState: make(map[string]any),
	}

	return &UIManager{
		screenList: screens,
	}
}

func (u *UIManager) HandleGameStateRender(state systems_data.GameState) {
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
		case systems_data.Start:
			startButtonPressed, exists := screenState["startButtonPressed"].(bool)
			if exists && startButtonPressed {
				screenState["startButtonPressed"] = false
				screen.Update(screenState)
				events.GetEventManagerInstance().Emit(events_data.ChangeGameState, events_data.ChangeStateData{
					NewState: systems_data.Playing,
				})
				u.playConfirmSound()
				break
			}
			exitButtonPressed, exists := screenState["exitButtonPressed"].(bool)
			if exists && exitButtonPressed {
				events.GetEventManagerInstance().Emit(events_data.ChangeGameState, events_data.ChangeStateData{
					NewState: systems_data.Exit,
				})
				u.playConfirmSound()
				break
			}
			settingsButtonPressed, exists := screenState["settingsButtonPressed"].(bool)
			if exists && settingsButtonPressed {
				events.GetEventManagerInstance().Emit(events_data.ChangeGameState, events_data.ChangeStateData{
					NewState: systems_data.Settings,
				})
				u.playConfirmSound()
				break
			}
			break
		case systems_data.Playing:
			break
		case systems_data.Paused:
			exitButtonPressed, exists := screenState["exitButtonPressed"].(bool)
			if exists && exitButtonPressed {
				events.GetEventManagerInstance().Emit(events_data.ChangeGameState, events_data.ChangeStateData{
					NewState: systems_data.Exit,
				})
				u.playConfirmSound()
				break
			}
			settingsButtonPressed, exists := screenState["settingsButtonPressed"].(bool)
			if exists && settingsButtonPressed {
				screenState["settingsButtonPressed"] = false
				events.GetEventManagerInstance().Emit(events_data.ChangeGameState, events_data.ChangeStateData{
					NewState: systems_data.Settings,
				})
				u.playConfirmSound()
				break
			}
			break
		case systems_data.GameOver:
			restartButtonPressed, exists := screenState["restartButtonPressed"].(bool)
			if exists && restartButtonPressed {
				screenState["restartButtonPressed"] = false
				screen.Update(screenState)
				events.GetEventManagerInstance().Emit(events_data.ChangeGameState, events_data.ChangeStateData{
					NewState: systems_data.Restart,
				})
				u.playConfirmSound()
				break
			}
			exitButtonPressed, exists := screenState["exitButtonPressed"].(bool)
			if exists && exitButtonPressed {
				events.GetEventManagerInstance().Emit(events_data.ChangeGameState, events_data.ChangeStateData{
					NewState: systems_data.Exit,
				})
				u.playConfirmSound()
				break
			}
			break
		case systems_data.Settings:
			backButtonPressed, exists := screenState["back"].(bool)
			if exists && backButtonPressed {
				screenState["back"] = false
				events.GetEventManagerInstance().Emit(events_data.ReturnGameState, events_data.ReturnStateData{})
				u.playConfirmSound()
				break
			}
			saveButtonPressed, exists := screenState["save"].(bool)
			if exists && saveButtonPressed {
				screenState["save"] = false
				settings, exists := screenState["settings"].(*systems_data.GameSettings)
				if exists {
					saveManager.GetInstance().UpdateSettings(settings)
				}
				events.GetEventManagerInstance().Emit(events_data.ReturnGameState, events_data.ReturnStateData{})
				u.playConfirmSound()
				break
			}
			break
		default:
			break
		}
	}

}

func (u *UIManager) playConfirmSound() {
	sound, ok := assets.GetAssetManagerInstance().GetSound("assets/sfx/Interface_Bleeps_Wav/Confirm_02.wav")
	if ok {
		rl.PlaySound(sound)
		rl.SetSoundVolume(sound, saveManager.GetInstance().Data.Settings.SfxVolume)
	}
}
