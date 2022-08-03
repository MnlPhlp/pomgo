package display

import (
	"fmt"
	"time"

	"github.com/MnlPhlp/pomgo/modes"
	"github.com/MnlPhlp/pomgo/parsing"
	"github.com/gen2brain/beeep"
)

func Notify(text string) {
	err := beeep.Notify("pomgo", text, "assets/information.png")
	if err != nil {
		panic(err)
	}
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

func ShowTimeOverview(intervals []parsing.Interval, iterations int) {
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

	fmt.Print("working time: ")
	printTime(workTime)
	finishTime := time.Now().Add(completeTime)
	fmt.Printf("finished at:  %v\n", finishTime.Local().Format("15:04"))
}

func ShowInfo(intervals []parsing.Interval, iterations int) {
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
	ShowTimeOverview(intervals, iterations)
}
