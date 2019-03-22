package samples

// Spike is used for collecting samples.
type Spike struct {
	Time  float64
	Value byte
	// This 'key' can represent anything, for example what color to
	// render the spike on a graph
	Key int

	// What the spike belongs to.
	Id int
}

func NewSpike() *Spike {
	s := new(Spike)
	return s
}
