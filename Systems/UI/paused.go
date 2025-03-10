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
	spacing := float32(10)
	buttonWidth := float32(rl.GetScreenWidth() / 15)
	startX := ((float32(rl.GetScreenWidth())) / 2) - buttonWidth - spacing
	exitButtonText := "QUIT"
	settingsButtonText := "SETTINGS"
	buttonHeight := 80

	rl.DrawText(
		"PAUSED",
		int32(rl.GetScreenWidth()/2)-rl.MeasureText("PAUSED", 40)/2,
		int32(rl.GetScreenHeight()/2),
		40,
		rl.Gray,
	)

	p.ScreenState["settingsButtonPressed"] = raygui.Button(rl.Rectangle{
		X:      startX,
		Y:      float32((rl.GetScreenHeight() / 4) * 3),
		Width:  float32(buttonWidth),
		Height: float32(buttonHeight),
	}, settingsButtonText)

	p.ScreenState["exitButtonPressed"] = raygui.Button(rl.Rectangle{
		X:      startX + buttonWidth + float32(spacing),
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
