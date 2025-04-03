package events_data

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	systems_data "github.com/prestonchoate/space-shmup/Systems/Data"
)

type ChangeStateData struct {
	NewState systems_data.GameState
}

type UpdateSettingsData struct {
	NewSettings systems_data.GameSettings
}

type ReturnStateData struct{}

type HighScoreData struct {
	Initials string
	Score    int64
}

type ScoreSubmissionCompleteData struct {
	Success bool
}

type AddMessageData struct {
	Message string
	Timer   float32
	Color   rl.Color
}

type StatUpgradeData struct {
	StatType string
	Value    float32
}
