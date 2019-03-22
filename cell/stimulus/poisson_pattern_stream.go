package stimulus

import (
	"fmt"
	"math/rand"
	"strings"

	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
)

// PoissonPatternStream emits a pattern inbetween ISI windows.
// When the pattern has completed emission a new ISI is generated
// as a delay.
type PoissonPatternStream struct {
	output byte

	ran  *rand.Rand
	seed int64

	// Poisson properties
	max    float64
	spread float64
	min    float64

	autoReset bool
	// When a pattern completes generate a new ISI automatically.
	autoGenerateISI bool

	// A collection of streams
	patterns *sll.List
	patItr   sll.Iterator
	// isi = spike intervals
	isi int // in milliseconds

	delayCnt int
}

func NewPoissonPatternStream(seed int64) *PoissonPatternStream {
	s := new(PoissonPatternStream)
	s.autoReset = true
	s.seed = seed
	s.ran = RanGen(seed)

	s.max = 300.0
	s.spread = 50.0
	s.min = 50.0

	s.isi = Generate(s.ran.Float64(), s.max, s.spread, s.min)
	fmt.Printf("isi: %d\n", s.isi)

	s.patterns = sll.New()
	return s
}

// SetISI sets the "ISI" between pattern applications. The pattern
// period itself could be longer than the interval.
func (nps *PoissonPatternStream) SetISI(isi int) {
	nps.isi = isi
}

func (nps *PoissonPatternStream) Add(strm IPatternStream) {
	nps.patterns.Add(strm)
}

func (nps *PoissonPatternStream) IsComplete() bool {
	return false
}

func (nps *PoissonPatternStream) EnableAutoReset() {
	nps.autoReset = true
}

func (nps *PoissonPatternStream) Reset() {
	fmt.Println("--------------- POI pattern RESETing")
	nps.ran.Seed(nps.seed)
	nps.patternReset()
}

func (nps *PoissonPatternStream) patternReset() {
	nps.delayCnt = 0
	nps.isi = Generate(nps.ran.Float64(), nps.max, nps.spread, nps.min)
	fmt.Printf("isi: %d\n", nps.isi)
	it := nps.patterns.Iterator()
	for it.Next() {
		stim := it.Value().(IPatternStream)
		stim.Reset()
	}
}

func (nps *PoissonPatternStream) Step() {
	// Step all the streams when the ISI had ended.
	// Once the pattern has completed we switch back to ISI.
	if nps.delayCnt > nps.isi {
		var complete bool
		it := nps.patterns.Iterator()
		for it.Next() {
			stim := it.Value().(IPatternStream)
			complete = complete || stim.Step()
		}

		if complete {
			nps.patternReset()
		}
	} else {
		nps.delayCnt++
	}
}

func (nps *PoissonPatternStream) Begin() bool {
	nps.patItr = nps.patterns.Iterator()
	return nps.patItr.First()
}

func (nps *PoissonPatternStream) Next() bool {
	return nps.patItr.Next()
}

// Stream returns the next available pattern stream.Begin
func (nps *PoissonPatternStream) Stream() IPatternStream {
	stream := nps.patItr.Value().(IPatternStream)
	return stream
}

func (nps PoissonPatternStream) String() string {
	var s strings.Builder

	it := nps.patterns.Iterator()
	for it.Next() {
		stim := it.Value().(*SpikeStream)
		s.WriteString(fmt.Sprintf("%s\n", stim.String()))
	}

	return s.String()
}
