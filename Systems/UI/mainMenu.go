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
	m.ScreenState["screenSize"] = rl.Vector2{X: float32(rl.GetScreenWidth()), Y: float32(rl.GetScreenHeight())}
}

func (m *MainMenuScreen) Draw() {

	startButtonText := "START"
	settingsButtonText := "SETTINGS"
	exitButtonText := "QUIT"

	buttonWidth := rl.MeasureText(settingsButtonText, 40)
	buttonHeight := 80
	screenSize, ok := m.ScreenState["screnSize"].(rl.Vector2)
	var screenWidth int
	var screenHeight int
	if !ok {
		screenWidth = rl.GetScreenWidth()
		screenHeight = rl.GetScreenHeight()
	} else {
		screenWidth = int(screenSize.X)
		screenHeight = int(screenSize.Y)
	}

	spacing := int32(20)
	totalWidth := (3 * buttonWidth) + (2 * spacing)
	startX := (screenWidth - int(totalWidth)) / 2
	buttonY := (screenHeight / 4) * 3

	mainText := fmt.Sprintf("%v", "UNTITLED SPACE SHOOTER")
	rl.DrawText(
		mainText,
		int32(screenWidth/2)-rl.MeasureText(mainText, 60)/2,
		int32(screenHeight/2),
		60,
		rl.NewColor(81, 191, 211, 255),
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
