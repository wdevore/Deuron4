package cell

import sll "github.com/emirpasic/gods/lists/singlylinkedlist"

// ICompartment collects synapses. It represents the functionality
// of a group of synapses located along the Dendrite.
// Compartments are either close to or farther away from the neuron's
// soma.
// TODO compartments can overlap causing effects to diffuse into neighboring
// compartments, for example, Ca dynamics.
type ICompartment interface {
	// Compartment properties
	AddSynapse(ISynapse)

	// Behaviors

	// Evaluates the total effective weight for the compartment.
	Integrate(t float64) float64

	Process()

	Reset()
}

type baseCompartment struct {
	// Collection of synapses
	synapses *sll.List
}

func (bc *baseCompartment) initialize() {
	bc.synapses = sll.New()
}

func (bc *baseCompartment) Reset() {
	// Reset properties

}

func (bc *baseCompartment) AddSynapse(syn ISynapse) {
	bc.synapses.Add(syn)
}
