package cell

// Global ID auto incrementing
var gid int

func GetNextId() int {
	g := gid
	gid++
	return g
}

// ICell represents a network wide cell of which there can be many
// implementations.
// Cells make connections with other cells via [IConnection]s
type ICell interface {
	ID() int
	SetID(int)

	// Output is either a 1 or 0
	Output() byte

	// Input comes from a IConnection
	AddInConnection(IConnection)

	AddOutConnection(IConnection)

	AttachDendrite(IDendrite)

	Integrate(dt float64) float64

	Process()

	Reset()

	Diagnostics(string)
}

// IConnection represents a connection between inputs and/or cells.
// Connections transport Data objects.
//
// A connection has a collection of inputs streams and output targets.
// A connection merges streams that target multiple cells.
//
// Multiple cells can connect their output to the connection's input
// and the connection's ouput can feed into multiple cell inputs.
//
// Some connections can delay their output behind their input.
// Some connections and have no delay (aka strait).
type IConnection interface {
	// If a connection has a delay then this will
	// step the delay
	Update()

	// Inject a Data (aka spike) into the connection.
	Input(byte)

	// Get connection's output
	Output() byte

	// Post pass after Process and Integrate
	Post()
}
