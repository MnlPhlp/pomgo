package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/MnlPhlp/pomgo/display"
	"github.com/MnlPhlp/pomgo/modes"
	"github.com/MnlPhlp/pomgo/parsing"
	"github.com/cheggaaa/pb/v3"
	"github.com/mattn/go-tty"
)

const (
	colorReset = "\033[0m"
	colorGreen = "\033[32m"
)

//go:embed help.txt
var help string

func remTimeStr(rem time.Duration) string {
	min := int(rem.Minutes())
	sec := int(rem.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d remaining", min, sec)
}

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

func main() {
	go watchExit()

	plan := ""
	iterations := 1
	showPlan := false
	showTime := false
	noNotifications := false

	flag.BoolVar(&showPlan, "p", false, "show plan and exit")
	flag.BoolVar(&showTime, "t", false, "show time and exit")
	flag.IntVar(&iterations, "r", 1, "set numer of iterations")
	flag.BoolVar(&noNotifications, "n", false, "disable notifications")
	flag.Usage = func() {
		fmt.Println(help)
	}

	flag.Parse()

	plan = "wswswswlwswswsw"
	if len(flag.Args()) > 0 {
		plan = ""
		for i := 0; i < len(flag.Args()); i++ {
			plan += flag.Args()[0]
		}
		if file, err := os.ReadFile(plan); err == nil {
			plan = strings.TrimSpace(string(file))
		}
	}

	if plan == "" {
		os.Exit(0)
	}
	intervals := parsing.ParsePlan(parsing.PlanString(plan))
	if showPlan {
		display.ShowInfo(intervals, iterations)
	} else if showTime {
		display.ShowTimeOverview(intervals, iterations)
	} else {
		display.ShowInfo(intervals, iterations)
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
				if !noNotifications {
					display.Notify(text)
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
