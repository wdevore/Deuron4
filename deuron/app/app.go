package app

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/wdevore/Deuron4/deuron/app/graphs"
	"github.com/wdevore/Deuron4/simulation/runreset"
	"github.com/wdevore/Deuron4/simulation/samples"
)

const (
	width  = 1024
	height = 1024
	xpos   = 100
	ypos   = 0
)

// App shows the plots and graphs.
// It receives commands for graphing and viewing various graphs.
type App struct {
	window   *sdl.Window
	surface  *sdl.Surface
	renderer *sdl.Renderer

	// The graphs and text are rendered to this texture.
	texture *sdl.Texture

	running bool

	opened bool

	nFont *Font

	txtSimStatus *Text
	status       string

	txtActiveProperty *Field

	// dynaTxt *DynaText

	spikeGraph  graphs.IGraph
	pointsGraph graphs.IGraph
	expoGraph   graphs.IGraph

	simType string // "runreset" or "continous"
	target  string // Path to simulation

	// comm channel to simulation
	statusComm    chan string
	propEventComm chan string
	requestComm   chan string
	runResetSim   *runreset.RunResetSim

	// Commands
	thing    string
	property string

	// 10 keyboard map layouts
	mode   string // "Main", "Entry"
	mapIdx int

	keyMaps []IKeyMap

	incSize float64
	decSize float64
}

// NewApp creates a new App and initializes it.
func NewApp() *App {
	v := new(App)
	v.opened = false
	v.simType = "runreset"
	v.status = "Stopped"
	v.mode = "Main"
	v.incSize = 10.0
	v.decSize = 10.0
	v.keyMaps = make([]IKeyMap, 10)
	v.keyMaps[0] = NewKeyMap0(v)
	v.keyMaps[1] = NewKeyMap1(v)

	return v
}

// Open shows the App and begin event polling
// (host deuron.IHost)
func (v *App) Open() {
	v.initialize()

	v.opened = true
}

// SetFont sets the font based on path and size.
func (v *App) SetFont(fontPath string, size int) {
	v.nFont = NewFont(fontPath, size)
}

func (v *App) SetText(field, value string) {
	v.txtActiveProperty.SetName(field + ": ")
	v.txtActiveProperty.SetValue(value)
	v.txtActiveProperty.SetName(field + ": ")
}

func (v *App) SetAppCommand(cmd []string) {
	// process command for app
	prop := cmd[0]

	switch prop {
	case "ExpoFunc":
		graph := v.expoGraph.(*graphs.ExpoGraph)
		sub := cmd[1]
		switch sub {
		case "A":
			val := cmd[2]
			fval, _ := strconv.ParseFloat(val, 64)
			graph.SetA(fval)
			break
		case "Tau":
			val := cmd[2]
			fval, _ := strconv.ParseFloat(val, 64)
			graph.SetTau(fval)
			break
		case "M":
			val := cmd[2]
			fval, _ := strconv.ParseFloat(val, 64)
			graph.SetM(fval)
			break
		case "WMax":
			val := cmd[2]
			fval, _ := strconv.ParseFloat(val, 64)
			graph.SetWMax(fval)
			break
		}
		break
	}
}

func (v *App) IncrementAppProperty(cmd []string) float64 {
	prop := cmd[0]

	switch prop {
	case "ExpoFunc":
		graph := v.expoGraph.(*graphs.ExpoGraph)
		what := cmd[1]
		switch what {
		case "A":
			val := graph.A()
			graph.SetA(val + v.incSize)
			return graph.A()
		case "Tau":
			val := graph.Tau()
			graph.SetTau(val + v.incSize)
			return graph.Tau()
		case "M":
			val := graph.M()
			graph.SetM(val + v.incSize)
			return graph.M()
		case "WMax":
			val := graph.WMax()
			graph.SetWMax(val + v.incSize)
			return graph.WMax()
		}
		break
	}

	return 0
}

func (v *App) DecrementAppProperty(cmd []string) {
	prop := cmd[0]

	switch prop {
	case "ExpoFunc":
		graph := v.expoGraph.(*graphs.ExpoGraph)
		what := cmd[1]
		switch what {
		case "A":
			val := graph.A()
			graph.SetA(val - v.decSize)
			break
		case "Tau":
			val := graph.Tau()
			graph.SetTau(val - v.decSize)
			break
		case "M":
			val := graph.M()
			graph.SetM(val - v.decSize)
			break
		case "WMax":
			val := graph.WMax()
			graph.SetWMax(val - v.decSize)
			break
		}
		break
	}
}

