package modes

import "time"

const (
	WORK = iota
	SHORT_BREAK
	LONG_BREAK
	CUSTOM
)

var ModeMap = map[rune]int{
	'w': WORK,
	's': SHORT_BREAK,
	'l': LONG_BREAK,
	'c': CUSTOM,
}

var (
	Time = []time.Duration{25 * time.Minute, 5 * time.Minute, 15 * time.Minute}
	Text = []string{"time to work", "take a short break", "take a long break", "unnamed custom Interval"}
	Desc = []string{"work", "short break", "long break", "custom"}
)
