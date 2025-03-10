package systems_data

type GameState int

const (
	None GameState = iota
	Start
	Settings
	Loading
	Paused
	Playing
	Menu
	GameOver
	Restart
	Exit
)
