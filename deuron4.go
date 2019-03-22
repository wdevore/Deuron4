package main

// Note: You may need a ram drive:
// diskutil erasevolume HFS+ 'RAMDisk' `hdiutil attach -nomount ram://2097152`

import (
	"fmt"

	"github.com/wdevore/Deuron4/deuron/app"
)

/*
	Deuron4 has a GUI for display.
	Keyboard input for control
	Debug output to the console.

	typical sequence:
	>con
	>load
	>start
*/
var gview *app.App // The GUI

// This is the main entry point for Deuron4.
// It starts both the GUI and TUI.
func main() {
	gview = app.NewApp()
	defer gview.Close()

	gview.Open()

	gview.SetFont("Roboto-Bold.ttf", 24)
	gview.Configure()

	// TODO replace with regexp in app.
	// go readConsole()

	gview.Run()
}

// TODO migrate to app scanning and handling.
// This causes issues during debugging.
func readConsole() {
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Println("Enter 'help' for console commands.")

	// var text = ""
	// fmt.Print("]")

	// for text != "quit\n" {
	// 	text, _ = reader.ReadString('\n')
	// 	text = strings.Trim(text, "\n")
	// 	args := strings.Split(text, " ")
	// 	// fmt.Printf("args: %v\n", args)
	// 	switch args[0] {
	// 	case "help":
	// 		printHelp()
	// 	default:
	// 		gview.Command(args)
	// 	}
	// 	fmt.Print("]")
	// }

	// fmt.Println("Console exited.")
}

func printHelp() {
	fmt.Println("------------------- Help ---------------------------------")
	fmt.Println("'quit' stops any simulation and exits app.")
	fmt.Println("'help' this help screen.")
	fmt.Println("'start' starts the current target simulation.")
	fmt.Println("'stop' stops the current target simulation.")
	fmt.Println("'load' reads target's json.")
	fmt.Println("'go' connects, loads and runs sim.")
	fmt.Println("'set sim-name' sets the target simulation, where")
	fmt.Println("   sim-name specifies a json file in the working directory.")
	fmt.Println("'con' connects to a target sim-name. It does NOT start it.")
	fmt.Println("'type' changes sim type: `runreset` or `continous`")
	fmt.Println("'ping' sends `ping` to target sim.")

	// fmt.Println("'p' activates property mode and lists available properties.")
	// fmt.Println("  you then enter <property number> and <value>")
	fmt.Println("'\\' shows what properties can be changed. To change, for example,")
	fmt.Println("   Poisson-min value enter 1 1 <value>")
	fmt.Println("----------------------------------------------------------")
}
