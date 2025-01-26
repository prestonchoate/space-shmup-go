package ui

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameOverScreen struct {
	ScreenState map[string]any
}

// Draw implements Screens.
func (g *GameOverScreen) Draw() {
	score, exists := g.ScreenState["score"]
	if !exists {
		score = 0;
	}

	rl.DrawText(
		"GAME OVER",
		int32(rl.GetScreenWidth()/2)-rl.MeasureText("GAME OVER", 40)/2,
		int32(rl.GetScreenHeight()/2),
		40,
		rl.Gray,
	)

	finalScore := fmt.Sprintf("%10v%06d", "Final Score: ", score)
	rl.DrawText(
		finalScore,
		int32(rl.GetScreenWidth()/2)-rl.MeasureText(finalScore, 20)/2,
		int32(rl.GetScreenHeight()/2)+50,
		20,
		rl.Gray,
	)

	restartButtonText := "RESTART"
	restartButtonX := rl.GetScreenWidth() / 3
	restartButtonY := (rl.GetScreenHeight() / 4) * 3

	g.ScreenState["restartButtonPressed"] = raygui.Button(
		rl.Rectangle{
			X: float32(restartButtonX),
			Y: float32(restartButtonY),
			Width: float32(rl.MeasureText(restartButtonText, 40)),
			Height: float32(rl.MeasureText(restartButtonText, 40)),
		}, 
		restartButtonText,
	)
}

// GetScreenState implements Screens.
func (g *GameOverScreen) GetScreenState() map[string]any {
	return g.ScreenState
}

// Update implements Screens.
func (g *GameOverScreen) Update(state map[string]any) {
	for key, val := range state {
		g.ScreenState[key] = val
	}
}

