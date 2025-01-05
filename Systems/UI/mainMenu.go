package ui

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type MainMenuScreen struct {
	ScreenState map[string]any
}

func (m *MainMenuScreen) Update(state *map[string]any) {

}

func (m *MainMenuScreen) Draw() {
	mainText := fmt.Sprintf("%v", "SPACE SHOOTER")
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

	m.ScreenState["startButtonPressed"] = raygui.Button(rl.Rectangle{
		X:      float32(startButtonX),
		Y:      float32(startButtonY),
		Width:  float32(rl.MeasureText(startButtonText, 40)),
		Height: float32(rl.MeasureText(startButtonText, 40)),
	},
		startButtonText,
	)
}

func (m *MainMenuScreen) GetScreenState() *map[string]any {
	return &m.ScreenState
}
