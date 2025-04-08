package ui

import (
	"strings"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	systems_data "github.com/prestonchoate/space-shmup/Systems/Data"
)

type ShopScreen struct {
	ScreenState map[string]any
}

func (s *ShopScreen) Update(state map[string]any) {
	for key, val := range state {
		s.ScreenState[key] = val
	}
}

func (s *ShopScreen) Draw() {
	w, h := s.getScreenDimensions()
	x := float32(w / 2)
	y := float32(h / 10)

	// Save default styling for later restoration
	defaultLabelAlign := raygui.GetStyle(raygui.DEFAULT, raygui.TEXT_ALIGNMENT)

	// Draw the shop header
	s.drawShopHeader(x, y, w)

	// Draw upgrade options if available
	s.drawUpgradeOptions(x, y, w, h)

	// Restore default styling
	raygui.SetStyle(raygui.DEFAULT, raygui.TEXT_ALIGNMENT, defaultLabelAlign)
}

// Gets the current screen dimensions based on fullscreen state
func (s *ShopScreen) getScreenDimensions() (int, int) {
	if rl.IsWindowFullscreen() {
		return rl.GetMonitorWidth(rl.GetCurrentMonitor()), rl.GetMonitorHeight(rl.GetCurrentMonitor())
	}
	return rl.GetScreenWidth(), rl.GetScreenHeight()
}

// Draws the "SHOP" header text
func (s *ShopScreen) drawShopHeader(x, y float32, w int) {
	defaultLabelTextSize := raygui.GetStyle(raygui.DEFAULT, raygui.TEXT_SIZE)
	raygui.SetStyle(raygui.DEFAULT, raygui.TEXT_SIZE, 60)

	shopText := "SHOP"
	shopTextFontSize := 60
	shopTextWidth := rl.MeasureText(shopText, int32(shopTextFontSize))

	shopTextRect := rl.Rectangle{
		X:      float32((w / 2) - int(shopTextWidth)/2),
		Y:      y - (y / 3),
		Width:  float32(shopTextWidth),
		Height: 30,
	}

	raygui.Label(shopTextRect, "Shop")
	raygui.SetStyle(raygui.DEFAULT, raygui.TEXT_SIZE, defaultLabelTextSize)
}

// Draws the upgrade options or a message if none are available
func (s *ShopScreen) drawUpgradeOptions(x, y float32, w, h int) {
	upgrades, exists := s.ScreenState["upgrades"].([]*systems_data.StatUpgrade)

	if !exists {
		s.drawNoUpgradesMessage(x, y, w)
		return
	}

	lastY := s.drawUpgradeCards(upgrades, x, y, float32(w), float32(h))
	s.drawRerollButton(x, lastY, float32(w), float32(h), upgrades)
}

// Draws a message when no upgrades are available
func (s *ShopScreen) drawNoUpgradesMessage(x, y float32, w int) {
	raygui.Label(rl.Rectangle{
		X:      x,
		Y:      y + 20,
		Width:  float32(w / 3),
		Height: 30,
	}, "No Upgrades Available")
}

// Draws the upgrade cards in a row
func (s *ShopScreen) drawUpgradeCards(upgrades []*systems_data.StatUpgrade, x, y, w, h float32) float32 {
	if len(upgrades) == 0 {
		return 0
	}

	const maxCardWidth = float32(200)
	const minPadding = float32(20)
	const rowSpacing = float32(30)
	const cardHeight = float32(140)

	startY := y + 30

	// Estimate how many cards can fit per row
	perRow := max(int((w+minPadding)/(maxCardWidth+minPadding)), 1)

	total := len(upgrades)
	rows := (total + perRow - 1) / perRow

	lastY := float32(0)

	for row := range rows {
		startIndex := row * perRow
		endIndex := min(startIndex+perRow, total)

		numCards := float32(endIndex - startIndex)

		// Compute actual card width to fit within screen width
		cardWidth := min((w-minPadding*(numCards+1))/numCards, maxCardWidth)

		// Recalculate horizontal padding for even spacing
		totalCardsWidth := cardWidth * numCards
		extraSpace := w - totalCardsWidth
		padding := extraSpace / (numCards + 1)

		currentX := padding
		currentY := startY + float32(row)*(cardHeight+rowSpacing)

		for i := startIndex; i < endIndex; i++ {
			upgrade := upgrades[i]
			if upgrade == nil {
				continue
			}

			s.drawUpgradeCard(upgrade, currentX, currentY, cardWidth, cardHeight)
			currentX += cardWidth + padding
		}
		lastY = startY + float32(row)*(cardHeight+rowSpacing) + cardHeight
	}
	return lastY
}

