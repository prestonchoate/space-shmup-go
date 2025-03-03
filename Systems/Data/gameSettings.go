package systems_data

type GameSettings struct {
	TargetFPS    int32    `json:"targetFPS"`
	ScreenWidth  int      `json:"screenWidth"`
	ScreenHeight int      `json:"screenHeight"`
	Fullscreen   bool     `json:"fullscreen"`
	Keys         InputMap `json:"keys"`
}
