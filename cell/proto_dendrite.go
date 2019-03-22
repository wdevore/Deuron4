package cell

type ProtoDendrite struct {
	baseDendrite

	neuron ICell
}

func NewProtoDendrite(cell ICell) IDendrite {
	n := new(ProtoDendrite)

	// Bidirectional associations
	n.neuron = cell
	n.baseDendrite.initialize()

	return n
}

// 1st pass

// Process handles post processing before Integration is performed.
func (d *ProtoDendrite) Process() {
	it := d.compartments.Iterator()
	for it.Next() {
		comp := it.Value().(ICompartment)
		comp.Process()
	}
}

// 2nd pass

// Integrate is the 2nd pass performing integration.
func (d *ProtoDendrite) Integrate(t float64) float64 {
	// TODO Integration is a behavior implemented by a Functor.

	w := 0.0

	it := d.compartments.Iterator()
	for it.Next() {
		comp := it.Value().(ICompartment)
		w += comp.Integrate(t)
	}

	return w
}
