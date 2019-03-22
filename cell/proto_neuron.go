package cell

// This neuron is for prototyping only.

type ProtoNeuron struct {
	baseCell

	// Soma threshold. When exceeded an AP is generated.
	threshold float64

	// --------------------------------------------------------
	// Action potential
	// --------------------------------------------------------
	// AP can travel back down the dendrite. The value decays
	// with distance.
	apDecay float64

	// The time-mark of the current AP.
	APt int
	// The previous time-mark of an AP
	preAPt int

	// The maximum value AP spikes to. This can be controlled by
	// meta-plasticity.
	maxAP float64

	// --------------------------------------------------------
	// STDP
	// --------------------------------------------------------
	// potentiation time-constant (decay)
	taoP float64 // both for pair and triplet

	// -----------------------------------
	// Depression pair-STDP
	// -----------------------------------
	// depression time-constant (decay)
	taoN float64

	// -----------------------------------
	// Potentiation triplet-STDP
	// -----------------------------------
	taoY float64
}

func NewProtoNeuron() ICell {
	n := new(ProtoNeuron)
	n.baseCell.initialize()
	return n
}

func (n *ProtoNeuron) Output() byte {
	return n.output
}

func (n *ProtoNeuron) AddInConnection(con IConnection) {
	n.inputs = append(n.inputs, con)
}

func (n *ProtoNeuron) AddOutConnection(con IConnection) {
	n.outputs = append(n.outputs, con)
}

func (n *ProtoNeuron) Integrate(t float64) float64 {
	return n.dendrite.Integrate(t)
}

func (n *ProtoNeuron) Process() {
	n.dendrite.Process()
}

func (n *ProtoNeuron) Reset() {

}
