package entropy

// Generator produces a fresh high-intensity entropy segment on demand.
type Generator interface {
	Generate() (*Entropy, error)
}
