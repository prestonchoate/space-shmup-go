package ui

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type MainMenuScreen struct {
	ScreenState map[string]any
}

func (m *MainMenuScreen) Update(state map[string]any) {

}

func (m *MainMenuScreen) Draw() {

	startButtonText := "START"
	settingsButtonText := "SETTINGS"
	exitButtonText := "QUIT"

	buttonWidth := rl.MeasureText(settingsButtonText, 40)
	buttonHeight := 80
	screenWidth := rl.GetScreenWidth()
	screenHeight := rl.GetScreenHeight()

	spacing := int32(20)
	totalWidth := (3 * buttonWidth) + (2 * spacing)
	startX := (screenWidth - int(totalWidth)) / 2
	buttonY := (screenHeight / 4) * 3

	mainText := fmt.Sprintf("%v", "UNTITLED SPACE SHOOTER")
	rl.DrawText(
		mainText,
		int32(rl.GetScreenWidth()/2)-rl.MeasureText(mainText, 60)/2,
		int32(rl.GetScreenHeight()/2),
		60,
		rl.White,
	)

	m.ScreenState["startButtonPressed"] = raygui.Button(rl.Rectangle{
		X:      float32(startX),
		Y:      float32(buttonY),
		Width:  float32(buttonWidth),
		Height: float32(buttonHeight),
	},
		startButtonText,
	)

	m.ScreenState["settingsButtonPressed"] = raygui.Button(rl.Rectangle{
		X:      float32(startX + int(spacing) + int(buttonWidth)),
		Y:      float32(buttonY),
		Width:  float32(buttonWidth),
		Height: float32(buttonHeight),
	},
		settingsButtonText,
	)

	m.ScreenState["exitButtonPressed"] = raygui.Button(rl.Rectangle{
		X:      float32(startX + 2*(int(buttonWidth)+int(spacing))),
		Y:      float32(buttonY),
		Width:  float32(buttonWidth),
		Height: float32(buttonHeight),
	},
		exitButtonText,
	)
}

func (m *MainMenuScreen) GetScreenState() map[string]any {
	return m.ScreenState
}
