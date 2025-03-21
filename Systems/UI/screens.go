package ui

type Screens interface {
	Update(state map[string]any)
	Draw()
	GetScreenState() map[string]any
}
