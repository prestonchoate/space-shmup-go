package ui

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlayingScreen struct {
	ScreenState map[string]any
}

// Update implements Screens.
func (p *PlayingScreen) Update(state map[string]any) {
	for key, val := range state {
		p.ScreenState[key] = val
	}
}

// Draw implements Screens.
func (p *PlayingScreen) Draw() {

	rl.DrawText(
		fmt.Sprintf("%10v%03d", "Health: ", p.getStateValue("health")),
		int32(rl.GetScreenWidth()-200),
		0,
		20,
		rl.Gray)

	rl.DrawText(
		fmt.Sprintf("%10v%06d", "Score: ", p.getStateValue("score")),
		int32(rl.GetScreenWidth()-200),
		30,
		20,
		rl.Gray)

	rl.DrawText(
		fmt.Sprintf("%10v%06d", "Enemies: ", p.getStateValue("enemyCount")),
		int32(rl.GetScreenWidth()-200),
		60,
		20,
		rl.Gray)
}

// GetScreenState implements Screens.
func (p *PlayingScreen) GetScreenState() map[string]any {
	return p.ScreenState
}

func (p *PlayingScreen) getStateValue(key string) any {
	val, exists := p.ScreenState[key]
	if exists {
		return val
	}

	return nil
}
