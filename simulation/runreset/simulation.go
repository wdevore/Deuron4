package runreset

import (
	"fmt"
	"math/rand"
	"strconv"

	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/wdevore/Deuron4/cell"
	"github.com/wdevore/Deuron4/cell/stimulus"
	"github.com/wdevore/Deuron4/simulation/samples"
)

type simulation struct {
	channel          chan string
	propEventChannel chan string

	neuron cell.ICell

	poiStreams  *sll.List
	stimStreams *sll.List
	syns        *sll.List
	cons        *sll.List

	lastCmd []string

	pattern1 *stimulus.PoissonPatternStream
}

func NewSimulation(channel, propEventChannel chan string) *simulation {
	s := new(simulation)
	s.channel = channel
	s.propEventChannel = propEventChannel
	return s
}

var ran = rand.New(rand.NewSource(1963))

func (s *simulation) initialize() int {
	s.neuron = cell.NewProtoNeuron()

	// A neuron has a dendrite
	den := cell.NewProtoDendrite(s.neuron)

	comp := cell.NewProtoCompartment(den)

	// Create 80% Excite and 20% Inhibit
	synCount := 10
	excite := int(float64(synCount) * 0.8)
	inhibit := int(float64(synCount) * 0.2)

	s.poiStreams = sll.New()
	s.stimStreams = sll.New()
	s.syns = sll.New()
	s.cons = sll.New()

	synId := 0
	poiId := 0

	s.createPatterns()
	s.pattern1.Begin()

	// For each synapse we attach a connection.
	// For this simulation each connection is also connected to
	// a poisson and pattern stream.
	for i := 0; i < excite; i++ {
		syn := cell.NewProtoSynapse(comp, cell.Excititory, synId)
		s.syns.Add(syn)

		con := cell.NewStraightConnection()
		s.cons.Add(con)

		seed := ran.Int63()
		poi := stimulus.NewPoissonStream(seed).(*stimulus.PoissonStream)
		poi.SetId(poiId)

		// Collect streams so we can step() it later.
		s.poiStreams.Add(poi)

		// Connect the pieces together:
		// stream -> connection -> synapse

		poi.Attach(con) // route noise stream into connection

		stim := s.pattern1.Stream()
		s.stimStreams.Add(stim)
		stim.Attach(con) // route stimulus into connection

		syn.Connect(con) // route connection to synapse

		synId++
		poiId++
	}

	for i := 0; i < inhibit; i++ {
		syn := cell.NewProtoSynapse(comp, cell.Inhibitory, synId)
		s.syns.Add(syn)

		con := cell.NewStraightConnection()
		s.cons.Add(con)

		seed := ran.Int63()
		poi := stimulus.NewPoissonStream(seed).(stimulus.IPatternStream)
		poi.SetId(poiId)

		s.poiStreams.Add(poi)
		// Connect stream to input of connection
		poi.Attach(con)

		stim := s.pattern1.Stream()
		s.stimStreams.Add(stim)
		stim.Attach(con) // route stimulus into connection

		syn.Connect(con) // attach connection into synapse

		synId++
		poiId++
	}

	s.neuron.AttachDendrite(den)

	fmt.Println("Sim: initialized")

	return synCount
}

func (s *simulation) reset() {
	it := s.poiStreams.Iterator()
	for it.Next() {
		poi := it.Value().(stimulus.IPatternStream)
		poi.Reset()
	}

	// Reset stimulus
	s.pattern1.Reset()

}

// A single pass of a simulation.
func (s *simulation) simulate(t float64) {
	// fmt.Printf("Pass: %f\n", t)
	s.pre()

	s.diagnostics(t) // Collect samples

	// Update learning rules (STDP and BTSP) and internal states/properties
	s.neuron.Process()

	// Now integrate
	epsp := s.neuron.Integrate(t)

	s.post()

	// Update app state.
	msg := fmt.Sprintf("Running (%d) epsp:(%f)...", int(t), epsp)

	// Update the app thread with a message.
	s.respond(msg)

	// time.Sleep(time.Millisecond * 100)
}

func (s *simulation) pre() {
	// Prep: Update streams first
	it := s.poiStreams.Iterator()
	for it.Next() {
		poi := it.Value().(stimulus.IPatternStream)
		poi.Step()
	}

	// Step all the stimulus streams
	s.pattern1.Step()
}

func (s *simulation) diagnostics(t float64) {
	// Capture the state at time "t".

	// Collect noise samples from the poisson streams.
	it := s.poiStreams.Iterator()
	for it.Next() {
		pois := it.Value().(stimulus.IPatternStream)
		samples.PoiSamples.Put(t, pois.Output(), pois.Id(), 3)
	}

	if s.pattern1.Begin() {
		more := true
		for more {
			stim := s.pattern1.Stream()
			samples.StimSamples.Put(t, stim.Output(), stim.Id(), 4)
			more = s.pattern1.Next()
		}
	}

	// Capture the cell's current output
	samples.CellSamples.Put(t, s.neuron.Output(), s.neuron.ID(), 0)
}

func (s *simulation) post() {
	// Post is a preperation for next pass.
	// This means we put all Data, on the output side of a connection,
	// back into the pool.

	// The values either source from noise streams, stimulus or other neuron outputs.
	it := s.cons.Iterator()
	for it.Next() {
		con := it.Value().(cell.IConnection)
		con.Post()
	}
}

func (s *simulation) respond(msg string) {
	// Send message back to the App
	s.channel <- msg
}

