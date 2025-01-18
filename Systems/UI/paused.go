package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type PausedScreen struct {
	ScreenState map[string]any
}

// Draw implements Screens.
func (p *PausedScreen) Draw() {
	rl.DrawText(
		"PAUSED",
		int32(rl.GetScreenWidth()/2)-rl.MeasureText("PAUSED", 40)/2,
		int32(rl.GetScreenHeight()/2),
		40,
		rl.Gray,
	)
}

// GetScreenState implements Screens.
func (p *PausedScreen) GetScreenState() map[string]any {
	return p.ScreenState
}

// Update implements Screens.
func (p *PausedScreen) Update(state map[string]any) {
	for key, val := range state {
		p.ScreenState[key] = val
	}
}

