package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	systems_data "github.com/prestonchoate/space-shmup/Systems/Data"
	"github.com/prestonchoate/space-shmup/Systems/saveManager"
)

type SettingsScreen struct {
	ScreenState        map[string]any
	settings           *systems_data.GameSettings
	editKey            *int32
	selectedResolution int32
	resolutions        []string
	resolutionList     string
	resolutionActive   bool
	editFps            bool
	fpsStr             string
}

func (s *SettingsScreen) Update(state map[string]any) {
	s.settings = &saveManager.GetInstance().Data.Settings
	if len(s.resolutions) == 0 {
		s.resolutions = []string{
			"3840x2160", "2560x1440", "1920x1080", "1600x900",
			"1366x768", "1280x720", "1024x768", "800x600",
			"5120x2160", "3840x1600", "3440x1440", "2560x1080",
			"7680x2160", "5120x1440", "3840x1080",
		}
		s.resolutionList = strings.Join(s.resolutions, ";")
		s.fpsStr = strconv.Itoa(int(s.settings.TargetFPS))
	}

	selectedIndex := 0
	currentWidth := s.settings.ScreenWidth
	currentHeight := s.settings.ScreenHeight
	if currentWidth == 0 {
		currentWidth = rl.GetMonitorWidth(rl.GetCurrentMonitor())
	}

	if currentHeight == 0 {
		currentHeight = rl.GetMonitorHeight(rl.GetCurrentMonitor())
	}

	current := fmt.Sprintf("%dx%d", currentWidth, currentHeight)
	for i, res := range s.resolutions {
		if res == current {
			selectedIndex = i
			break
		}
	}

	s.selectedResolution = int32(selectedIndex)
	s.ScreenState["settings"] = s.settings
	s.ScreenState["windowSize"] = rl.Vector2{X: float32(rl.GetRenderWidth()), Y: float32(rl.GetRenderHeight())}

}

func (s *SettingsScreen) Draw() {
	screenWidth := s.ScreenState["windowSize"].(rl.Vector2).X
	screenHeight := s.ScreenState["windowSize"].(rl.Vector2).Y
	vSpace := int32(20)
	startX := int32(screenWidth / 2)
	startY := int32(screenHeight / 3)
	buttonWidth := int32(screenWidth / 10)
	buttonHeight := int32(screenHeight / 30)
	labelOffset := int32(screenWidth / 15)
	startX -= (buttonWidth + labelOffset) / 2

	raygui.Label(rl.Rectangle{X: float32(startX + (buttonWidth+labelOffset)/2), Y: float32(screenHeight / 4), Width: float32(buttonWidth+labelOffset) * 2, Height: float32(buttonHeight) * 2}, "SETTINGS")

	// I don't love this but it prevents multiple button presses when the drop down is active
	if !s.resolutionActive {

		// FPS Setting - need to initialize with setting data once
		raygui.Label(rl.Rectangle{X: float32(startX), Y: float32(startY), Width: float32(labelOffset), Height: float32(buttonHeight)}, "Target FPS")
		if raygui.TextBox(rl.Rectangle{
			X:      float32(startX + labelOffset),
			Y:      float32(startY),
			Width:  float32(buttonWidth),
			Height: float32(buttonHeight),
		}, &s.fpsStr, 20, s.editFps) {
			s.editFps = !s.editFps
		}

		if !s.editFps {
			if val, err := strconv.Atoi(s.fpsStr); err == nil {
				s.settings.TargetFPS = int32(val)
			}
		}

		// Fullscreen Checkbox
		raygui.Label(rl.Rectangle{X: float32(startX), Y: float32((startY + (buttonHeight * 2) + vSpace) + buttonHeight - 10), Width: float32(labelOffset), Height: float32(buttonHeight)}, "Fullscreen:")
		s.settings.Fullscreen = raygui.CheckBox(
			rl.Rectangle{X: float32(startX + labelOffset), Y: float32((startY + (buttonHeight * 2) + vSpace) + buttonHeight - 10), Width: float32(20), Height: float32(20)},
			"", s.settings.Fullscreen,
		)

		// Keybindings
		keyStartY := startY + (buttonHeight * 3) + vSpace
		keyBindings := []struct {
			label string
			key   *int32
			y     int32
		}{
			{"Move Left:", &s.settings.Keys.KeyLeft, keyStartY + buttonHeight + vSpace},
			{"Move Right:", &s.settings.Keys.KeyRight, keyStartY + (buttonHeight * 2) + vSpace},
			{"Move Up:", &s.settings.Keys.KeyUp, keyStartY + (buttonHeight * 3) + vSpace},
			{"Move Down:", &s.settings.Keys.KeyDown, keyStartY + (buttonHeight * 4) + vSpace},
			{"Fire:", &s.settings.Keys.KeyFire, keyStartY + (buttonHeight * 5) + vSpace},
		}

		for _, bind := range keyBindings {
			raygui.Label(rl.Rectangle{X: float32(startX), Y: float32(bind.y), Width: float32(labelOffset), Height: float32(buttonHeight)}, bind.label)
			btnText := saveManager.KeyToString(*bind.key)
			if s.editKey == bind.key {
				btnText = "Press any key..."
				newKey := rl.GetKeyPressed()
				if newKey != 0 {
					*bind.key = newKey
					s.editKey = nil
				}
			}

			if raygui.Button(rl.Rectangle{X: float32(startX + labelOffset), Y: float32(bind.y), Width: float32(buttonWidth), Height: float32(buttonHeight / 2)}, btnText) {
				s.editKey = bind.key
			}
		}

		// Save & Back Buttons
		if raygui.Button(rl.Rectangle{X: float32((startX) + ((buttonWidth + labelOffset) / 2) - 160), Y: float32(startY + (buttonHeight * 10) + vSpace), Width: 140, Height: float32(buttonHeight)}, "Back") {
			s.ScreenState["back"] = true
		}
		if raygui.Button(rl.Rectangle{X: float32(startX + ((buttonWidth + labelOffset) / 2) + 160), Y: float32(startY + (buttonHeight * 10) + vSpace), Width: 140, Height: float32(buttonHeight)}, "Save") {
			s.ScreenState["save"] = true
		}
	}

	// Draw the dropdown last so it appears on top of other menu items
	// Screen Resolution Dropdown
	raygui.Label(rl.Rectangle{X: float32(startX), Y: float32(startY + buttonHeight + vSpace), Width: float32(labelOffset), Height: float32(buttonHeight)}, "Resolution:")
	if raygui.DropdownBox(
		rl.Rectangle{X: float32(startX + labelOffset), Y: float32(startY + buttonHeight + vSpace), Width: float32(buttonWidth), Height: float32(buttonHeight)},
		s.resolutionList, &s.selectedResolution, s.resolutionActive,
	) {
		s.resolutionActive = !s.resolutionActive // Toggle dropdown state
	}

	// Apply resolution change when dropdown is closed
	if !s.resolutionActive {
		parts := strings.Split(s.resolutions[s.selectedResolution], "x")
		if len(parts) == 2 {
			fmt.Sscanf(parts[0], "%d", &s.settings.ScreenWidth)
			fmt.Sscanf(parts[1], "%d", &s.settings.ScreenHeight)
		}
	}

}

func (s *SettingsScreen) GetScreenState() map[string]any {
	return s.ScreenState
}
