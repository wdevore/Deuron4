package app

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

// filterEvent returns false if it handled the event. Returning false
// prevents the event from being added to the queue.
func (v *App) filterEvent(e sdl.Event, userdata interface{}) bool {

	switch t := e.(type) {
	case *sdl.QuitEvent:
		fmt.Println("SDL Quit event")
		v.running = false
		return false // We handled it. Don't allow it to be added to the queue.
	case *sdl.MouseMotionEvent:
		// fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
		// 	t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
		return false // We handled it. Don't allow it to be added to the queue.
	// case *sdl.MouseButtonEvent:
	// 	fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
	// 		t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
	case *sdl.MouseWheelEvent:
		// -x = fingers moving left
		// fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
		// 	t.Timestamp, t.Type, t.Which, t.X, t.Y)
		return false
	case *sdl.KeyboardEvent:

		if t.State == sdl.RELEASED {
			// fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tScode:%d\tstate:%d\trepeat:%d\n",
			// 	t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.Keysym.Scancode, t.State, t.Repeat)

			switch t.Keysym.Scancode {
			case sdl.SCANCODE_ESCAPE:
				v.running = false
				break
			default:
				if v.mode == "Main" {
					// switch maps
					kmode := codeToInt(t.Keysym.Scancode)
					if kmode < 0 {
						// leave main mode
						v.mode = "Entry"
						// Entry mode. Route to active key map
						kmap := v.keyMaps[v.mapIdx]
						if kmap != nil {
							kmap.reset()
							km := kmap.handle(t.Keysym.Scancode, v.mode)

							if km != "" {
								v.mode = km
							}
						} else {
							fmt.Printf("No key map for (%d)\n", v.mapIdx)
						}
						return false
					}

					if kmode != v.mapIdx {
						fmt.Printf("Switching from (%d) map to (%d) map\n", v.mapIdx, kmode)
						v.mapIdx = kmode
					}
				} else {
					// Entry mode. Route to key active map
					kmap := v.keyMaps[v.mapIdx]
					km := kmap.handle(t.Keysym.Scancode, v.mode)
					if km != "" {
						v.mode = km
					}
				}
				break
			}
		}
		return false // We handled it. Don't allow it to be added to the queue.
	}

	return true
}

type IKeyMap interface {
	handle(code sdl.Scancode, mode string) string
	reset()
}

type keyMapBase struct {
	mode     string
	value    string
	App      *App
	property string
	field    string
	cmd      []string
}

func (km *keyMapBase) handle(code sdl.Scancode, mode string) string {
	// Are they cancelling while in Entry mode
	km.mode = mode

	// fmt.Printf("mode %s\n", mode)
	if mode == "Entry" {
		switch code {
		case sdl.SCANCODE_RIGHTBRACKET:
			fmt.Println("Cancelling, changing back to Main")
			km.App.SetText(km.property+" "+km.field, "--")
			km.mode = "Main"
			break
		}
	}

	return km.mode
}

func (km *keyMapBase) reset() {
	km.value = ""
}

func codeToInt(code sdl.Scancode) int {
	switch code {
	case sdl.SCANCODE_0:
		return 0
	case sdl.SCANCODE_1:
		return 1
	case sdl.SCANCODE_2:
		return 2
	case sdl.SCANCODE_3:
		return 3
	case sdl.SCANCODE_4:
		return 4
	case sdl.SCANCODE_5:
		return 5
	case sdl.SCANCODE_6:
		return 6
	case sdl.SCANCODE_7:
		return 7
	case sdl.SCANCODE_8:
		return 8
	case sdl.SCANCODE_9:
		return 9
	}
	return -1
}

func codeToString(code sdl.Scancode) string {
	switch code {
	case sdl.SCANCODE_0:
		return "0"
	case sdl.SCANCODE_1:
		return "1"
	case sdl.SCANCODE_2:
		return "2"
	case sdl.SCANCODE_3:
		return "3"
	case sdl.SCANCODE_4:
		return "4"
	case sdl.SCANCODE_5:
		return "5"
	case sdl.SCANCODE_6:
		return "6"
	case sdl.SCANCODE_7:
		return "7"
	case sdl.SCANCODE_8:
		return "8"
	case sdl.SCANCODE_9:
		return "9"
	case sdl.SCANCODE_PERIOD:
		return "."
	}
	return ""
}
