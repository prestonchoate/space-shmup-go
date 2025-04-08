package ui

import (
	"strings"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	systems_data "github.com/prestonchoate/space-shmup/Systems/Data"
)

type ShopScreen struct {
	ScreenState   map[string]any
	glowShader    rl.Shader
	shaderLoaded  bool
	locGlowColor  int32
	locResolution int32
	locRectBounds int32
	locSoftness   int32
	locIntensity  int32
}

func (s *ShopScreen) Update(state map[string]any) {
	for key, val := range state {
		s.ScreenState[key] = val
	}
	if h, exists := s.ScreenState["hoverStates"]; !exists || h == nil {
		s.ScreenState["hoverStates"] = make(map[*systems_data.StatUpgrade]float32)
	}

	if !s.shaderLoaded {
		s.initShaders()
	}
}

// Initialize shaders used for glow effects
func (s *ShopScreen) initShaders() {
	// Vertex shader remains unchanged from default
	vsCode := `
	#version 330
	in vec3 vertexPosition;
	in vec2 vertexTexCoord;
	in vec4 vertexColor;
	out vec2 fragTexCoord;
	out vec4 fragColor;
	uniform mat4 mvp;
	void main() {
		fragTexCoord = vertexTexCoord;
		fragColor = vertexColor;
		gl_Position = mvp*vec4(vertexPosition, 1.0);
	}`

	// Fragment shader for rectangular glow effect
	fsCode := `
	#version 330
	in vec2 fragTexCoord;
	in vec4 fragColor;
	out vec4 finalColor;
	uniform vec4 glowColor;
	uniform vec2 resolution;
	uniform vec4 rectBounds;  // x, y, width, height
	uniform float softness;
	uniform float intensity;
	
	// Smoothstep function for soft edges
	float boxSDF(vec2 p, vec2 size) {
		vec2 d = abs(p) - size;
		return length(max(d, 0.0)) + min(max(d.x, d.y), 0.0);
	}
	
	void main() {
		// Current fragment position (correcting Y-axis)
		vec2 pixelPos = vec2(gl_FragCoord.x, resolution.y - gl_FragCoord.y);
		
		// Rectangle center
		vec2 rectCenter = vec2(rectBounds.x + rectBounds.z/2.0, rectBounds.y + rectBounds.w/2.0);
		
		// Rectangle half size
		vec2 halfSize = vec2(rectBounds.z, rectBounds.w) / 2.0;
		
		// Calculate distance from current pixel to rectangle
		float dist = boxSDF(pixelPos - rectCenter, halfSize);
		
		// Create soft glow with configurable softness
		float glow = exp(-dist * dist / (softness * softness)) * intensity;
		
		// Blend the glow with the base color
		finalColor = mix(vec4(0,0,0,0), glowColor, glow);
		finalColor.a = min(finalColor.a, glowColor.a * glow);
	}`

	s.glowShader = rl.LoadShaderFromMemory(vsCode, fsCode)

	// Get shader uniform locations directly
	s.locGlowColor = rl.GetShaderLocation(s.glowShader, "glowColor")
	s.locResolution = rl.GetShaderLocation(s.glowShader, "resolution")
	s.locRectBounds = rl.GetShaderLocation(s.glowShader, "rectBounds")
	s.locSoftness = rl.GetShaderLocation(s.glowShader, "softness")
	s.locIntensity = rl.GetShaderLocation(s.glowShader, "intensity")

	s.shaderLoaded = true
}

// Clean up resources when the screen is closed or unloaded
func (s *ShopScreen) CleanUp() {
	if s.shaderLoaded {
		rl.UnloadShader(s.glowShader)
		s.shaderLoaded = false
	}
}

