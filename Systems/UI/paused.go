package ui

import (
	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type PausedScreen struct {
	ScreenState map[string]any
}

// Draw implements Screens.
func (p *PausedScreen) Draw() {
	exitButtonText := "QUIT"
	buttonWidth := rl.MeasureText(exitButtonText, 40)
	buttonHeight := 80

	rl.DrawText(
		"PAUSED",
		int32(rl.GetScreenWidth()/2)-rl.MeasureText("PAUSED", 40)/2,
		int32(rl.GetScreenHeight()/2),
		40,
		rl.Gray,
	)

	p.ScreenState["exitButtonPressed"] = raygui.Button(rl.Rectangle{
		X:      float32(rl.GetScreenWidth()-int(buttonWidth)) / 2,
		Y:      float32((rl.GetScreenHeight() / 4) * 3),
		Width:  float32(buttonWidth),
		Height: float32(buttonHeight),
	}, exitButtonText)
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