func (v *App) RequestAppProperty(property string) string {
	switch property {
	case "ExpoFunc A":
		graph := v.expoGraph.(*graphs.ExpoGraph)
		return fmt.Sprintf("%0.4f", graph.A())
	case "ExpoFunc Tau":
		graph := v.expoGraph.(*graphs.ExpoGraph)
		return fmt.Sprintf("%0.4f", graph.Tau())
	case "ExpoFunc M":
		graph := v.expoGraph.(*graphs.ExpoGraph)
		return fmt.Sprintf("%0.4f", graph.M())
	case "ExpoFunc WMax":
		graph := v.expoGraph.(*graphs.ExpoGraph)
		return fmt.Sprintf("%0.4f", graph.WMax())
	}

	return "??"
}

func (v *App) SetCommand(cmd []string) {
	v.runResetSim.SetCommand(cmd)
}

func (v *App) SetValue(value string) {
	v.txtActiveProperty.SetValue(value)
}

func (v *App) IncrementProperty() {
	cmd := []string{"prop", "up", fmt.Sprintf("%0.2f", v.incSize)}
	v.Command(cmd)
}

func (v *App) DecrementProperty() {
	cmd := []string{"prop", "down", fmt.Sprintf("%0.2f", v.decSize)}
	v.Command(cmd)
}

// Run starts the polling event loop. This must run on
// the main thread.
func (v *App) Run() {
	// log.Println("Starting App polling")
	v.running = true
	sdl.SetEventFilterFunc(v.filterEvent, nil)

	// sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "linear")
	v.renderer.SetDrawColor(64, 64, 64, 255)
	v.renderer.Clear()

	// v.pointsGraph.SetSeries(v.pointsGraph.Accessor())
	v.spikeGraph.SetSeries(nil)

	for v.running {
		sdl.PumpEvents()

		v.renderer.Clear()

		// v.pointsGraph.MarkDirty(true)
		v.spikeGraph.MarkDirty(true)

		if samples.PoiSamples != nil {
			draw := v.spikeGraph.Check()
			if draw {
				v.spikeGraph.DrawAt(0, 100)
			}
		}

		v.expoGraph.DrawAt(0, 300)

		v.txtSimStatus.Draw()
		v.txtActiveProperty.Draw()

		v.window.UpdateSurface()

		// sdl.Delay(17)
		time.Sleep(time.Millisecond * 100)
	}

	v.shutdown()
}

// Quit stops the gui from running, effectively shutting it down.
func (v *App) Quit() {
	v.running = false
}

// Close closes the App.
// Be sure to setup a "defer x.Close()"
func (v *App) Close() {
	if !v.opened {
		return
	}

	log.Println("\nClosing App...")

	log.Println("Destroying font")
	v.nFont.Destroy()
	log.Println("Destroying text(s)")
	v.txtSimStatus.Destroy()
	// v.dynaTxt.Destroy()
	v.txtActiveProperty.Destroy()

	log.Println("Destroying graphs")

	v.spikeGraph.Destroy()
	v.expoGraph.Destroy()

	log.Println("Destroying texture")
	err := v.texture.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Destroying renderer")
	err = v.renderer.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	// log.Println("Shutting down App")
	err = v.window.Destroy()
	sdl.Quit()

	if err != nil {
		log.Fatal(err)
	}
}

func (v *App) initialize() {
	var err error

	err = sdl.Init(sdl.INIT_TIMER | sdl.INIT_VIDEO | sdl.INIT_EVENTS)
	if err != nil {
		panic(err)
	}

	v.window, err = sdl.CreateWindow("Deuron4 Graph", xpos, ypos, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	v.surface, err = v.window.GetSurface()
	if err != nil {
		panic(err)
	}

	v.renderer, err = sdl.CreateSoftwareRenderer(v.surface)
	if err != nil {
		panic(err)
	}

	v.texture, err = v.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		panic(err)
	}

	v.spikeGraph = graphs.NewSpikesGraph(v.renderer, v.texture, 1024, 300)

	v.expoGraph = graphs.NewExpoGraph(v.renderer, v.texture, 1024, 600)
}

// Configure view with draw objects
func (v *App) Configure() {
	fmt.Println("App configuring...")

	v.txtSimStatus = NewText(v.nFont, v.renderer)
	err := v.txtSimStatus.SetText("Status: "+v.status, sdl.Color{R: 127, G: 64, B: 0, A: 255})
	if err != nil {
		v.Close()
		panic(err)
	}

	v.txtActiveProperty = NewField(v.nFont, v.renderer)
	v.txtActiveProperty.SetPosition(5, 30)

	// v.dynaTxt = NewDynaText(v.nFont, v.renderer)
}

