package systems_data

type GameState int

const (
	Start GameState = iota
	Settings
	Loading
	Paused
	Playing
	Menu
	GameOver
	Restart
	Exit
)
