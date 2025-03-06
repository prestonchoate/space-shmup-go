package saveManager

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func KeyToString(key int32) string {
	switch key {
	case rl.KeySpace:
		return "SPACE"
	case rl.KeyEscape:
		return "ESCAPE"
	case rl.KeyEnter:
		return "ENTER"
	case rl.KeyTab:
		return "TAB"
	case rl.KeyBackspace:
		return "BACKSPACE"
	case rl.KeyDelete:
		return "DELETE"
	case rl.KeyRight:
		return "RIGHT"
	case rl.KeyLeft:
		return "LEFT"
	case rl.KeyDown:
		return "DOWN"
	case rl.KeyUp:
		return "UP"
	case rl.KeyF1:
		return "F1"
	case rl.KeyF2:
		return "F2"
	case rl.KeyF3:
		return "F3"
	case rl.KeyF4:
		return "F4"
	case rl.KeyF5:
		return "F5"
	case rl.KeyF6:
		return "F6"
	case rl.KeyF7:
		return "F7"
	case rl.KeyF8:
		return "F8"
	case rl.KeyF9:
		return "F9"
	case rl.KeyF10:
		return "F10"
	case rl.KeyF11:
		return "F11"
	case rl.KeyF12:
		return "F12"
	case rl.KeyA:
		return "A"
	case rl.KeyB:
		return "B"
	case rl.KeyC:
		return "C"
	case rl.KeyD:
		return "D"
	case rl.KeyE:
		return "E"
	case rl.KeyF:
		return "F"
	case rl.KeyG:
		return "G"
	case rl.KeyH:
		return "H"
	case rl.KeyI:
		return "I"
	case rl.KeyJ:
		return "J"
	case rl.KeyK:
		return "K"
	case rl.KeyL:
		return "L"
	case rl.KeyM:
		return "M"
	case rl.KeyN:
		return "N"
	case rl.KeyO:
		return "O"
	case rl.KeyP:
		return "P"
	case rl.KeyQ:
		return "Q"
	case rl.KeyR:
		return "R"
	case rl.KeyS:
		return "S"
	case rl.KeyT:
		return "T"
	case rl.KeyU:
		return "U"
	case rl.KeyV:
		return "V"
	case rl.KeyW:
		return "W"
	case rl.KeyX:
		return "X"
	case rl.KeyY:
		return "Y"
	case rl.KeyZ:
		return "Z"
	case rl.KeyZero:
		return "0"
	case rl.KeyOne:
		return "1"
	case rl.KeyTwo:
		return "2"
	case rl.KeyThree:
		return "3"
	case rl.KeyFour:
		return "4"
	case rl.KeyFive:
		return "5"
	case rl.KeySix:
		return "6"
	case rl.KeySeven:
		return "7"
	case rl.KeyEight:
		return "8"
	case rl.KeyNine:
		return "9"
	case rl.KeyLeftShift:
		return "LEFT SHIFT"
	case rl.KeyRightShift:
		return "RIGHT SHIFT"
	case rl.KeyLeftControl:
		return "LEFT CTRL"
	case rl.KeyRightControl:
		return "RIGHT CTRL"
	case rl.KeyLeftAlt:
		return "LEFT ALT"
	case rl.KeyRightAlt:
		return "RIGHT ALT"
	case rl.KeyLeftSuper:
		return "LEFT SUPER"
	case rl.KeyRightSuper:
		return "RIGHT SUPER"
	case rl.KeyCapsLock:
		return "CAPS LOCK"
	case rl.KeyScrollLock:
		return "SCROLL LOCK"
	case rl.KeyNumLock:
		return "NUM LOCK"
	case rl.KeyPrintScreen:
		return "PRINT SCREEN"
	case rl.KeyPause:
		return "PAUSE"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", key)
	}
}
