package systems

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	entities "github.com/prestonchoate/space-shmup/Entities"
)

type UIManager struct {
	player          *entities.Player
	enemyManager    *entities.EnemyManager
	health          int
	score           int
	enemyCount      int
	transitionReady bool
	transitionState GameState
}

func CreateUIManager(p *entities.Player, em *entities.EnemyManager) *UIManager {
	return &UIManager{
		player:       p,
		enemyManager: em,
	}
}

func (u *UIManager) HandleGameStateRender(state GameState) {
	switch state {
	case Start:
		u.drawStartState()
	case Playing:
		u.drawPlayState()
	case GameOver:
		u.drawGameOverState()
	case Paused:
		u.drawPausedState()
	}
}

func (u *UIManager) Update(state GameState) {
	if state == Playing {
		u.updatePlayState()
	}
	if state == Start {
		u.updateStartState()
	}
}

func (u *UIManager) updatePlayState() {
	u.health = u.player.GetHealth()
	u.score = u.player.GetScore()
	u.enemyCount = u.enemyManager.GetEnemyCount()
}

func (u *UIManager) updateStartState() {
	// Check if mouse click happens on start button

	// Check if mouse click happens on settings button
}

func (u *UIManager) drawStartState() {
	mainText := fmt.Sprintf("%v", "SPACE SHOOTER")
	subText := fmt.Sprintf("%v", "Press ENTER to start")
	rl.DrawText(
		mainText,
		int32(rl.GetScreenWidth()/2)-rl.MeasureText(mainText, 60)/2,
		int32(rl.GetScreenHeight()/2),
		60,
		rl.White,
	)

	// Draw Gray Rectangle with START text
	startButtonText := "START"
	startButtonX := rl.GetScreenWidth() / 3
	startButtonY := (rl.GetScreenHeight() / 4) * 3
	rl.DrawRectangle(
		int32(startButtonX),
		int32(startButtonY),
		rl.MeasureText(startButtonText, 40)+40,
		rl.MeasureText(startButtonText, 40)+40,
		rl.NewColor(0, 66, 37, 255),
	)

	rl.DrawText(
		subText,
		int32(startButtonX+20),
		int32(startButtonY+20),
		20,
		rl.White,
	)

}

func (u *UIManager) drawPausedState() {
	rl.DrawText(
		"PAUSED",
		int32(rl.GetScreenWidth()/2)-rl.MeasureText("PAUSED", 40)/2,
		int32(rl.GetScreenHeight()/2),
		40,
		rl.Gray,
	)
}

func (u *UIManager) drawGameOverState() {
	rl.DrawText(
		"GAME OVER",
		int32(rl.GetScreenWidth()/2)-rl.MeasureText("GAME OVER", 40)/2,
		int32(rl.GetScreenHeight()/2),
		40,
		rl.Gray,
	)

	finalScore := fmt.Sprintf("%10v%06d", "Final Score: ", u.score)
	rl.DrawText(
		finalScore,
		int32(rl.GetScreenWidth()/2)-rl.MeasureText(finalScore, 20)/2,
		int32(rl.GetScreenHeight()/2)+50,
		20,
		rl.Gray,
	)
	rl.DrawText(
		"PRESS ENTER TO RESTART",
		int32(rl.GetScreenWidth()/2)-rl.MeasureText("PRESS ENTER TO RESTART", 20)/2,
		int32(rl.GetScreenHeight()/2)+80,
		20,
		rl.Gray,
	)
}

func (u *UIManager) drawPlayState() {

	rl.DrawText(
		fmt.Sprintf("%10v%03d", "Health: ", u.health),
		int32(rl.GetScreenWidth()-200),
		0,
		20,
		rl.Gray)

	rl.DrawText(
		fmt.Sprintf("%10v%06d", "Score: ", u.score),
		int32(rl.GetScreenWidth()-200),
		30,
		20,
		rl.Gray)

	rl.DrawText(
		fmt.Sprintf("%10v%06d", "Enemies: ", u.enemyCount),
		int32(rl.GetScreenWidth()-200),
		60,
		20,
		rl.Gray)
}