// Command handles messages from the console.
func (v *App) Command(args []string) {
	switch args[0] {
	case "quit":
		v.shutdown()
	case "set":
		v.target = args[1]
	case "type":
		v.simType = args[1]
		fmt.Printf("Type switched to `%s`\n", v.simType)
	case "con":
		v.connect()
	case "ping":
		v.runResetSim.Send("ping")
	case "start":
		v.start()
	case "go":
		go v.doit()
	case "stop":
		if v.runResetSim == nil {
			panic("Not connected. Please connect first. Use 'help'.")
		}
		v.runResetSim.Send("stop")
		response := <-v.statusComm
		if response == "Stopped" {
			fmt.Println("Sim requested to stop")
		}
		break
	case "prop":
		// A command relating to a property
		v.runResetSim.Command(args)
		break
	case "step":
		v.runResetSim.Step()
		break
	case "runPause":
		v.runResetSim.RunPause()
		break
	case "reset":
		v.runResetSim.Reset()
		break
	case "create":
		v.create()
		break
	}
}

func (v *App) create() {
	// Connect
	v.connect()

	// Load
	v.runResetSim.Send("load")
	response := <-v.statusComm // wait for response

	if response != "loaded" {
		panic("Unable to load parameters")
	}

	fmt.Println("Loaded")

	v.runResetSim.Create()

	go v.pollForMessage()
	go v.pollForPropertyEvents()
}

func (v *App) doit() {
	// Connect
	v.connect()

	// Load
	v.runResetSim.Send("load")
	response := <-v.statusComm // wait for response

	if response != "loaded" {
		panic("Unable to load parameters")
	}

	fmt.Println("Loaded")

	// Run
	v.status = "Starting..."
	v.txtSimStatus.SetText("Status: "+v.status, sdl.Color{R: 127, G: 64, B: 0, A: 255})

	// fmt.Println("Starting")
	v.start()
}

func (v *App) start() {
	go v.pollForMessage()
	go v.pollForPropertyEvents()

	// This will cause the sim to start the simulation in a coroutine.
	v.runResetSim.Send("start")
	// We start async because we can't lock the app thread
	// from receiveing system events (ex: keyboard)
}

// Runs in a coroutine.
func (v *App) pollForMessage() {
	var response string
	fmt.Print("Waiting for messages from sim...\n")

	for response != "Stopped" {
		response = <-v.statusComm // Wait for response

		v.status = response
		msg := "Status: " + v.status
		v.txtSimStatus.SetText(msg, sdl.Color{R: 255, G: 127, B: 0, A: 255})
	}

	fmt.Printf("Polling exited from: (%s)\n", response)
}

func (v *App) pollForPropertyEvents() {
	var response string
	fmt.Print("Waiting for property events from sim...\n")

	for {
		response = <-v.propEventComm // Wait for response
		split := strings.Split(response, ",")
		fmt.Printf("poll prop: %s, %v\n", response, split)
		v.SetText(split[0], split[1])
	}

}

func (v *App) RequestProperty(property string) string {
	switch property {
	case "Inc Size":
		return fmt.Sprintf("%0.2f", v.incSize)
	case "Dec Size":
		return fmt.Sprintf("%0.2f", v.decSize)
	}

	return v.runResetSim.RequestProperty(property)
}

func (v *App) connect() {
	fmt.Printf("Connecting to `%s`...\n", v.simType)

	if v.runResetSim == nil {
		fmt.Println("Creating sim")
		v.runResetSim = runreset.NewRunResetSim()

		fmt.Println("Creating comm channels")

		v.statusComm = make(chan string)
		v.propEventComm = make(chan string)
		v.requestComm = make(chan string)
		v.runResetSim.Connect(v.statusComm, v.propEventComm, v.requestComm)

		// Test connection to sim
		v.runResetSim.Send("ping")
		response := <-v.statusComm
		if response == "pong" {
			fmt.Println("Connected.")
		} else {
			fmt.Printf("Sim didn't respond to connection correctly (%s)\n", response)
		}
	}
}

func (v *App) shutdown() {
	fmt.Println("Shutting down...")

	if v.runResetSim != nil {
		fmt.Println("Sending simulation the `stop` command...")
		v.runResetSim.Send("stop")
	}

	v.Quit()

	fmt.Println("Done.")
}
