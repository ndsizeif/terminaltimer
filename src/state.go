package main

import (
	"encoding/json"
	"time"
	"path/filepath"
	"os"
)
// initialize and/or load state data structure along with a history logger
func InitializeState() (*State, error) {
	var s State
	debug, err := InitializeDebugLog()
	if err != nil {
		s.Debug.Enable(false) // don't attempt to log
	}
	s.Debug = debug
	s.Load()
	return &s, nil
}

func (s *State) UseDefaults() {
	s.TimeStart = time.Now()
	s.TimePause = time.Time{}
	s.TimeInterval = time.Duration(25 * time.Minute)
	s.TimeBreak = time.Duration(5 * time.Minute)
	s.TimeAlert = time.Duration(2 * time.Minute)
	s.Task = ""
}

type State struct {
	TimeStart    time.Time     `json:"start"`
	TimePause    time.Time     `json:"pause"`
	TimeInterval time.Duration `json:"interval"`
	TimeBreak    time.Duration `json:"break"`
	TimeAlert    time.Duration `json:"alert"`
	Task         string        `json:"task"`
	Debug        *History      `json:"-"`
}

func (s *State) GetTask() string { return s.Task }

// zero start time returns stopped timmer
func (s *State) TimerIsStopped() bool { return s.TimeStart.IsZero() }

// non-zero pause time returns paused timer
func (s *State) TimerIsPaused() bool { return !s.TimePause.IsZero() }

// returns if timer is inside alert window
func (s *State) TimerOnAlert() bool {
	return time.Since(s.TimeStart) > s.TimeInterval-s.TimeAlert && time.Since(s.TimeStart) < s.TimeInterval
}

// returns if elapsed time is greater than interval but less than interval+break; if on break, subtract pause time
func (s *State) TimerOnBreak() bool {
	if s.TimePause.IsZero() {
		return time.Since(s.TimeStart) > s.TimeInterval && time.Since(s.TimeStart) < (s.TimeInterval+s.TimeBreak)
	}
	if !s.TimePause.IsZero() {
		t := time.Since(s.TimeStart) - time.Since(s.TimePause)
		return t > s.TimeInterval && t < s.TimeInterval+s.TimeBreak
	}
	return false
}

// elapsed time is greater than interval + break; timer is not paused
func (s *State) TimerHasExpired() bool {
	return time.Since(s.TimeStart) >= s.TimeInterval+s.TimeBreak && s.TimePause.IsZero()
}

// time since s.TimeStart aka time.Now().Sub(t.TimeStart)
func (s *State) GetTotal() time.Duration {
	return s.TimeInterval + s.TimeBreak
}

// time since s.TimeStart aka time.Now().Sub(t.TimeStart)
func (s *State) GetElapsed() time.Duration {
	return time.Since(s.TimeStart)
}
func (s *State) GetElapsedBreak() time.Duration {
	return time.Since(s.TimeStart) - s.TimeInterval
}
func (s *State) GetElapsedPaused() time.Duration {
	return time.Since(s.TimeStart) - time.Since(s.TimePause)
}
func (s *State) GetElapsedPausedBreak() time.Duration {
	return time.Since(s.TimeStart) - time.Since(s.TimePause) - s.TimeInterval
}

// interval time remaining plus pause time, can be negative
func (s *State) GetRemainingPaused() time.Duration {
	return (s.TimeInterval + time.Since(s.TimePause)) - time.Since(s.TimeStart)
}
func (s *State) GetRemainingPausedBreak() time.Duration {
	return (s.TimeInterval + s.TimeBreak + time.Since(s.TimePause)) - time.Since(s.TimeStart)
}
func (s *State) GetRemainingBreak() time.Duration {
	return (s.TimeInterval + s.TimeBreak) - time.Since(s.TimeStart)
}

// interval time remaining, can be negative
func (s *State) GetIntervalRemaining() time.Duration {
	return s.TimeInterval - time.Since(s.TimeStart)
}

func (s *State) SetStart(v time.Time)        { s.TimeStart = v }
func (s *State) SetPause(v time.Time)        { s.TimePause = v }
func (s *State) SetInterval(v time.Duration) { s.TimeInterval = v }
func (s *State) SetBreak(v time.Duration)    { s.TimeBreak = v }
func (s *State) SetAlert(v time.Duration)    { s.TimeAlert = v }
func (s *State) SetTask(v string)            { s.Task = v }

func (s *State) ClearTask() error { // return nil error to satisfy map[string]func() err
	s.Task = ""
	return nil
}

// convert bytes (from file) to state struct
func (s *State) Unmarshal(bytes []byte) error {
	err := json.Unmarshal(bytes, s)
	if err != nil {
		s.Debug.Print(err)
		return err
	}
	return nil
}
// convert state structure to bytes
func (s *State) Marshal() ([]byte, error) {
	json, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		s.Debug.Print(err)
		return nil, err
	}
	return json, nil
}

// save state to file using path priority, create if necessary
func (s *State) Save() error {
	bytes, err := s.Marshal()
	if err != nil {
		s.Debug.Print(err)
		s.Debug.Print(s.Debug.Trace())
		return err
	}

	path, err := os.UserCacheDir()
	if err != nil {
			s.Debug.Print("UserCacheDir() failed, using binary path", err)
		path, err = os.Executable()
		if err != nil {
			s.Debug.Print("os.Executable() failed", err)
			return err
		}
	}

	saveFile := filepath.Join(path, programName, StateFile)

	_, err = checkFilePath(filepath.Dir(saveFile)) // is directory present
	if err != nil {
		err = createDirectory(filepath.Dir(saveFile))
		if err != nil {
			s.Debug.Print("could not create parent", err)
			return err
		}
	}
	s.Debug.Print("writing file:", saveFile)
	err = writeFile(saveFile, bytes)
	if err != nil {
		s.Debug.Print(err)
		return err
	}
	return nil
}

func (s *State) Load() error {
	var path, loadFile string

	path, err := os.UserCacheDir()
	if err != nil {
		s.Debug.Print("UserCacheDir() failed, using binary path", err)
		path, err = os.Executable()
		if err != nil {
			s.Debug.Print("os.Executable() failed to get binary path", err)
			return err
		}
	}

	loadFile = filepath.Join(path, programName, StateFile)
	bytes, err := readFile(loadFile)
	if err != nil {
		s.UseDefaults()
		return nil // using defaults corrects read error
	}
	err = s.Unmarshal(bytes)
	if err != nil {
		s.UseDefaults()
		return nil // using defaults corrects unmarshal error
	}
	return nil
}
