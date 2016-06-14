package main

import "time"

// A TimeInterval represents a span of time. The start is inclusive, the end
// is exclusive.
type TimeInterval struct {
	Start, End time.Time
}

// NewTimeInterval constructs a TimeInterval value after checking for logic
// errors.
func NewTimeInterval(start, end time.Time) TimeInterval {
	if end.Before(start) {
		panic("logic error in NewTimeInterval: bad input times")
	}
	return TimeInterval{Start: start.UTC(), End: end.UTC()}
}

// Overlap detects whether this TimeInterval overlaps with another
// TimeInterval.
func (ti *TimeInterval) Overlap(other *TimeInterval) bool {
	S := ti.Start.UnixNano()
	E := ti.End.UnixNano()

	s := other.Start.UnixNano()
	e := other.End.UnixNano()

	// special case 1 of 2: when boundaries match exactly, maintain the
	// property that end is exclusive but start is inclusive:
	if E == s {
		return false
	}

	// special case 2 of 2: when boundaries match exactly, maintain the
	// property that start is inclusive but end is exclusive:
	if e == S {
		return false
	}

	// *{--[--]--}* (surrounds other)
	if S <= s && e <= E {
		return true
	}

	// *{--[--}*--] (overlaps other start)
	if S <= s && s <= E {
		return true
	}

	// [--*{--]--}* (overlaps other end)
	if S <= e && e <= E {
		return true
	}

	// *[--{--}--]* (contained within other)
	if s <= S && E <= e {
		return true
	}

	return false
}

// bucketTimeIntervals is a helper that creates a slice of TimeInterval
// over the given span of time, in chunks of duration `window`.
func bucketTimeIntervals(start, end time.Time, window time.Duration) []TimeInterval {
	if end.Before(start) {
		panic("logic error in bucketTimeIntervals: bad input times")
	}
	ret := []TimeInterval{}

	start = start.Truncate(window)
	for start.Before(end) {
		ti := NewTimeInterval(start, start.Add(window))
		ret = append(ret, ti)
		start = start.Add(window)
	}

	return ret
}
