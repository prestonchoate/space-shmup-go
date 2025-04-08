package systems_data

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type UpgradeType int32
type UpgradeTier struct {
	Name   string
	Color  rl.Color
	Chance float32
}

func (t *UpgradeTier) ColorAsInt() int64 {
	// Convert raylib.Color to int32
	return int64(uint32(t.Color.R)<<24 | uint32(t.Color.G)<<16 | uint32(t.Color.B)<<8 | uint32(t.Color.A))
}

const (
	FlatAmount UpgradeType = iota
	Multiplier
)

type StatUpgrade struct {
	StatName string
	Amount   float32
	Type     UpgradeType
	Tier     UpgradeTier
}

func (s *StatUpgrade) String() string {
	pct := ""
	amt := s.Amount
	firstWord := ""
	if s.StatName == "Size" {
		firstWord = "Decrease"
	} else {
		firstWord = "Increase"
	}
	if s.Type == Multiplier {
		pct = "%"
		amt = s.Amount * 100
	}
	return fmt.Sprintf("%s %s by %.2f%s", firstWord, s.StatName, amt, pct)
}
