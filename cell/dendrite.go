package cell

import sll "github.com/emirpasic/gods/lists/singlylinkedlist"

// IDendrite collects and manages ICompartments.
type IDendrite interface {
	AddCompartment(ICompartment)

	Integrate(t float64) float64

	Process()

	Reset()
}

type baseDendrite struct {
	// Collection of dendrite compartments.
	// Proximal, Apical and Distal
	compartments *sll.List
}

func (bc *baseDendrite) initialize() {
	bc.compartments = sll.New()

}

func (bc *baseDendrite) Reset() {
	// Reset properties

}

func (bc *baseDendrite) AddCompartment(comp ICompartment) {
	bc.compartments.Add(comp)
}
