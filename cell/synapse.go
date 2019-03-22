package cell

// Weight behavior

// ISynapse is a common interface.
// A synapse is associated with a connection and one or more behaviors.
// One such behavior is STDP for evaluating the Long-term effect.
// Each synapse is part of a group called a Compartment.
type ISynapse interface {
	Id() int

	Connect(IConnection)
	GetConnection() IConnection

	// Evaluates the total effective weight for the synapse.
	Integrate(dt float64) float64

	Process()

	// The current input on the synapse. The input is feed from an IConnection's
	// output.
	Input() byte

	IsExcititory() bool
}

type SynapseType bool

const (
	Excititory SynapseType = true
	Inhibitory             = false
)

type baseSynapse struct {
	id int

	synType SynapseType

	// initial intrinsic weight learned during the
	// prenatal period
	wI float64

	// TODO review/deprecate???
	// potentiation or depression (decays by lamba)
	// Over the long term this value is tranferred to wI
	// via consolidation.
	wP float64

	// A synapse will read this input on each integration pass
	conn IConnection

	// The compartment this synaspe resides in.
	comp ICompartment
}

func (bs *baseSynapse) initialize() {
}

func (bs *baseSynapse) IsExcititory() bool {
	return bs.synType == Excititory
}

func (bs *baseSynapse) Id() int {
	return bs.id
}

func (bs *baseSynapse) SetId(id int) {
	bs.id = id
}

func (bs *baseSynapse) Connect(con IConnection) {
	bs.conn = con
}

func (bs *baseSynapse) GetConnection() IConnection {
	return bs.conn
}

// Input returns the input feeding into this synapse.
// The synapse has an IConnection for input. This method returns
// that connection's output.
// i.e. the data entering the synapse from the connection.
func (bs *baseSynapse) Input() byte {
	return bs.conn.Output()
}
