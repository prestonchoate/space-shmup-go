package systems

import (
	"fmt"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
	assets "github.com/prestonchoate/space-shmup/Systems/Assets"
	systems_data "github.com/prestonchoate/space-shmup/Systems/Data"
	events "github.com/prestonchoate/space-shmup/Systems/Events"
	events_data "github.com/prestonchoate/space-shmup/Systems/Events/Data"
	ui "github.com/prestonchoate/space-shmup/Systems/UI"
	"github.com/prestonchoate/space-shmup/Systems/saveManager"
)

// TODO: possibility for concurrent read and write to messages map. Might need to add a mutex handler
type UIManager struct {
	screenList         map[systems_data.GameState]ui.Screens
	messages           []*events_data.AddMessageData
	showWaveComplete   bool
	waveCompleteTicker float32
}

type UIUpdate struct {
	health     float32
	score      int
	enemyCount int
	state      systems_data.GameState
	delta      float32
}

func CreateUIManager() *UIManager {

	u := &UIManager{
		messages: make([]*events_data.AddMessageData, 10),
	}

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

	screens[systems_data.Shop] = &ui.ShopScreen{
		ScreenState: make(map[string]any),
	}

	u.screenList = screens
	events.GetEventManagerInstance().Subscribe(events_data.AddMessage, u.handleAddMessageEvent)
	events.GetEventManagerInstance().Subscribe(events_data.SubmitHighScore, u.handleScoreSubmissionEvent)
	events.GetEventManagerInstance().Subscribe(events_data.WaveCompleteEvent, u.handleWaveCompleteEvent)
	return u
}

func (u *UIManager) HandleGameStateRender(state systems_data.GameState) {
	screen, exists := u.screenList[state]
	if exists {
		screen.Draw()
	}

	msgX := rl.GetScreenWidth()
	msgY := rl.GetScreenHeight() - 30

	if rl.IsWindowFullscreen() {
		msgX = rl.GetMonitorWidth(rl.GetCurrentMonitor())
		msgY = rl.GetMonitorHeight(rl.GetCurrentMonitor()) - 30
	}
	padding := 20

	for _, msg := range u.messages {
		if msg == nil {
			continue
		}
		if msg.Timer >= 0 {
			rl.DrawText(msg.Message, int32(msgX)-rl.MeasureText(msg.Message, 20)-10, int32(msgY), 20, msg.Color)
			msgY -= padding
		}
	}

	if u.showWaveComplete && state == systems_data.Playing {
		u.drawWaveCompleteText()
	}
}

func (u *UIManager) drawWaveCompleteText() {
	var w, h int
	if rl.IsWindowFullscreen() {
		w = rl.GetMonitorWidth(rl.GetCurrentMonitor())
		h = rl.GetMonitorHeight(rl.GetCurrentMonitor())
	} else {
		w = rl.GetScreenWidth()
		h = rl.GetScreenHeight()
	}

	text := "WAVE COMPLETE"
	fontSize := int32(60)
	textWidth := rl.MeasureText(text, fontSize)
	x := w/2 - int(textWidth)
	y := h / 2

	rl.DrawText(text, int32(x), int32(y), fontSize, rl.NewColor(81, 191, 211, 255))
	rl.DrawText(fmt.Sprintf("%.02f", u.waveCompleteTicker), int32(x), int32(y)+10, 20, rl.White)
}

func (u *UIManager) Update(update UIUpdate) {
	screenUpdate := map[string]any{
		"health":     update.health,
		"score":      update.score,
		"enemyCount": update.enemyCount,
	}

	u.updateMessageTimers(update.delta)

	if update.state != systems_data.GameOver {
		u.resetGameOverScreen()
	}

	screen, exists := u.screenList[update.state]
	if !exists {
		return
	}

	screen.Update(screenUpdate)
	screenState := screen.GetScreenState()

	// Handle state-specific logic
	switch update.state {
	case systems_data.Start:
		u.handleStartScreen(screen, screenState)
	case systems_data.Playing:
		u.handlePlayingScreen(screen, screenState, update)
	case systems_data.Paused:
		u.handlePausedScreen(screen, screenState)
	case systems_data.GameOver:
		u.handleGameOverScreen(screen, screenState, update)
	case systems_data.Settings:
		u.handleSettingsScreen(screen, screenState)
	case systems_data.Shop:
		u.handleShopScreen(screen, screenState)
	}

	u.cleanupMessageQueue()

	if u.showWaveComplete && update.state == systems_data.Playing {
		u.waveCompleteTicker -= update.delta
		if u.waveCompleteTicker <= 0 {
			u.showWaveComplete = false
			u.changeGameState(systems_data.Shop)
		}
	}
}

