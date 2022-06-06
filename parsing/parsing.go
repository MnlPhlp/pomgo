package parsing

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/MnlPhlp/pomgo/modes"
)

type PlanString string

type Interval struct {
	Mode int
	Text string
	Time time.Duration
}

func (s PlanString) nextChar(i int) rune {
	if len(s) >= i+2 {
		return []rune(s)[i+1]
	} else {
		return '_'
	}
}

func ParsePlan(plan PlanString) []Interval {
	parsingTime := false
	parsingText := false
	ok := false
	currentMode := 0
	timeStr := ""
	text := ""
	result := []Interval{}
	for i, c := range plan {
		// parse custom time for a mode
		if parsingTime {
			if unicode.IsNumber(c) {
				timeStr += string(c)
			} else {
				parsingTime = false
				if currentMode == modes.CUSTOM {
					parsingText = true
				} else {
					minutes, _ := strconv.Atoi(timeStr)
					result = append(result, Interval{
						Mode: currentMode,
						Text: "",
						Time: time.Duration(time.Duration(minutes) * time.Minute),
					})
					timeStr = ""
				}
			}
		}
		// parse custom text for a custom interval
		if parsingText {
			if c == ':' {
				parsingText = false
				minutes, _ := strconv.Atoi(timeStr)
				text = strings.Replace(text, "_", " ", -1)
				result = append(result, Interval{
					Mode: currentMode,
					Text: text,
					Time: time.Duration(time.Duration(minutes) * time.Minute),
				})
				text = ""
				continue
			} else {
				text += string(c)
			}
		}
		// start to parse a new mode
		if !parsingTime && !parsingText {
			if currentMode, ok = modes.ModeMap[c]; ok {
				if unicode.IsNumber(plan.nextChar(i)) {
					parsingTime = true
				} else {
					if currentMode == modes.CUSTOM {
						panic("Error: custom intervals need a specified time")
					}
					result = append(result, Interval{
						Mode: currentMode,
						Text: "",
						Time: modes.Time[currentMode],
					})
				}
			} else {
				panic(fmt.Errorf("invalid mode %v", c))
			}
		}
	}
	//check if parsing finished successfully
	if parsingTime {
		minutes, _ := strconv.Atoi(timeStr)
		result = append(result, Interval{
			Mode: currentMode,
			Text: "",
			Time: time.Duration(time.Duration(minutes) * time.Minute),
		})
	}
	if parsingText {
		panic("Error: missing ':' to end custom task")
	}
	return result
}
