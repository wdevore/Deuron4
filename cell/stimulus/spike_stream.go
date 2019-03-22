package stimulus

import (
	"fmt"
	"strings"

	"github.com/wdevore/Deuron4/cell"
)

// SpikeStream provides a spiking stimulus stream
// For a pattern to form spatially you will need 2 or more streams
//
// This type of stream outputs a pattern at a certain frequency.
// For example, at 60Hz the pattern is presented every 17ms which means
// the pattern's width needs to be at least < than the period.
//
// ---------------------|--pattern--|----...----|--pattern--|---...-----|--pattern--|
// ^                                ^
// |  period = 17ms                 |
//
type SpikeStream struct {
	basePatternStream

	complete  bool
	autoReset bool

	// Spike pattern. The pattern is fixed in size, for now.
	pattern []byte
	idx     int
}

// [size] is in milliseconds
func NewSpikeStream() IPatternStream {
	s := new(SpikeStream)
	s.baseInitialize()

	s.autoReset = false

	s.Reset()
	s.EnableAutoReset()

	return s
}

func (ss *SpikeStream) Input(v byte) {
}

func (ss *SpikeStream) Output() byte {
	return ss.value
}

func (ss *SpikeStream) IsComplete() bool {
	return ss.complete
}

func (ss *SpikeStream) EnableAutoReset() {
	ss.autoReset = true
}

func (ss *SpikeStream) Reset() {
	ss.complete = false
	ss.idx = len(ss.pattern) - 1
	ss.value = 0
}

func (ss *SpikeStream) Step() bool {
	// Only step the pattern after the delay.

	// Step the pattern from end to start so the
	// pattern in code looks the same on the display.
	if ss.autoReset && ss.idx < 0 {
		ss.Reset()
		return true // complete
	}

	ss.value = ss.pattern[ss.idx]

	// Place stream's current output value onto the
	// associated connection(s) input
	it := ss.cons.Iterator()
	for it.Next() {
		conn := it.Value().(cell.IConnection)
		conn.Input(ss.value)
	}

	ss.idx--

	return false // not complete yet
}

func (ss *SpikeStream) Set(t int) {
	if t > len(ss.pattern) {
		fmt.Println("SpikeStream: bad t position")
		return
	}

	ss.pattern[t] = 1
}

func (ss *SpikeStream) SetRange(ts []int) {
	for _, ti := range ts {

		if ti > len(ss.pattern) {
			fmt.Println("SpikeStream: bad t position")
			return
		}

		ss.pattern[ti] = 1
	}
}

func (ss *SpikeStream) SetSpikes(sp []byte) {
	ss.pattern = make([]byte, len(sp))

	for t, spik := range sp {

		if t > len(ss.pattern) {
			fmt.Println("SpikeStream: bad set position")
			return
		}

		ss.pattern[t] = spik
	}
	ss.Reset()
}

func (ss *SpikeStream) Clear(t int) {
	if t > len(ss.pattern) {
		fmt.Println("SpikeStream: bad clear position")
		return
	}
	ss.pattern[t] = 0
}

func (ss SpikeStream) String() string {
	var s strings.Builder

	for j := len(ss.pattern) - 1; j >= 0; j-- {
		s.WriteString(fmt.Sprintf("%d", ss.pattern[j]))
	}

	return s.String()
}
