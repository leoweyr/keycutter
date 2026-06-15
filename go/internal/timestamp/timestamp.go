package timestamp

// Timestamp is the optional creation marker embedded as the fourth prefix component.
// It pairs a Unix epoch second count with its Base36 rendering, whose digits remain
// inside the permitted prefix character set.
type Timestamp struct {
	seconds uint64
	encoded string
}

// NewTimestamp pairs a Unix epoch second count with its Base36 rendering.
func NewTimestamp(seconds uint64, encoded string) *Timestamp {
	return &Timestamp{
		seconds: seconds,
		encoded: encoded,
	}
}

// Seconds returns the embedded Unix epoch second count.
func (timestamp *Timestamp) Seconds() uint64 {
	return timestamp.seconds
}

// Value returns the Base36 rendering placed inside the token prefix.
func (timestamp *Timestamp) Value() string {
	return timestamp.encoded
}

// String returns the Base36 rendering placed inside the token prefix.
func (timestamp *Timestamp) String() string {
	return timestamp.encoded
}
