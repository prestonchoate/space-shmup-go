package ui

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameOverScreen struct {
	ScreenState map[string]any
}

// Draw implements Screens.
func (g *GameOverScreen) Draw() {
	rl.DrawText(
		"GAME OVER",
		int32(rl.GetScreenWidth()/2)-rl.MeasureText("GAME OVER", 40)/2,
		int32(rl.GetScreenHeight()/2),
		40,
		rl.Gray,
	)

	finalScore := fmt.Sprintf("%10v%06d", "Final Score: ", 0)
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

