package runreset

import (
	"fmt"
	"strings"

	"github.com/wdevore/Deuron4/simulation/samples"
)

/*
This simulation simulates a single neuron by repeatedly
applying stimulus for a time course and then resetting and repeating.
This type of simulation is for tuning the neuron's parameters such that
the neuron maintains and operates within a critical functional region.
*/

type RunResetSim struct {
	statusChannel    chan string
	propEventChannel chan string
	requestChannel   chan string

	stopped bool

	workingPath string

	// Sim ticks at 1ms resolution
	dt float64
	t  int

	// How long to run before resetting.
	runDuration int

	sim *simulation
}

func NewRunResetSim() *RunResetSim {
	s := new(RunResetSim)
	s.stopped = true
	return s
}

func (s *RunResetSim) Connect(statusChannel, propEventChannel, requestChannel chan string) {
	// Send message back to the App/Viewer in the
	// responseLoop coroutine.
	s.statusChannel = statusChannel
	s.propEventChannel = propEventChannel
	s.requestChannel = requestChannel
}

func (s *RunResetSim) Command(args []string) {
	// fmt.Printf("Send msg: %s\n", msg)
	switch args[0] {
	case "start":
		s.stopped = false
		s.start()
	case "stop":
		s.stopped = true
		go s.respond("Stopped")
	case "ping":
		go s.respond("pong")
	case "load":
		// res := fmt.Sprintf("loading `%s`", args[1])
		go s.respond("loaded")
	case "prop":
		s.changeProperty(args[1:])
	}
}

func (s *RunResetSim) Send(msg string) {
	args := strings.Split(msg, " ")
	s.Command(args)
}

// Sends a msg back through channel async
func (s *RunResetSim) respond(msg string) {
	s.statusChannel <- msg
}

func (s *RunResetSim) respondPropEvent(msg string) {
	s.propEventChannel <- msg
}

func (s *RunResetSim) start() {
	s.Create()
	fmt.Println("Starting...")

	// Start the simulation loop in a coroutine.
	go s.run()
}

func (s *RunResetSim) Create() {
	fmt.Println("Creating...")
	s.runDuration = 1000 // 1000ms
	s.t = 0
	s.dt = 0.0

	s.sim = NewSimulation(s.statusChannel, s.propEventChannel)
	synCnt := s.sim.initialize()

	fmt.Printf("Syn cnt: %d, duration: %d\n", synCnt, s.runDuration)

	samples.PoiSamples = samples.NewDatSamples(synCnt, s.runDuration)
	samples.StimSamples = samples.NewDatSamples(synCnt, s.runDuration)
	samples.CellSamples = samples.NewNeuronSamples(s.runDuration)

	fmt.Println("Launched.")
}

// This runs in a "Go"routine.
func (s *RunResetSim) run() {
	// Run the sim for a fixed amount of time and then reset.

	for !s.stopped {
		if s.t >= s.runDuration {
			s.Reset()
		} else {
			s.Step()
		}
	}

	fmt.Println("RunReset: run() loop exited")
	s.respond("Stopped")
}

func (s *RunResetSim) Reset() {
	// Reset
	s.t = 0
	s.dt = 0.0
	// Reset random seeds.
	s.sim.reset()
}

func (s *RunResetSim) Step() {
	fmt.Printf("Step: (%d), %f\n", s.t, s.dt)
	s.sim.simulate(s.dt)
	s.t++
	s.dt += 1.0
}

func (s *RunResetSim) RunPause() {
	for s.t < s.runDuration {
		// Run
		s.sim.simulate(s.dt)
		s.t++
		s.dt += 1.0
	}
}

func (s *RunResetSim) changeProperty(args []string) {
	s.sim.changeProperty(args)
}

func (s *RunResetSim) RequestProperty(property string) string {
	return s.sim.requestProperty(property)
}

func (s *RunResetSim) SetCommand(cmd []string) {
	s.sim.SetCommand(cmd)
}
