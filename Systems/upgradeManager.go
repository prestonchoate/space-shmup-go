package systems

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
	systems_data "github.com/prestonchoate/space-shmup/Systems/Data"
)

var upgrader_instance *UpgradeManager

type UpgradeManager struct {
	upgradeTiers []systems_data.UpgradeTier
}

func GetUpgrader() *UpgradeManager {
	if upgrader_instance == nil {
		upgrader_instance = &UpgradeManager{
			upgradeTiers: make([]systems_data.UpgradeTier, 4),
		}
		upgrader_instance.createUpgradeTiers()
	}

	return upgrader_instance
}

func (u *UpgradeManager) createUpgradeTiers() {
	// Common, Rare, Epic, Legendary
	c := systems_data.UpgradeTier{
		Name:   "Common",
		Color:  rl.NewColor(159, 183, 216, 200),
		Chance: 0.6,
	}
	r := systems_data.UpgradeTier{
		Name:   "Rare",
		Color:  rl.NewColor(65, 168, 95, 255),
		Chance: 0.25,
	}
	e := systems_data.UpgradeTier{
		Name:   "Epic",
		Color:  rl.NewColor(57, 21, 184, 255),
		Chance: 0.1,
	}
	l := systems_data.UpgradeTier{
		Name:   "Legendary",
		Color:  rl.NewColor(250, 197, 28, 255),
		Chance: 0.05,
	}

	u.upgradeTiers = append(u.upgradeTiers, c, r, e, l)
}

func (u *UpgradeManager) GetUpgrade() *systems_data.StatUpgrade {
	tier := u.getRandomUpgradeTier()
	stat := u.getRandomStatName()
	upgradeType := u.getRandomUpgradeType()
	var amount float32

	// Define base amounts per stat and tier (very basic for demo)
	base := map[string]float32{
		"Speed":     10,
		"MaxHealth": 5,
		"FireRate":  3,
		"Damage":    1.5,
		"Size":      0.1,
	}

	modifiers := map[string]float32{
		"Common":    1.0,
		"Rare":      1.5,
		"Epic":      2.0,
		"Legendary": 3.0,
	}

	baseAmount := base[stat] * modifiers[tier.Name]

	if upgradeType == systems_data.Multiplier {
		amount = baseAmount / 10 // e.g. 1.15 = +15%
	} else {
		amount = baseAmount
	}

	return &systems_data.StatUpgrade{
		StatName: stat,
		Amount:   amount,
		Type:     upgradeType,
		Tier:     tier,
	}
}

func (u *UpgradeManager) GetUpgrades(num int32) []*systems_data.StatUpgrade {
	upgrades := make([]*systems_data.StatUpgrade, 0)
	for range int(num) {
		upgrades = append(upgrades, u.GetUpgrade())
	}
	return upgrades
}

func (u *UpgradeManager) getRandomUpgradeTier() systems_data.UpgradeTier {
	roll := GetRandomizer().GetPercentage()

	sum := float32(0)

	for _, tier := range u.upgradeTiers {
		sum += tier.Chance
		if roll <= sum {
			return tier
		}
	}

	if sum > 1.0 {
		log.Println("Upgrade Manager: Upgrade Tier chances equal more than 100%. Total chance: ", sum)
	}

	return u.upgradeTiers[0]
}

func (u *UpgradeManager) getRandomUpgradeType() systems_data.UpgradeType {
	if GetRandomizer().GetPercentage() < 0.7 {
		return systems_data.FlatAmount
	}
	return systems_data.Multiplier
}

func (u *UpgradeManager) getRandomStatName() string {
	stats := []string{"Speed", "MaxHealth", "FireRate", "Damage", "Size"}
	return stats[GetRandomizer().GetNumber(int32(len(stats)))]
}