func (s *simulation) propertyChangeEvent(msg string) {
	// Send message back to the App
	s.propEventChannel <- msg
}

func (s *simulation) requestProperty(property string) string {
	switch property {
	case "Poisson Max":
		it := s.poiStreams.Iterator()
		if it.Next() {
			poi := it.Value().(*stimulus.PoissonStream)
			return fmt.Sprintf("%f", poi.Max())
		}
		break
	case "Poisson Min":
		it := s.poiStreams.Iterator()
		if it.Next() {
			poi := it.Value().(*stimulus.PoissonStream)
			return fmt.Sprintf("%f", poi.Min())
		}
		break
	case "Poisson Spread":
		it := s.poiStreams.Iterator()
		if it.Next() {
			poi := it.Value().(*stimulus.PoissonStream)
			return fmt.Sprintf("%f", poi.Spread())
		}
		break
	}

	return ""
}

func (s *simulation) SetCommand(cmd []string) {
	s.lastCmd = []string{}
	for _, sa := range cmd {
		s.lastCmd = append(s.lastCmd, sa)
	}
}

func (s *simulation) changeProperty(args []string) {
	// fmt.Printf("changeProperty: %v\n", args)

	// Could be "up" or "down" with inc value
	thing := args[0]

	if thing == "up" || thing == "down" {
		// change value of last property
		// fmt.Printf("lascmd: %v\n", s.lastCmd)
		if s.lastCmd == nil {
			return
		}

		lastValue, _ := strconv.ParseFloat(s.lastCmd[2], 64)
		value, err2 := strconv.ParseFloat(args[1], 64)
		if err2 != nil {
			fmt.Println("RunReset:changeProperty up/down unrecognized.")
			return
		}

		if thing == "up" {
			lastValue += value
		} else {
			lastValue -= value
		}

		thing = s.lastCmd[0]
		args = []string{thing, s.lastCmd[1], fmt.Sprintf("%f", lastValue)}
		// fmt.Printf("up/down: %v\n", args)
	}

	switch thing {
	case "Poisson":
		property := args[1]
		value, err := strconv.ParseFloat(args[2], 64)
		s.SetCommand(args)

		if err != nil {
			fmt.Printf("RunReset:changeProperty command properties correct: %s\n", args[3])
			return
		}

		switch property {
		case "Max":
			it := s.poiStreams.Iterator()
			for it.Next() {
				poi := it.Value().(*stimulus.PoissonStream)
				poi.SetMax(value)
			}

			s.propertyChangeEvent("Poisson Max," + args[2])
		case "Min":
			it := s.poiStreams.Iterator()
			for it.Next() {
				poi := it.Value().(*stimulus.PoissonStream)
				poi.SetMin(value)
			}
			s.propertyChangeEvent("Poisson Min," + args[2])
		case "Spread":
			it := s.poiStreams.Iterator()
			for it.Next() {
				poi := it.Value().(*stimulus.PoissonStream)
				poi.SetSpread(value)
			}
			s.propertyChangeEvent("Poisson Spread," + args[2])
		}
	}
}

func (s *simulation) createPatterns() {
	// ------------------------------------------------------------
	// Create collection
	s.pattern1 = stimulus.NewPoissonPatternStream(123)
	// s.pattern1.Period(100, 25) // Pattern will be applied at 30Hz or every 33ms

	// Create patterns
	spk := stimulus.NewSpikeStream().(*stimulus.SpikeStream)
	spk.SetId(0)
	spk.SetSpikes([]byte{0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0})
	s.pattern1.Add(spk)

	spk = stimulus.NewSpikeStream().(*stimulus.SpikeStream)
	spk.SetId(1)
	spk.SetSpikes([]byte{0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0})
	// spk.SetSpikes([]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1})
	s.pattern1.Add(spk)

	spk = stimulus.NewSpikeStream().(*stimulus.SpikeStream)
	spk.SetId(2)
	spk.SetSpikes([]byte{1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1})
	s.pattern1.Add(spk)

	spk = stimulus.NewSpikeStream().(*stimulus.SpikeStream)
	spk.SetId(3)
	spk.SetSpikes([]byte{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0})
	s.pattern1.Add(spk)

	spk = stimulus.NewSpikeStream().(*stimulus.SpikeStream)
	spk.SetId(4)
	spk.SetSpikes([]byte{0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0})
	s.pattern1.Add(spk)

	spk = stimulus.NewSpikeStream().(*stimulus.SpikeStream)
	spk.SetId(5)
	spk.SetSpikes([]byte{0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0})
	s.pattern1.Add(spk)

	spk = stimulus.NewSpikeStream().(*stimulus.SpikeStream)
	spk.SetId(6)
	spk.SetSpikes([]byte{1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0})
	s.pattern1.Add(spk)

	spk = stimulus.NewSpikeStream().(*stimulus.SpikeStream)
	spk.SetId(7)
	spk.SetSpikes([]byte{0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1})
	s.pattern1.Add(spk)

	spk = stimulus.NewSpikeStream().(*stimulus.SpikeStream)
	spk.SetId(8)
	spk.SetSpikes([]byte{0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1})
	s.pattern1.Add(spk)

	spk = stimulus.NewSpikeStream().(*stimulus.SpikeStream)
	spk.SetId(9)
	spk.SetSpikes([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1})
	// spk.SetSpikes([]byte{1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1})
	s.pattern1.Add(spk)

	// fmt.Printf("createPatterns: \n%s\n", s.pattern1)
}
