package timestamp

// Clock supplies the current wall-clock instant used to stamp a token at generation time.
type Clock interface {
	NowUnixSeconds() uint64
}