// Helper function to update message timers
func (u *UIManager) updateMessageTimers(delta float32) {
	for idx, msg := range u.messages {
		if msg == nil {
			continue
		}
		msg.Timer = msg.Timer - delta
		if msg.Timer <= 0 {
			u.messages[idx] = nil
			continue
		}
		u.messages[idx] = msg
	}
}

// Helper function to reset game over screen
func (u *UIManager) resetGameOverScreen() {
	s, exists := u.screenList[systems_data.GameOver]
	if exists {
		state := s.GetScreenState()
		delete(state, "scoreSubmitted")
		s.Update(state)
	}
}

// Handle start screen state
func (u *UIManager) handleStartScreen(screen ui.Screens, screenState map[string]any) {
	if startButtonPressed, exists := screenState["startButtonPressed"].(bool); exists && startButtonPressed {
		screenState["startButtonPressed"] = false
		screen.Update(screenState)
		u.changeGameState(systems_data.Playing)
		return
	}

	if exitButtonPressed, exists := screenState["exitButtonPressed"].(bool); exists && exitButtonPressed {
		u.changeGameState(systems_data.Exit)
		return
	}

	if settingsButtonPressed, exists := screenState["settingsButtonPressed"].(bool); exists && settingsButtonPressed {
		u.changeGameState(systems_data.Settings)
		return
	}
}

// Handle playing screen state
func (u *UIManager) handlePlayingScreen(screen ui.Screens, screenState map[string]any, update UIUpdate) {
	// Currently empty as there's no logic in the original function
}

// Handle paused screen state
func (u *UIManager) handlePausedScreen(screen ui.Screens, screenState map[string]any) {
	if exitButtonPressed, exists := screenState["exitButtonPressed"].(bool); exists && exitButtonPressed {
		u.changeGameState(systems_data.Exit)
		return
	}

	if settingsButtonPressed, exists := screenState["settingsButtonPressed"].(bool); exists && settingsButtonPressed {
		screenState["settingsButtonPressed"] = false
		screen.Update(screenState)
		u.changeGameState(systems_data.Settings)
		return
	}
}

// Handle game over screen state
func (u *UIManager) handleGameOverScreen(screen ui.Screens, screenState map[string]any, update UIUpdate) {
	if restartButtonPressed, exists := screenState["restartButtonPressed"].(bool); exists && restartButtonPressed {
		screenState["restartButtonPressed"] = false
		screen.Update(screenState)
		u.changeGameState(systems_data.Restart)
		return
	}

	if exitButtonPressed, exists := screenState["exitButtonPressed"].(bool); exists && exitButtonPressed {
		u.changeGameState(systems_data.Exit)
		return
	}

	if submitButtonPressed, exists := screenState["submitButtonPressed"].(bool); exists && submitButtonPressed {
		initials, exists := screenState["playerInitials"].(string)
		if exists && len(initials) > 0 {
			events.GetEventManagerInstance().Emit(events_data.SubmitHighScore, events_data.HighScoreData{
				Initials: initials,
				Score:    int64(update.score),
			})
			u.playConfirmSound()
			screenState["scoreSubmitted"] = true
			screenState["submitButtonPressed"] = false
			screen.Update(screenState)
		}
	}
}

