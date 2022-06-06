package main

import (
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/cheggaaa/pb/v3"
)

const (
	WORK = iota
	SHORT_BREAK
	LONG_BREAK
	CUSTOM
	useNotifications = true
	colorReset       = "\033[0m"
	colorGreen       = "\033[32m"
)

var (
	modeTime = []time.Duration{25 * time.Minute, 5 * time.Minute, 15 * time.Minute}
	modeText = []string{"time to work", "take a short break", "take a long break", "unnamed custom Interval"}
	modeDesc = []string{"work", "short break", "long break", "custom"}
)

var modes = map[rune]int{
	'w': WORK,
	's': SHORT_BREAK,
	'l': LONG_BREAK,
	'c': CUSTOM,
}

//go:embed help.txt
var help string

func remTimeStr(rem time.Duration) string {
	min := int(rem.Minutes())
	sec := int(rem.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d remaining", min, sec)
}

func notify(text string) {}

func runPart(runTime time.Duration) {
	remTime := runTime
	seconds := int(runTime.Seconds())

	tmpl := `{{ bar . "<" "#" "#" "." ">"}} {{percent . "%3.f%%"}} {{string . "remaining" | green}}`
	var bar = pb.ProgressBarTemplate(tmpl).Start(seconds)
	bar.Set("remaining", remTimeStr(remTime))
	bar.SetMaxWidth(100)

	for i := 0; i < seconds; i++ {
		time.Sleep(time.Second)
		bar.Increment()
		remTime -= time.Second
		bar.Set("remaining", remTimeStr(remTime))
	}

	bar.Finish()
	fmt.Println()

	if useNotifications {
		/*import notify
		  proc notify(text: string) =
		    var n: Notification = newNotification("pomTimer", text, "dialog-information")
		    n.timeout = 10000
		    discard n.show()
		*/
	}
}

type Interval struct {
	mode int
	text string
	time time.Duration
}

func printTime(t time.Duration) {
	min := int(t.Minutes()) % 60
	hour := int(t.Hours())
	hourText := ""
	if hour > 0 {
		hourText = fmt.Sprintf("%vh:", hour)
	}
	fmt.Printf("%v%vmin\n", hourText, min)
}

func showTimeOverview(intervals []Interval, iterations int) {
	completeTime := time.Duration(0)
	workTime := time.Duration(0)
	for _, interval := range intervals {
		completeTime += interval.time
		if interval.mode == WORK {
			workTime += interval.time
		}
	}
	completeTime *= time.Duration(iterations)
	workTime *= time.Duration(iterations)
	fmt.Print("\ntotal time:   ")
	printTime(completeTime)

	fmt.Print("\nworking time:   ")
	printTime(workTime)
	finishTime := time.Now().Add(completeTime)
	fmt.Printf("finished at:  %v\n", finishTime.Local().Format("15:04"))
}

func showInfo(intervals []Interval, iterations int) {
	fmt.Println("\nyour plan:")
	for _, interval := range intervals {
		time := interval.time
		text := ""
		if interval.text != "" {
			text = fmt.Sprintf("text: %v", interval.text)
		}
		mode := modeDesc[interval.mode]
		fmt.Printf("  mode: %-12s  time: %v min  %v\n", mode, time.Minutes(), text)
	}
	fmt.Printf("\n  iterations: %v\n", iterations)
	showTimeOverview(intervals, iterations)
}

type planString string

func (s planString) nextChar(i int) rune {
	if len(s) >= i+2 {
		return []rune(s)[i+1]
	} else {
		return '_'
	}
}

func isValidMode(c rune) bool {
	for mode, _ := range modes {
		if c == mode {
			return true
		}
	}
	return false
}

func parsePlan(plan planString) []Interval {
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
				if currentMode == CUSTOM {
					parsingText = true
				} else {
					minutes, _ := strconv.Atoi(timeStr)
					result = append(result, Interval{
						mode: currentMode,
						text: "",
						time: time.Duration(time.Duration(minutes) * time.Minute),
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
					mode: currentMode,
					text: "",
					time: time.Duration(time.Duration(minutes) * time.Minute),
				})
				text = ""
				continue
			} else {
				text += string(c)
			}
		}
		// start to parse a new mode
		if !parsingTime && !parsingText {
			if currentMode, ok = modes[c]; ok {
				if unicode.IsNumber(plan.nextChar(i)) {
					parsingTime = true
				} else {
					if currentMode == CUSTOM {
						panic("Error: custom intervals need a specified time")
					}
					result = append(result, Interval{
						mode: currentMode,
						text: "",
						time: modeTime[currentMode],
					})
				}
			} else {
				panic(fmt.Errorf("Error: invalid mode %v", c))
			}
		}
	}
	//check if parsing finished successfully
	if parsingTime {
		minutes, _ := strconv.Atoi(timeStr)
		result = append(result, Interval{
			mode: currentMode,
			text: "",
			time: time.Duration(time.Duration(minutes) * time.Minute),
		})
	}
	if parsingText {
		panic("Error: missing ':' to end custom task")
	}
	return result
}

func main() {
	plan := ""
	iterations := 1
	showPlan := false
	showTime := false
	/*if paramCount() >= 1{}
	    for opt in getopt():
	      if opt.kind == cmdArgument:
	        if fileExists(opt.key):
	          # read tasks from the file
	          var tmp = readFile(paramStr(1))
	          for c in Whitespace:
	            tmp = tmp.replace(&"{c}","")
	          plan &= tmp
	        else:
	          # read tasks from the comandline
	          plan &= opt.key
	      else:
	        # parse Options
	        case opt.key:
	          of "h","help":
	            fmt.Println(help)
	            quit(0)
	          of "p","plan":
	            showPlan = true
	          of "t","time":
	            showTime = true
	          of "r","repeat":
	            try:
	              iterations = opt.val.parseInt()
	            except ValueError:
	              quit("Error: invalid value " & opt.val & " for option " & opt.key)
	          else:
	            quit("Error: invalid option " & opt.key)
	  else:
	*/
	plan = "wswswswlwswswsw"

	if plan == "" {
		os.Exit(0)
	}
	intervals := parsePlan(planString(plan))
	if showPlan {
		showInfo(intervals, iterations)
	} else if showTime {
		showTimeOverview(intervals, iterations)
	} else {
		showInfo(intervals, iterations)
		for i := 1; i <= iterations; i++ {
			if iterations > 1 {
				fmt.Printf("iteration %v of %v\n", i, iterations)
			}
			for _, interval := range intervals {
				time := interval.time
				text := modeText[interval.mode]
				if interval.text != "" {
					text = interval.text
				}
				fmt.Println(colorGreen + text + colorReset)
				if useNotifications {
					notify(text)
				}
				runPart(time)
			}
		}
		fmt.Println("")
	}
}
