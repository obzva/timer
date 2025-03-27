package timer

import "time"

type State int

const (
	stateIdle State = iota
	stateActive
	stateExpired
)

// The Timer type represents a single event. When the Timer expires,
// the current time will be sent on C, unless the Timer was created by AfterFunc.
// A Timer must be created with NewTimer or AfterFunc.
type Timer struct {
	C         <-chan time.Time
	c         chan time.Time
	Duration  time.Duration
	state     State
	fn        func()
	StartedAt time.Time
	t         *time.Timer
}

// AfterFunc waits after calling its Start method for the duration
// to elapse and then calls f in its own goroutine.
// It returns a Timer that can be used to cancel the call using its Stop method,
// or pause using its Pause method
func AfterFunc(d time.Duration, f func()) *Timer {
	t := NewTimer(d)
	t.fn = func() {
		t.state = stateExpired
		t.c <- time.Now()
		f()
	}
	return t
}

// NewTimer creates a new Timer.
// It returns a Timer that can be used to cancel the call using its Stop method,
// or pause using its Pause method
func NewTimer(d time.Duration) *Timer {
	c := make(chan time.Time, 1)
	t := new(Timer)
	t.C = c
	t.c = c
	t.Duration = d
	t.fn = func() {
		t.state = stateExpired
		t.c <- time.Now()
	}
	return t
}

// Pause pauses current timer until Start method will be called.
// Next Start call will wait rest of duration.
func (t *Timer) Pause() bool {
	if t.state != stateActive {
		return false
	}
	if !t.t.Stop() {
		t.state = stateExpired
		return false
	}
	t.state = stateIdle
	dur := time.Since(t.StartedAt)
	t.Duration = t.Duration - dur
	return true
}

// Start starts Timer that will send the current time on its channel after at least duration d.
func (t *Timer) Start() bool {
	if t.state != stateIdle {
		return false
	}
	t.StartedAt = time.Now()
	t.state = stateActive
	t.t = time.AfterFunc(t.Duration, t.fn)
	return true
}

// Stop prevents the Timer from firing. It returns true if the call stops the timer,
// false if the timer has already expired or been stopped.
// Stop does not close the channel, to prevent a read from the channel succeeding incorrectly.
func (t *Timer) Stop() bool {
	if t.state != stateActive {
		return false
	}
	t.StartedAt = time.Now()
	t.state = stateExpired
	t.t.Stop()
	return true
}

// Restart restarts the timer with a new duration.
func (t *Timer) Restart(d time.Duration) {
	t.Duration = d
	t.StartedAt = time.Now()
	t.state = stateActive
	t.t = time.AfterFunc(t.Duration, t.fn)
}
