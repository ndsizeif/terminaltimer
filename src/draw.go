package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

const redrawRate = 1 * time.Second   // inline redraw
var userInput bool                   // user input string accepted

func (t *Task) RunInline() error {
	end := make(chan struct{})
	quitting := make(chan os.Signal, 1)
	signal.Notify(quitting, os.Interrupt)
	redraw := time.NewTicker(redrawRate)
	Stty("-echo") // turn stty echo off while running inline
	fmt.Printf("%v%v%v", cursorPrevLine, clearLine, cursorHide)
	t.GetTime()
	t.Render()

	go func() {
		for {
			t.Read()
		}
	}()

	go func() {
		for {
			select {
			case <-quitting: // Ctrl + c
				Stty("echo")
				fmt.Printf("%v%v%v", clearLine, carriageReturn, cursorShow)
				os.Exit(0)
				return
			case <-redraw.C: // render at ticker rate unless user inputs command
				if userInput {
					userInput = false
					break
				}
				fmt.Printf("%v%v", clearLine, carriageReturn)
				t.GetTime()
				t.Render()
			}
		}
	}()

	<-end // never recieve an end struct, goroutines to run indefintely
	return nil
}

// receive user input when run inline, perform basic actions based on input
func (t *Task) Read() {
	var input string
	fmt.Scan(&input)
	if input != "" {
		// quit for any q* string to cover "quit"
		if strings.HasPrefix(input, "q") {
			Stty("echo")
			fmt.Printf("%v%v%v", clearLine, carriageReturn, cursorShow)
			os.Exit(0)
		}
		if input == "b" || input == "break" {
			userInput = true
			fmt.Printf("%v%vtake a break", clearLine, carriageReturn)
			t.Break()
			updatedTask, err := InitializeTimer()
			if err != nil {
				os.Exit(0)
			}
			t = &updatedTask // replace Task object, garbage collector should deal with the old object
			*t = updatedTask
			return
		}
		if input == "t" || input == "stop" {
			userInput = true
			fmt.Printf("%v%vstop timer", clearLine, carriageReturn)
			t.Stop()
			updatedTask, err := InitializeTimer()
			if err != nil {
				os.Exit(0)
			}
			t = &updatedTask
			return
		}
		if input == "s" || input == "start" {
			userInput = true
			fmt.Printf("%v%vstart timer", clearLine, carriageReturn)
			t.Start()
			updatedTask, err := InitializeTimer()
			if err != nil {
				os.Exit(0)
			}
			t = &updatedTask
			return
		}
		if input == "p" || input == "pause" {
			userInput = true
			fmt.Printf("%v%vpause timer", clearLine, carriageReturn)
			t.Pause()
			updatedTask, err := InitializeTimer()
			if err != nil {
				os.Exit(0)
			}
			t = &updatedTask
			return
		}
		if input == "r" || input == "resume" {
			userInput = true
			fmt.Printf("%v%vresume timer", clearLine, carriageReturn)
			t.Resume()
			updatedTask, err := InitializeTimer()
			if err != nil {
				os.Exit(0)
			}
			t = &updatedTask
			return
		}
		if input == "c" || input == "clear" { // clear terminal screen, but leave scrollback
			cmd := exec.Command("clear", "-x")
			cmd.Stdout = os.Stdout
			err := cmd.Run()
			if err != nil {
				t.Log.Print(err)
				return
			}
			return
		}
		// default case, clear line, go back to start previous line, and render timer immediately
		fmt.Printf("%v%v", clearLine, carriageReturn)
		t.Render()
	}
}

// prints timer output to terminal
func (t *Task) Render() {
	fmt.Printf("%v%v%v%v%v%v", t.DrawIcon(), t.DrawTask(), t.DrawBar(), t.DrawTime(), t.DrawPercent(), t.RingBell())
}

// display user task string
func (t *Task) DrawTask() string {
	task := t.State.GetTask()
	if len(task) == 0 || t.Config.HideTask {
		return ""
	}
	if len(task) > t.Config.TaskLength {
		return fmt.Sprintf("%v ", task[:t.Config.TaskLength])
	}
	return fmt.Sprintf("%v ", t.State.GetTask())
}

// display clock
func (t *Task) DrawTime() string {
	if t.Config.HideTime {
		return ""
	}
	if t.State.TimerIsStopped() {
		return ""
	}
	switch {
	case t.Config.ReverseTime:
		switch {
		case t.State.TimerHasExpired():
			return fmt.Sprintf(" %v", t.FormatTime(t.State.TimeInterval-t.State.TimeInterval))
		case t.State.TimerIsPaused() && t.State.TimerOnBreak():
			return fmt.Sprintf(" %v", t.FormatTime(t.State.GetRemainingPausedBreak()))
		case t.State.TimerIsPaused() && !t.State.TimerOnBreak():
			return fmt.Sprintf(" %v", t.FormatTime(t.State.GetRemainingPaused()))
		case !t.State.TimerIsPaused() && t.State.TimerOnBreak():
			return fmt.Sprintf(" %v", t.FormatTime(t.State.GetRemainingBreak()))
		default:
			return fmt.Sprintf(" %v", t.FormatTime(t.State.GetIntervalRemaining()))
		}
	default:
		switch {
		case t.State.TimerHasExpired():
			return fmt.Sprintf(" %v", t.FormatTime(t.State.TimeInterval))
		case t.State.TimerIsPaused() && t.State.TimerOnBreak():
			return fmt.Sprintf(" %v", t.FormatTime(t.State.GetElapsedPausedBreak()))
		case t.State.TimerIsPaused() && !t.State.TimerOnBreak():
			return fmt.Sprintf(" %v", t.FormatTime(t.State.GetElapsedPaused()))
		case !t.State.TimerIsPaused() && t.State.TimerOnBreak():
			return fmt.Sprintf(" %v", t.FormatTime(t.State.GetElapsedBreak()))
		default:
			return fmt.Sprintf(" %v", t.FormatTime(t.State.GetElapsed()))
		}
	}
}