// Draws a single upgrade card with glow effect
func (s *ShopScreen) drawUpgradeCard(upgrade *systems_data.StatUpgrade, x, y, width, height float32) {
	raygui.SetStyle(raygui.LABEL, raygui.TEXT_ALIGNMENT, raygui.TEXT_ALIGN_CENTER)

	panelRect := rl.Rectangle{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}

	// Draw the glow effect
	s.drawCardGlowEffect(panelRect, upgrade.Tier.Color)

	// Draw the panel and its contents
	raygui.Panel(panelRect, upgrade.Tier.Name)
	s.drawCardContent(x, y, width, height, panelRect, upgrade)
}

// Draws the glow effect for a card
func (s *ShopScreen) drawCardGlowEffect(panelRect rl.Rectangle, color rl.Color) {
	glowSize := float32(10)
	glowColor := color
	//glowColor.A = 200 // Semi-transparent

	rl.DrawRectangle(
		int32(panelRect.X-glowSize),
		int32(panelRect.Y-glowSize),
		int32(panelRect.Width+2*glowSize),
		int32(panelRect.Height+2*glowSize),
		glowColor,
	)
}

// Draws the content of an upgrade card
func (s *ShopScreen) drawCardContent(x, y, width, height float32, panelRect rl.Rectangle, upgrade *systems_data.StatUpgrade) {
	titleFontSize := raygui.GetStyle(raygui.DEFAULT, raygui.TEXT_SIZE)
	sidePadding := float32(10)
	topPadding := float32(10)

	textRect := rl.Rectangle{
		X:      x + sidePadding,
		Y:      y + float32(titleFontSize) + topPadding*1.5,
		Width:  panelRect.Width - sidePadding*2,
		Height: height - float32(titleFontSize) - topPadding*2.5,
	}

	text := upgrade.String()
	s.drawWrappedCenteredText(text, textRect, int32(raygui.GetFont().BaseSize))
}

// Draws the reroll button
func (s *ShopScreen) drawRerollButton(x, y, w, h float32, upgrades []*systems_data.StatUpgrade) {
	buttonWidth := float32(150)
	buttonHeight := float32(30)
	paddingTop := float32(30)
	buttonX := x - buttonWidth/2
	buttonY := y + paddingTop
	s.ScreenState["rerollButtonPressed"] = raygui.Button(rl.Rectangle{
		X:      buttonX,
		Y:      buttonY,
		Width:  buttonWidth,
		Height: buttonHeight,
	}, "Reroll Upgrades")
}

func (s *ShopScreen) GetScreenState() map[string]any {
	return s.ScreenState
}

// Custom function to draw wrapped text centered in a panel
func (s *ShopScreen) drawWrappedCenteredText(text string, rect rl.Rectangle, fontSize int32) {
	lines := s.wordWrapText(text, int(rect.Width)-20, fontSize)

	// Calculate starting Y position to center text block vertically
	totalTextHeight := float32(int32(len(lines)) * (fontSize + 2))
	startY := rect.Y + (rect.Height-totalTextHeight)/2

	// Draw each line centered horizontally
	for i, line := range lines {
		textWidth := rl.MeasureText(line, fontSize)
		posX := rect.X + (rect.Width-float32(textWidth))/2
		posY := startY + float32(int32(i)*(fontSize+2))

		raygui.Label(rl.Rectangle{
			X:      posX,
			Y:      posY,
			Width:  float32(textWidth),
			Height: float32(fontSize),
		}, line)
	}
}

// Helper function to word wrap text
func (s *ShopScreen) wordWrapText(text string, maxWidth int, fontSize int32) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	lines := make([]string, 0)
	currentLine := ""

	for _, word := range words {
		// Check if adding the word would exceed max width
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		lineWidth := rl.MeasureText(testLine, fontSize)

		if float32(lineWidth) <= float32(maxWidth) || currentLine == "" {
			// Word fits on current line, add it
			currentLine = testLine
		} else {
			// Word doesn't fit, start a new line
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	// Add the last line if it contains text
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}
