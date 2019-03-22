package samples

import (
	"fmt"

	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
)

// Samples is a fixed size array scaled to the length of the simulation.
// The graph scans and renders the array.
// The sim populates the array based on a moving index.

// Samples is a 2D list of samples for synapses
type Samples struct {
	// List of SamplesLane(s), for all passes.
	// This data is used for rendering by the graphs.
	// Each lane is a fixed size array.
	lanes *sll.List

	// scanIdx int
	synCnt int
	size   int // typically the length of simulation

	// deprecated
	// mutex *sync.Mutex
}

// Lanes are trains of spikes for a given synapse.
type SamplesLane struct {
	Id      int
	Samples []*Spike
}

func NewSamples(synCnt, size int) *Samples {
	s := new(Samples)
	s.lanes = sll.New()
	s.synCnt = synCnt
	s.size = size

	// Pre expand collection
	for i := 0; i < synCnt; i++ {
		l := new(SamplesLane)
		l.Id = i
		l.Samples = make([]*Spike, size)
		for s := 0; s < size; s++ {
			l.Samples[s] = NewSpike()
		}
		s.lanes.Add(l)
	}

	// s.mutex = &sync.Mutex{}
	return s
}

func (s *Samples) GetLanes() *sll.List {
	return s.lanes
}

func (s *Samples) Size() int {
	return s.size
}

func (s *Samples) Put(time float64, value byte, sid, key int) {
	// sid is usually synId
	_, lif := s.lanes.Find(func(id int, v interface{}) bool {
		return v.(*SamplesLane).Id == sid
	})

	l := lif.(*SamplesLane)
	sp := l.Samples[int(time)]
	sp.Time = time
	sp.Value = value
	sp.Id = sid
	sp.Key = key

	// s.scanIdx = (s.scanIdx + 1) % s.size
}

func (s *Samples) Print() {
	it := s.lanes.Iterator()
	for it.Next() {
		lane := it.Value().(*SamplesLane)
		fmt.Printf("(%d) %v\n", lane.Id, lane.Samples)
	}
}

// ---------------------------------------------------------
// Data samples
// ---------------------------------------------------------
var PoiSamples *DatSamples
var StimSamples *DatSamples

type DatSamples struct {
	// List of SamplesLane(s), for all passes.
	// This data is used for rendering by the graphs.
	// Each lane is a fixed size array.
	lanes *sll.List

	synCnt int
	size   int // typically the length of simulation
}

func NewDatSamples(synCnt, size int) *DatSamples {
	s := new(DatSamples)
	s.lanes = sll.New()
	s.synCnt = synCnt
	s.size = size

	// Pre expand collection
	for i := 0; i < synCnt; i++ {
		l := new(SamplesLane)
		l.Id = i
		l.Samples = make([]*Spike, size)
		for s := 0; s < size; s++ {
			l.Samples[s] = NewSpike()
		}
		s.lanes.Add(l)
	}

	return s
}

func (s *DatSamples) GetLanes() *sll.List {
	return s.lanes
}

func (s *DatSamples) Size() int {
	return s.size
}

func (s *DatSamples) Put(time float64, value byte, sid, key int) {
	// sid is usually synId
	_, lif := s.lanes.Find(func(id int, v interface{}) bool {
		return v.(*SamplesLane).Id == sid
	})

	l := lif.(*SamplesLane)
	sp := l.Samples[int(time)]
	sp.Time = time
	sp.Value = value
	sp.Id = sid
	sp.Key = key
}

// ---------------------------------------------------------
// Neuron samples
// ---------------------------------------------------------
var CellSamples *NeuronSamples

type NeuronSamples struct {
	Samples []*Spike
}

func NewNeuronSamples(size int) *NeuronSamples {
	ns := new(NeuronSamples)

	// Pre expand collection
	ns.Samples = make([]*Spike, size)
	for s := 0; s < size; s++ {
		ns.Samples[s] = NewSpike()
	}
	return ns
}

func (ns *NeuronSamples) Put(time float64, value byte, sid, key int) {
	sp := ns.Samples[int(time)]
	sp.Time = time
	sp.Value = value
	sp.Id = sid
	sp.Key = key
}

// func (s *Samples) Use() *sll.List {
// 	s.mutex.Lock()
// 	return s.lanes
// }

// func (s *Samples) Release() {
// 	s.mutex.Unlock()
// }