// show time in % complete
func (t *Task) DrawPercent() string {
	if !t.Config.Percent {
		return ""
	}
	var percent int
	switch {
	case t.State.TimerIsStopped():
		return ""
	case t.State.TimerHasExpired():
		percent = 100
	case t.State.TimerOnBreak() && t.State.TimerIsPaused():
		percent = int(float64(t.State.GetElapsedPausedBreak()) / float64(t.State.TimeBreak) * 100.0)
	case t.State.TimerOnBreak() && !t.State.TimerIsPaused():
		percent = int(float64(t.State.GetElapsedBreak()) / float64(t.State.TimeBreak) * 100.0)
	case !t.State.TimerOnBreak() && t.State.TimerIsPaused():
		percent = int(float64(t.State.GetElapsedPaused()) / float64(t.State.TimeInterval) * 100.0)
	default:
		percent = int(float64(t.State.GetElapsed()) / float64(t.State.TimeInterval) * 100.0)
	}
	return fmt.Sprintf("%v%%", FormatPercent(percent))
}

// render progress as filled bar characters
func (t *Task) DrawBar() string {
	if t.Config.HideBar {
		return ""
	}
	var bar, done, todo string
	var scale int
	switch {
	case t.State.TimerIsStopped():
		return fmt.Sprintf("%v", bar)
	case t.State.TimerHasExpired():
		bar = strings.Repeat(t.Progress["done"], t.Config.BarSize)
		return fmt.Sprintf("%v", bar)
	case t.State.TimerIsPaused():
		if t.State.TimerOnBreak() {
			scale = int(float64(t.Config.BarSize) * float64(t.State.GetElapsedPausedBreak()) / float64(t.State.TimeBreak))
		}
		if !t.State.TimerOnBreak() {
			scale = int(float64(t.Config.BarSize) * float64(t.State.GetElapsedPaused()) / float64(t.State.TimeInterval))
		}
	case !t.State.TimerIsPaused():
		if t.State.TimerOnBreak() {
			scale = int(float64(t.Config.BarSize) * float64(t.State.GetElapsedBreak()) / float64(t.State.TimeBreak))
		}
		if !t.State.TimerOnBreak() {
			scale = int(float64(t.Config.BarSize) * float64(t.State.GetElapsed()) / float64(t.State.TimeInterval))
		}
	}
	if scale < 0 { // NOTE prevents negative repeat crashes.  This could happen if system time is changed while program is running
		return ""
	}
	done = strings.Repeat(t.Progress["done"], scale)
	if (t.Config.BarSize - scale) < 0 {
		t.State.Debug.Print("NEGATIVE BAR")
		t.State.Debug.Print("t.Config.Barsize: ", t.Config.BarSize, "scale: ", scale)
	}
	todo = strings.Repeat(t.Progress["todo"], t.Config.BarSize-scale)
	return fmt.Sprintf("%v%v", done, todo)
}

// display the timer's mode with an icon
func (t *Task) DrawIcon() string {
	if t.Config.HideIcon {
		return ""
	}
	switch {
	case t.State.TimerIsStopped():
		return fmt.Sprintf("%v ", t.Symbols["stopped"])
	case t.State.TimerHasExpired():
		return fmt.Sprintf("%v ", t.Symbols["expired"])
	case t.State.TimerIsPaused() && !t.State.TimerOnBreak():
		return fmt.Sprintf("%v ", t.Symbols["paused"])
	case t.State.TimerIsPaused() && t.State.TimerOnBreak():
		return fmt.Sprintf("%v ", t.Symbols["breakp"])
	case !t.State.TimerIsPaused() && t.State.TimerOnBreak():
		return fmt.Sprintf("%v ", t.Symbols["break"])
	case t.State.TimerOnAlert():
		return fmt.Sprintf("%v ", t.Symbols["warning"])
	default:
		return fmt.Sprintf("%v ", t.Symbols["on"])
	}
}

// sound old school terminal bell
func (t *Task) RingBell() string {
	// skip if disabled bell in config, timer is stopped or paused
	if !t.Config.Bell || t.State.TimerIsStopped() || t.State.TimerIsPaused() {
		return ""
	}
	// set threshold a few milliseconds greater than the refresh rate to ensure the bell triggers at least once
	switch {
	case t.State.TimerOnBreak():
		if (t.State.GetTotal()-t.State.GetElapsed()) < time.Duration(redrawRate+(100*time.Millisecond)) && (t.State.GetTotal()-t.State.GetElapsed()) > time.Duration(0*time.Second) {
			t.State.Debug.Print("BREAK COMPLETE BELL")
			return fmt.Sprintf("\a")
		}
	case !t.State.TimerOnBreak():
		if t.State.GetIntervalRemaining() < time.Duration(redrawRate+(100*time.Millisecond)) && t.State.GetIntervalRemaining() > time.Duration(0*time.Second) {
			t.State.Debug.Print("INTERVAL COMPLETE BELL")
			return fmt.Sprintf("\a")
		}
	}
	return ""
}