func (s *ShopScreen) Draw() {
	w, h := s.getScreenDimensions()
	x := float32(w / 2)
	y := float32(h / 4)

	// Save default styling for later restoration
	defaultLabelAlign := raygui.GetStyle(raygui.DEFAULT, raygui.TEXT_ALIGNMENT)

	// Draw the shop header
	s.drawShopHeader(y, w)

	// Draw upgrade options if available
	s.drawUpgradeOptions(x, y, w)

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
func (s *ShopScreen) drawShopHeader(y float32, w int) {
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
func (s *ShopScreen) drawUpgradeOptions(x, y float32, w int) {
	upgrades, exists := s.ScreenState["upgrades"].([]*systems_data.StatUpgrade)

	if !exists {
		s.drawNoUpgradesMessage(x, y, w)
		return
	}

	lastY := s.drawUpgradeCards(upgrades, y, float32(w))
	s.drawRerollButton(x, lastY)
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
func (s *ShopScreen) drawUpgradeCards(upgrades []*systems_data.StatUpgrade, y, w float32) float32 {
	if len(upgrades) == 0 {
		return 0
	}

	const maxCardWidth = float32(200)
	const minPadding = float32(50)
	const rowSpacing = float32(60)
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
	hoverStates, exists := s.ScreenState["hoverStates"].(map[*systems_data.StatUpgrade]float32)
	if !exists || hoverStates == nil {
		hoverStates = map[*systems_data.StatUpgrade]float32{}
		s.ScreenState["hoverStates"] = hoverStates
	}

	raygui.SetStyle(raygui.LABEL, raygui.TEXT_ALIGNMENT, raygui.TEXT_ALIGN_CENTER)

	mouse := rl.GetMousePosition()
	cardRect := rl.Rectangle{X: x, Y: y, Width: width, Height: height}
	hovered := rl.CheckCollisionPointRec(mouse, cardRect)

	// Animate hover
	const animSpeed = float32(0.1)
	progress := hoverStates[upgrade]
	if hovered {
		progress += animSpeed
		if progress > 1 {
			progress = 1
		}
	} else {
		progress -= animSpeed
		if progress < 0 {
			progress = 0
		}
	}
	hoverStates[upgrade] = progress

	// Apply scaling
	scale := 1 + 0.05*progress
	scaledWidth := width * scale
	scaledHeight := height * scale
	scaledX := x - (scaledWidth-width)/2
	scaledY := y - (scaledHeight-height)/2
	scaledRect := rl.Rectangle{
		X:      scaledX,
		Y:      scaledY,
		Width:  scaledWidth,
		Height: scaledHeight,
	}

	// Handle click
	if hovered && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		s.ScreenState["selectedUpgrade"] = upgrade
	}

	// More glow when hovered
	glowAlpha := uint8(100 + 80*progress)
	glowColor := upgrade.Tier.Color
	glowColor.A = glowAlpha
	if upgrade.Tier.Name == "Common" {
		s.drawCardGlowEffect(scaledRect, glowColor)
	} else {
		s.drawCardGlowEffectWithShader(scaledRect, glowColor)
	}

	// Draw the panel
	raygui.Panel(scaledRect, upgrade.Tier.Name)

	// Draw content inside
	s.drawCardContent(scaledX, scaledY, scaledHeight, scaledRect, upgrade)
	s.ScreenState["hoverStates"] = hoverStates
}

// Draws the glow effect for a card using shader
func (s *ShopScreen) drawCardGlowEffectWithShader(panelRect rl.Rectangle, color rl.Color) {
	if !s.shaderLoaded {
		return
	}

	// Save current blend mode and set to additive for glow effect
	rl.BeginBlendMode(rl.BlendAdditive)

	// Expanded rectangle for the glow effect rendering area
	glowSize := float32(40) // Larger glow size for more visible effect
	glowRect := rl.Rectangle{
		X:      panelRect.X - glowSize*.5,
		Y:      panelRect.Y - glowSize*1.75,
		Width:  panelRect.Width + 2*glowSize,
		Height: panelRect.Height + 2*glowSize,
	}

	// Get screen dimensions for shader resolution
	w, h := s.getScreenDimensions()
	resolution := [2]float32{float32(w), float32(h)}

	// Convert color to normalized floats for shader
	colorVec := [4]float32{
		float32(color.R) / 255.0,
		float32(color.G) / 255.0,
		float32(color.B) / 255.0,
		float32(color.A) / 255.0,
	}

	// Set rectangle bounds for the shader
	rectBounds := [4]float32{
		panelRect.X, panelRect.Y,
		panelRect.Width, panelRect.Height,
	}

	// Set shader uniforms using stored locations
	rl.SetShaderValue(s.glowShader, s.locGlowColor, colorVec[:], rl.ShaderUniformVec4)
	rl.SetShaderValue(s.glowShader, s.locResolution, resolution[:], rl.ShaderUniformVec2)
	rl.SetShaderValue(s.glowShader, s.locRectBounds, rectBounds[:], rl.ShaderUniformVec4)

	// Control the softness of the glow falloff
	softness := float32(30.0) // Higher value = softer glow
	rl.SetShaderValue(s.glowShader, s.locSoftness, []float32{softness}, rl.ShaderUniformFloat)

	// Control the intensity of the glow
	intensity := float32(1.8) // Higher value = brighter glow
	rl.SetShaderValue(s.glowShader, s.locIntensity, []float32{intensity}, rl.ShaderUniformFloat)

	// Begin shader mode
	rl.BeginShaderMode(s.glowShader)

	// Draw a rectangle with the shader applied
	rl.DrawRectangleRec(glowRect, rl.White)

	// End shader mode
	rl.EndShaderMode()

	// Restore previous blend mode
	rl.EndBlendMode()
}

// Draw normal "glow" effect
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
func (s *ShopScreen) drawCardContent(x, y, height float32, panelRect rl.Rectangle, upgrade *systems_data.StatUpgrade) {
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
func (s *ShopScreen) drawRerollButton(x, y float32) {
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
