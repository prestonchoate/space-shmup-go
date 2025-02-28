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
		score = 0
	}

	restartButtonText := "RESTART"
	exitButtonText := "QUIT"

	buttonWidth := rl.MeasureText(restartButtonText, 40)
	buttonHeight := 80
	screenWidth := rl.GetScreenWidth()
	screenHeight := rl.GetScreenHeight()

	spacing := int32(20)
	totalWidth := (2 * buttonWidth) + spacing
	startX := (screenWidth - int(totalWidth)) / 2
	buttonY := (screenHeight / 4) * 3

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

	g.ScreenState["restartButtonPressed"] = raygui.Button(rl.Rectangle{
		X:      float32(startX),
		Y:      float32(buttonY),
		Width:  float32(buttonWidth),
		Height: float32(buttonHeight),
	},
		restartButtonText,
	)
	g.ScreenState["exitButtonPressed"] = raygui.Button(rl.Rectangle{
		X:      float32(startX + int(spacing) + int(buttonWidth)),
		Y:      float32(buttonY),
		Width:  float32(buttonWidth),
		Height: float32(buttonHeight),
	},
		exitButtonText,
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
