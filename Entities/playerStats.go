package entities

// PlayerStats encapsulates all statistics for the player
type PlayerStats struct {
	Speed     float32
	Health    int
	MaxHealth int
	FireRate  float32
	Damage    float32
	Size      float32
	// You can add more stats here as needed
	// Examples:
	// CritChance  float32
	// ShieldValue int
	// ReloadSpeed float32
}

// NewDefaultPlayerStats creates a new PlayerStats with default values
func NewDefaultPlayerStats() PlayerStats {
	return PlayerStats{
		Speed:     350,
		Health:    100,
		MaxHealth: 100,
		FireRate:  DEFAULT_FIRE_RATE,
		Damage:    10.0,
		Size:      1.0,
	}
}
