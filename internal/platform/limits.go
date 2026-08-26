package platform

type Limits struct {
	MaximumBodyBytes int64
	MaximumDepth     int
	MaximumResults   int
	MaximumBatch     int
}

func DefaultLimits() Limits {
	return Limits{MaximumBodyBytes: 8 << 20, MaximumDepth: 32, MaximumResults: 10000, MaximumBatch: 1000}
}
func (l Limits) Valid() bool {
	return l.MaximumBodyBytes > 0 && l.MaximumDepth > 0 && l.MaximumResults > 0 && l.MaximumBatch > 0
}
func (l Limits) ClampDepth(value int) int {
	if value < 1 {
		return 1
	}
	if value > l.MaximumDepth {
		return l.MaximumDepth
	}
	return value
}
func (l Limits) ClampResults(value int) int {
	if value < 1 {
		return 1
	}
	if value > l.MaximumResults {
		return l.MaximumResults
	}
	return value
}
