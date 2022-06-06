package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MnlPhlp/pomgo/modes"
	"github.com/MnlPhlp/pomgo/parsing"
	"github.com/cheggaaa/pb/v3"
	"github.com/mattn/go-tty"
)

const (
	useNotifications = true
	colorReset       = "\033[0m"
	colorGreen       = "\033[32m"
)

//go:embed help.txt
var help string

func remTimeStr(rem time.Duration) string {
	min := int(rem.Minutes())
	sec := int(rem.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d remaining", min, sec)
}

func notify(text string) {}

func runPart(runTime time.Duration) {
	start := time.Now()
	remTime := runTime
	seconds := int(runTime.Seconds())

	tmpl := `{{ bar . "[" "=" ">" "." "]"}} {{percent . "%3.f%%"}} {{string . "remaining" | green}}`
	var bar = pb.ProgressBarTemplate(tmpl).Start(seconds)
	bar.Set("remaining", remTimeStr(remTime))
	bar.SetMaxWidth(100)

	sleepTime := time.Second - time.Since(start)
	for i := 0; i < seconds; i++ {
		start = time.Now()
		time.Sleep(sleepTime)
		bar.Increment()
		remTime -= time.Second
		bar.Set("remaining", remTimeStr(remTime))
		// calculate sleep time to adjust for delays
		sleepTime = 2*time.Second - time.Since(start)
	}

	bar.Finish()
	fmt.Println()
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

func showTimeOverview(intervals []parsing.Interval, iterations int) {
	completeTime := time.Duration(0)
	workTime := time.Duration(0)
	for _, interval := range intervals {
		completeTime += interval.Time
		if interval.Mode == modes.WORK {
			workTime += interval.Time
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

func showInfo(intervals []parsing.Interval, iterations int) {
	fmt.Println("\nyour plan:")
	for _, interval := range intervals {
		time := interval.Time
		text := ""
		if interval.Text != "" {
			text = fmt.Sprintf("text: %v", interval.Text)
		}
		mode := modes.Desc[interval.Mode]
		fmt.Printf("  mode: %-12s  time: %v min  %v\n", mode, time.Minutes(), text)
	}
	fmt.Printf("\n  iterations: %v\n", iterations)
	showTimeOverview(intervals, iterations)
}

func main() {
	plan := ""
	iterations := 1
	showPlan := false
	showTime := false

	go watchExit()
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
	if len(os.Args) == 2 {
		plan = os.Args[1]
	}

	if plan == "" {
		os.Exit(0)
	}
	intervals := parsing.ParsePlan(parsing.PlanString(plan))
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
				time := interval.Time
				text := modes.Text[interval.Mode]
				if interval.Text != "" {
					text = interval.Text
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

func watchExit() {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	for {
		r, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		if r == 'q' {
			fmt.Println("")
			os.Exit(0)
		}
	}
}
