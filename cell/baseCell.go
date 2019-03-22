package cell

type baseCell struct {
	id int

	output byte // 1 or 0

	// A cell will read its input connections on each pass
	inputs  []IConnection
	outputs []IConnection

	dendrite IDendrite
}

func (bc *baseCell) initialize() {
	bc.inputs = []IConnection{}
	bc.outputs = []IConnection{}
}

func (bc *baseCell) AttachDendrite(den IDendrite) {
	bc.dendrite = den
}

func (bc *baseCell) Diagnostics(msg string) {
}

func (bc *baseCell) ID() int {
	return bc.id
}

func (bc *baseCell) SetID(id int) {
	bc.id = id
}
