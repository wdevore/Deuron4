package cell

type ProtoCompartment struct {
	baseCompartment

	den IDendrite
}

func NewProtoCompartment(den IDendrite) ICompartment {
	n := new(ProtoCompartment)

	// Bidirectional associations
	n.den = den
	n.baseCompartment.initialize()

	return n
}

func (c *ProtoCompartment) Process() {
	it := c.synapses.Iterator()
	for it.Next() {
		synapse := it.Value().(ISynapse)
		synapse.Process()
	}
}

func (c *ProtoCompartment) Integrate(t float64) float64 {
	// TODO Integration is a behavior implemented by a Functor.

	w := 0.0

	it := c.synapses.Iterator()
	for it.Next() {
		synapse := it.Value().(ISynapse)
		w += synapse.Integrate(t)
	}

	return w
}