// Handle settings screen state
func (u *UIManager) handleSettingsScreen(screen ui.Screens, screenState map[string]any) {
	if backButtonPressed, exists := screenState["back"].(bool); exists && backButtonPressed {
		screenState["back"] = false
		screen.Update(screenState)
		events.GetEventManagerInstance().Emit(events_data.ReturnGameState, events_data.ReturnStateData{})
		u.playConfirmSound()
		return
	}

	if saveButtonPressed, exists := screenState["save"].(bool); exists && saveButtonPressed {
		screenState["save"] = false
		screen.Update(screenState)
		if settings, exists := screenState["settings"].(*systems_data.GameSettings); exists {
			saveManager.GetInstance().UpdateSettings(settings)
		}
		events.GetEventManagerInstance().Emit(events_data.ReturnGameState, events_data.ReturnStateData{})
		u.playConfirmSound()
		return
	}
}

// Handle shop screen state
func (u *UIManager) handleShopScreen(screen ui.Screens, screenState map[string]any) {
	spawned, exists := screenState["spawnedUpgrades"].(bool)
	if !exists || !spawned {
		// TODO: get number of upgrades from somewhere?
		upgrades := GetUpgrader().GetUpgrades(3)
		screenState["upgrades"] = upgrades
		screenState["spawnedUpgrades"] = true
		screen.Update(screenState)
	}

	rerollButtonPressed, exists := screenState["rerollButtonPressed"].(bool)
	if exists && rerollButtonPressed {
		upgrades := GetUpgrader().GetUpgrades(3)
		screenState["upgrades"] = upgrades
		screenState["spawnedUpgrades"] = true
		screen.Update(screenState)
	}

	upgrade, exists := screenState["selectedUpgrade"].(*systems_data.StatUpgrade)
	if exists && upgrade != nil {
		events.GetEventManagerInstance().Emit(events_data.StatUpgradeEvent, events_data.StatUpgradeData{
			Upgrade: upgrade,
		})
		u.playConfirmSound()
		screenState["selectedUpgrade"] = nil
		screenState["spawnedUpgrades"] = false
		screen.Update(screenState)
		u.changeGameState(systems_data.Playing)
	}

	skip, exists := screenState["skipButtonPressed"].(bool)
	if exists && skip {
		u.playConfirmSound()
		screenState["skipButtonPressed"] = false
		screen.Update(screenState)
		u.changeGameState(systems_data.Playing)
	}
}

// Helper function to change game state
func (u *UIManager) changeGameState(newState systems_data.GameState) {
	events.GetEventManagerInstance().Emit(events_data.ChangeGameState, events_data.ChangeStateData{
		NewState: newState,
	})
	u.playConfirmSound()
}

func (u *UIManager) playConfirmSound() {
	sound, ok := assets.GetAssetManagerInstance().GetSound("assets/sfx/Interface_Bleeps_Wav/Confirm_02.wav")
	if ok {
		rl.PlaySound(sound)
		rl.SetSoundVolume(sound, saveManager.GetInstance().Data.Settings.SfxVolume)
	}
}

func (u *UIManager) cleanupMessageQueue() {
	newMsgs := make([]*events_data.AddMessageData, len(u.messages))

	for _, msg := range u.messages {
		if msg != nil {
			newMsgs = append(newMsgs, msg)
		}
	}

	u.messages = newMsgs
}

func (u *UIManager) handleAddMessageEvent(event events.Event) {
	if data, ok := event.Data.(events_data.AddMessageData); ok {
		u.messages = append(u.messages, &data)
	}
}

func (u *UIManager) handleScoreSubmissionEvent(event events.Event) {
	if data, ok := event.Data.(events_data.ScoreSubmissionCompleteData); ok {
		s, exists := u.screenList[systems_data.GameOver]
		if exists && !data.Success {
			log.Println("Removing score submission flag from game over state")
			state := s.GetScreenState()
			delete(state, "scoreSubmitted")
			s.Update(state)
		}

	}
}

func (u *UIManager) handleWaveCompleteEvent(e events.Event) {
	if _, ok := e.Data.(events_data.WaveCompleteData); ok {
		u.waveCompleteTicker = 2.0
		u.showWaveComplete = true
	}
}
