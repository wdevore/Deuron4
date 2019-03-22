package app

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

// -----------------------------------------------------------------
// Key map 0
// -----------------------------------------------------------------

type keyMap0 struct {
	keyMapBase
}

func NewKeyMap0(App *App) IKeyMap {
	km := new(keyMap0)
	km.App = App
	return km
}

func (km *keyMap0) handle(code sdl.Scancode, mode string) string {
	kmo := km.keyMapBase.handle(code, mode)
	if kmo == "Main" || kmo == "handled" {
		return km.mode
	}

	switch code {
	case sdl.SCANCODE_Q:
		km.property = "Poisson"
		km.field = "Max"
		fmt.Println("Enter poisson max value")
		v := km.App.RequestProperty(km.property + " " + km.field)
		km.App.SetText(km.property+" "+km.field, v)
		km.App.SetCommand([]string{km.property, km.field, v})
		break
	case sdl.SCANCODE_W:
		km.property = "Poisson"
		km.field = "Min"
		fmt.Println("Enter poisson min value")
		v := km.App.RequestProperty(km.property + " " + km.field)
		km.App.SetText(km.property+" "+km.field, v)
		km.App.SetCommand([]string{km.property, km.field, v})
		break
	case sdl.SCANCODE_E:
		km.field = "Spread"
		km.property = "Poisson"
		fmt.Println("Enter poisson spread value")
		v := km.App.RequestProperty(km.property + " " + km.field)
		km.App.SetText(km.property+" "+km.field, v)
		km.App.SetCommand([]string{km.property, km.field, v})
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
		km.App.IncrementProperty()
		return "handled"
	case sdl.SCANCODE_DOWN:
		// This message eventually reaches the simulation.
		// Increment the current active property.
		km.App.DecrementProperty()
		return "handled"
	case sdl.SCANCODE_G:
		cmd := []string{"go"}
		km.App.Command(cmd)
		break
	case sdl.SCANCODE_APOSTROPHE: // '
		// Prepare the sim for running simulations.
		km.App.Command([]string{"create"})
		break
	case sdl.SCANCODE_SEMICOLON: // ,
		// Make a single step in a paused simulation
		cmd := []string{"step"}
		km.App.Command(cmd)
		break
	case sdl.SCANCODE_PERIOD:
		// Run a single pass through a complete simulation then pause.
		cmd := []string{"runPause"}
		km.App.Command(cmd)
		break
	case sdl.SCANCODE_COMMA:
		// Force reset any simulation in progress.
		cmd := []string{"reset"}
		km.App.Command(cmd)
		break
	case sdl.SCANCODE_SLASH:
		// Pause a running simulation
		cmd := []string{"pause"}
		km.App.Command(cmd)
		break
	case sdl.SCANCODE_RETURN:
		fmt.Printf("Entered: (%s), changing back to main.\n", km.value)
		// Send property to sim.
		cmd := []string{"prop", km.property, km.field, km.value}
		km.App.Command(cmd)

		km.mode = "Main"
	default:
		km.value = km.value + codeToString(code)
		// Update gui
		km.App.SetValue(km.value)

		// km.App.SetText(km.property+" "+km.field, km.value)
		break
	}

	return km.mode
}
