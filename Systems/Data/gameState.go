package systems_data

import "fmt"

type GameState int

const (
	None GameState = iota
	Start
	Settings
	Loading
	Paused
	Playing
	Shop
	GameOver
	Restart
	Exit
)

// String implements the fmt.Stringer interface for GameState
func (g GameState) String() string {
	switch g {
	case None:
		return "None"
	case Start:
		return "Start"
	case Settings:
		return "Settings"
	case Loading:
		return "Loading"
	case Paused:
		return "Paused"
	case Playing:
		return "Playing"
	case Shop:
		return "Shop"
	case GameOver:
		return "GameOver"
	case Restart:
		return "Restart"
	case Exit:
		return "Exit"
	default:
		return fmt.Sprintf("GameState(%d)", int(g))
	}
}
