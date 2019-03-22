package app

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

// -----------------------------------------------------------------
// Key map 1
// -----------------------------------------------------------------

type keyMap1 struct {
	keyMapBase
}

func NewKeyMap1(App *App) IKeyMap {
	km := new(keyMap1)
	km.App = App
	return km
}

func (km *keyMap1) handle(code sdl.Scancode, mode string) string {
	kmo := km.keyMapBase.handle(code, mode)
	if kmo == "Main" || kmo == "handled" {
		return km.mode
	}

	switch code {
	case sdl.SCANCODE_Q:
		km.property = "ExpoFunc"
		km.field = "A"
		fmt.Println("Enter Expo `A` value")
		v := km.App.RequestAppProperty(km.property + " " + km.field)
		km.App.SetText(km.property+" "+km.field, v)
		km.cmd = []string{km.property, km.field, v}
		km.App.SetAppCommand(km.cmd)
		break
	case sdl.SCANCODE_W:
		km.property = "ExpoFunc"
		km.field = "Tau"
		fmt.Println("Enter Expo tau value")
		v := km.App.RequestAppProperty(km.property + " " + km.field)
		km.App.SetText(km.property+" "+km.field, v)
		km.cmd = []string{km.property, km.field, v}
		km.App.SetAppCommand(km.cmd)
		break
	case sdl.SCANCODE_E:
		km.property = "ExpoFunc"
		km.field = "M"
		fmt.Println("Enter Expo M value")
		v := km.App.RequestAppProperty(km.property + " " + km.field)
		km.App.SetText(km.property+" "+km.field, v)
		km.cmd = []string{km.property, km.field, v}
		km.App.SetAppCommand(km.cmd)
		break
	case sdl.SCANCODE_R:
		km.property = "ExpoFunc"
		km.field = "WMax"
		fmt.Println("Enter Expo WMax value")
		v := km.App.RequestAppProperty(km.property + " " + km.field)
		km.App.SetText(km.property+" "+km.field, v)
		km.cmd = []string{km.property, km.field, v}
		km.App.SetAppCommand(km.cmd)
		break
	case sdl.SCANCODE_T:
		km.property = "Inc"
		km.field = "Size"
		fmt.Println("Enter Increment size value")
		v := km.App.RequestProperty(km.property + " " + km.field)
		km.App.SetText(km.property+" "+km.field, v)
		break
	case sdl.SCANCODE_Y:
		km.property = "Dec"
		km.field = "Size"
		fmt.Println("Enter Decrement size value")
		v := km.App.RequestProperty(km.property + " " + km.field)
		km.App.SetText(km.property+" "+km.field, v)
		break
	case sdl.SCANCODE_UP:
		km.App.IncrementAppProperty(km.cmd)
		v := km.App.RequestAppProperty(km.property + " " + km.field)
		km.App.SetText(km.property+" "+km.field, v)
		break
	case sdl.SCANCODE_DOWN:
		// This message eventually reaches the simulation.
		// Increment the current active property.
		km.App.DecrementAppProperty(km.cmd)
		v := km.App.RequestAppProperty(km.property + " " + km.field)
		km.App.SetText(km.property+" "+km.field, v)
		break
	case sdl.SCANCODE_G:
		cmd := []string{"go"}
		km.App.Command(cmd)
		break
	case sdl.SCANCODE_APOSTROPHE: // '
		break
	case sdl.SCANCODE_SEMICOLON: // ,
		break
	case sdl.SCANCODE_COMMA:
		break
	case sdl.SCANCODE_SLASH:
		break
	case sdl.SCANCODE_RETURN:
		fmt.Printf("Entered: (%s), changing back to main.\n", km.value)
		// Send property to sim.
		cmd := []string{km.property, km.field, km.value}
		km.App.SetAppCommand(cmd)

		km.mode = "Main"
	default:
		km.value = km.value + codeToString(code)
		// fmt.Printf("val: %s\n", km.value)
		// Update gui
		km.App.SetText(km.property+" "+km.field, km.value)
		break
	}

	return km.mode
}
