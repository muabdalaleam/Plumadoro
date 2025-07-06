package phase

import (
	"time"
	"errors"
)

type PhaseType byte

const (
	Focus PhaseType = iota
	ShortBreak
	LongBreak
)

type Phase struct {
	remainingTime  time.Duration
	pausedTime     time.Duration
	phaseType      PhaseType
	running        bool
	n              uint8 
}

var phaseDurations = map[PhaseType]time.Duration {
	Focus:      time.Minute * 25,
	ShortBreak: time.Minute * 5,
	LongBreak:  time.Minute * 20,
}

var phaseAutostart = false

func isValidPhaseType(pt PhaseType) bool {
	switch pt {
	case Focus, ShortBreak, LongBreak:
		return true
	default:
		return false
	}
}

func SetPhaseDurations(pd map[PhaseType]time.Duration) error {
	var err error = nil

	if pd == nil {
		err = errors.New("Passed phase durations map is empty")
	}

	// HACK: I am looping over an enum values using hard coded number this is probably
	// not the best thing to do
	for i := 0; i < 3; i++ {
		if pd[PhaseType(i)] >= time.Minute * 1000 ||
		   pd[PhaseType(i)] <= time.Second * 10 {
			pd[PhaseType(i)] = phaseDurations[PhaseType(i)]
			err = errors.New("Phase durations can't exceed 1000 Minutes not be lower than 10 seconds")
		}
	}

	phaseDurations[Focus]      = pd[Focus]
	phaseDurations[ShortBreak] = pd[ShortBreak]
	phaseDurations[LongBreak]  = pd[LongBreak]

	return err
}

func SetPhaseAutostart(au bool) {
	phaseAutostart = au
}

func InitFirstPhase() (Phase) {
	return Phase {
		remainingTime: phaseDurations[Focus],
		pausedTime: time.Duration(0),
		phaseType: Focus,
		running: phaseAutostart,
		n: 1, // NOTE: the index of phases is one based
	}
}

func (p *Phase) Toggle() {
	if p.running == false {
		p.running = true
	} else {
		p.running = false
	}
}

func (p *Phase) Tick(d time.Duration) {
	if p.running {
		p.remainingTime -= d
	} else {
		p.pausedTime += d
	}

	if p.remainingTime <= time.Duration(0) {
		p.Next()
	}
}

func (p *Phase) Reset() {
	p.remainingTime = phaseDurations[p.phaseType]
	p.pausedTime    = time.Duration(0)
	p.running       = phaseAutostart
}

func (p *Phase) Next() {
	p.Reset()

	// Updating phase type
	var newPhaseType PhaseType

	// NOTE: be careful n is updated first
	p.n += 1

	switch p.phaseType {
	case ShortBreak, LongBreak:
		newPhaseType = Focus
		break

	// HACK: Checking next phase should be more configurable and not hard coded
	case Focus:
		if p.n == 7 || (p.n - 1) % 7 == 0 { // because for every 4 focus phases there are three short breaks in between
			newPhaseType = LongBreak
		} else {
			newPhaseType = ShortBreak
		}
		break
	}

	p.pausedTime    = time.Duration(0)
	p.remainingTime = phaseDurations[newPhaseType]
	p.phaseType     = newPhaseType
	p.running       = phaseAutostart
}

func (p *Phase) GetRemainingTime() time.Duration {
	return p.remainingTime
}

func (p *Phase) GetPausedTime() time.Duration {
	return p.pausedTime
}

func (p *Phase) GetPhaseType() PhaseType {
	return p.phaseType
}

func (p *Phase) IsRunning() bool {
	return p.running
}

func (p *Phase) GetN() uint8 {
	return p.n
}

func (p *Phase) GetProgress() float64 {
	return float64(p.remainingTime) /
		   float64(phaseDurations[p.phaseType])
}

