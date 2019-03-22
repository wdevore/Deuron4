package cell

// This synapse is for prototyping only.

type ProtoSynapse struct {
	baseSynapse

	// The time-mark at which a spike arrived at a synapse
	preT int

	// -----------------------------------
	// Depression pair-STDP
	// -----------------------------------
	// denominator, positive window time decay
	taoP float64
	// denominator, negative window time decay
	taoN float64

	// -----------------------------------
	// Potentiation triplet-STDP
	// -----------------------------------
}

func NewProtoSynapse(comp ICompartment, synType SynapseType, id int) ISynapse {
	n := new(ProtoSynapse)
	n.comp = comp
	n.synType = synType
	n.SetId(id)
	comp.AddSynapse(n)
	n.baseSynapse.initialize()
	return n
}

// Process handles post processing after Integrate has
// completed. It is considered the 1st pass of the simulation per time step.
// Internal values are 'moved' to the outputs.
// Learning rules are applied.
func (n *ProtoSynapse) Process() {
	// If this is a
}

// Integrate is the 2nd pass and handles integration.
// The effects pre/post synaptic spikes are felt here.
func (n *ProtoSynapse) Integrate(t float64) float64 {
	// TODO Integration is a behavior implemented by a Functor.

	return 0.0
}
