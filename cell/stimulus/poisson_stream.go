package stimulus

import (
	"math"
	"math/rand"

	"github.com/wdevore/Deuron4/cell"
)

func RanGen(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

// PoissonStream generates spikes with poisson distribution.
// The outputs are generally routed into StraitConnections.
type PoissonStream struct {
	basePatternStream

	ran *rand.Rand

	// Random seed
	seed int64

	// The Inter spike interval (ISI) counter is populated by a value.
	// When the counter reaches 0 a spike is placed on the output.
	isi int

	// Poisson properties
	max    float64
	spread float64
	min    float64
}

// NewPoissonStream creates a stream
func NewPoissonStream(seed int64) IPatternStream {
	s := new(PoissonStream)
	s.baseInitialize()

	s.seed = seed
	s.ran = RanGen(seed)

	s.max = 300.0
	s.spread = 50.0
	s.min = 7.0

	s.isi = s.generate(s.max, s.spread, s.min)
	return s
}

func (ss *PoissonStream) Initialize(max, spread, min float64) {
	ss.max = max
	ss.spread = spread
	ss.min = min
	ss.Reset()
}

// Attach connections to this stream
// The given IConnection will have spikes routed into it.
func (ss *PoissonStream) Attach(con cell.IConnection) {
	ss.cons.Add(con)
}

func (ss *PoissonStream) ISI() int {
	return ss.isi
}

func (ss *PoissonStream) SetMax(v float64) {
	ss.max = v
}

func (ss *PoissonStream) Max() float64 {
	return ss.max
}

func (ss *PoissonStream) Min() float64 {
	return ss.min
}

func (ss *PoissonStream) Spread() float64 {
	return ss.spread
}

func (ss *PoissonStream) SetMin(v float64) {
	ss.min = v
}

func (ss *PoissonStream) SetSpread(v float64) {
	ss.spread = v
}

func Generate(rand, scale, div, min float64) int {
	return int(scale*math.Pow(math.E, -rand*scale/div) + min)
}

// generate tends to spread spikes a bit more.
// Typical values of: 15.0, 3.0 yield ISIs 5-7 with occasional 50-100s,
// or 50.0,15.0,2.0
func (ss *PoissonStream) generate(scale, div, min float64) int {
	return Generate(ss.ran.Float64(), scale, div, min)
}

// ----------------------------------------------
// IPatternStream methods
// ----------------------------------------------

func (ss *PoissonStream) EnableAutoReset() {
	// Not applicable
}

// Reset generates a new ISI
func (ss *PoissonStream) Reset() {
	ss.ran.Seed(ss.seed)
	ss.isi = ss.generate(ss.max, ss.spread, ss.min)
}

func (ss *PoissonStream) Step() bool {

	// Check ISI counter
	if ss.isi == 0 {
		// Time to generate a spike
		ss.value = 1
		ss.isi = ss.generate(ss.max, ss.spread, ss.min)
	} else {
		ss.value = 0
		ss.isi--
	}

	// Place stream's current output value onto the
	// associated connection(s) input
	it := ss.cons.Iterator()
	for it.Next() {
		conn := it.Value().(cell.IConnection)
		conn.Input(ss.value)
	}

	return false
}

func (ss *PoissonStream) IsComplete() bool {
	return false // This type of stream never completes
}

// ----------------------------------------------
// IBitStream methods
// ----------------------------------------------

func (ss *PoissonStream) Input(v byte) {
	// Not applicable.
}

func (ss *PoissonStream) Output() byte {
	return ss.value
}
