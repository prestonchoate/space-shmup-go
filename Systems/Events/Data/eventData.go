package events_data

import (
	systems_data "github.com/prestonchoate/space-shmup/Systems/Data"
)

type ChangeStateData struct {
	NewState systems_data.GameState
}
