package ui

import (
	"fmt"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	events "github.com/prestonchoate/space-shmup/Systems/Events"
	events_data "github.com/prestonchoate/space-shmup/Systems/Events/Data"
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
	submitButtonText := "SUBMIT"

	buttonWidth := rl.MeasureText(restartButtonText, 40)
	buttonHeight := 80
	screenWidth := rl.GetScreenWidth()
	screenHeight := rl.GetScreenHeight()

	spacing := int32(20)
	totalWidth := (3 * buttonWidth) + spacing
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

	// Initialize initials field if it doesn't exist yet
	if _, ok := g.ScreenState["playerInitials"]; !ok {
		g.ScreenState["playerInitials"] = ""
	}

	// Cast the initials to string to ensure proper type
	initialsStr, ok := g.ScreenState["playerInitials"].(string)
	if !ok {
		initialsStr = ""
		g.ScreenState["playerInitials"] = initialsStr
	}

	// Label for initials input
	rl.DrawText(
		"Enter your initials:",
		int32(rl.GetScreenWidth()/2)-rl.MeasureText("Enter your initials:", 20)/2,
		int32(rl.GetScreenHeight()/2)+90,
		20,
		rl.Gray,
	)

	// Input field for initials
	inputRect := rl.Rectangle{
		X:      float32(rl.GetScreenWidth()/2 - 50),
		Y:      float32(rl.GetScreenHeight()/2) + 120,
		Width:  100,
		Height: 40,
	}

	// Character limit message
	rl.DrawText(
		"(3 characters max)",
		int32(rl.GetScreenWidth()/2)-rl.MeasureText("(3 characters max)", 16)/2,
		int32(inputRect.Y+inputRect.Height+10),
		16,
		rl.Gray,
	)

	initialsEditMode, exists := g.ScreenState["initialsEditMode"].(bool)
	if !exists {
		initialsEditMode = false
		g.ScreenState["initialsEditMode"] = false
	}
	// TextBox for initials input with 3 char limit + the null terminator for C style strings
	if raygui.TextBox(inputRect, &initialsStr, 4, initialsEditMode) {
		initialsEditMode = !initialsEditMode
		g.ScreenState["initialsEditMode"] = initialsEditMode
	}

	// Limit initials to 3 characters
	if len(initialsStr) > 3 {
		initialsStr = initialsStr[:3]
	}

	// Force uppercase for initials
	if len(initialsStr) > 0 {
		for i := range len(initialsStr) {
			if initialsStr[i] >= 'a' && initialsStr[i] <= 'z' {
				initialsRunes := []rune(initialsStr)
				initialsRunes[i] = rune(initialsStr[i] - 32) // Convert to uppercase
				initialsStr = string(initialsRunes)
			}
		}
	}

	g.ScreenState["playerInitials"] = initialsStr

	g.ScreenState["restartButtonPressed"] = raygui.Button(rl.Rectangle{
		X:      float32(startX),
		Y:      float32(buttonY),
		Width:  float32(buttonWidth),
		Height: float32(buttonHeight),
	},
		restartButtonText,
	)

	// Submit button - only active if initials are entered
	initialsEntered := len(initialsStr) > 0
	if !initialsEntered {
		// Draw a disabled button
		raygui.SetState(raygui.STATE_DISABLED)
		raygui.Button(rl.Rectangle{
			X:      float32(startX + int(spacing) + int(buttonWidth)),
			Y:      float32(buttonY),
			Width:  float32(buttonWidth),
			Height: float32(buttonHeight),
		},
			submitButtonText,
		)
		raygui.SetState(raygui.STATE_NORMAL)
	} else {
		submit := raygui.Button(rl.Rectangle{
			X:      float32(startX + int(spacing) + int(buttonWidth)),
			Y:      float32(buttonY),
			Width:  float32(buttonWidth),
			Height: float32(buttonHeight),
		},
			submitButtonText,
		)

		data, ok := g.ScreenState["scoreSubmitted"]
		var scoreSubmitted bool
		if ok {
			scoreSubmitted = data.(bool)
		}

		if !ok || !scoreSubmitted {
			if submit {
				// Here we would submit the score with the initials
				g.ScreenState["submitButtonPressed"] = true
				g.ScreenState["scoreSubmitted"] = true
			}
		} else if scoreSubmitted && submit {
			events.GetEventManagerInstance().Emit(events_data.AddMessage, events_data.AddMessageData{
				Message: "Score submission already attempted",
				Timer:   2.0,
				Color:   rl.Red,
			})
		}
	}

	/*

		submit := raygui.Button(rl.Rectangle{
			X:      float32(startX + int(spacing) + int(buttonWidth)),
			Y:      float32(buttonY),
			Width:  float32(buttonWidth),
			Height: float32(buttonHeight),
		},
			submitButtonText,
		)

		data, ok := g.ScreenState["scoreSubmitted"]

		var scoreSubmitted bool
		if ok {
			scoreSubmitted = data.(bool)
		}

		if !ok || !scoreSubmitted {
			g.ScreenState["submitButtonPressed"] = submit
		} else if scoreSubmitted && submit {
			events.GetEventManagerInstance().Emit(events_data.AddMessage, events_data.AddMessageData{
				Message: "Score submission already attempted",
				Timer:   2.0,
				Color:   rl.Red,
			})
		}

	*/

	g.ScreenState["exitButtonPressed"] = raygui.Button(rl.Rectangle{
		X:      float32(startX + (int(spacing)+int(buttonWidth))*2),
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
