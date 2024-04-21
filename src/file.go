package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	StateFile     = "state.json"  // timer state file
	ConfigFile    = "config.json" // configuration file
	TimerFile     = "timer.log"   // timer log file
	DebugFile     = "debug.log"   // debug log file
	EnableDebug   = false         // false removes debug printing to log
	LogTimeFormat = time.Kitchen
)

var TimeStamp = time.Now().Format(LogTimeFormat)

// return an error if a path is invalid
func checkFilePath(path string) (string, error) {
	_, err := os.Stat(path)
	if err != nil {
		return path, err
	}
	return path, nil
}

// create a directory
func createDirectory(path string) error {
	err := os.Mkdir(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

// return bytes of file at given path
func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// write bytes to a file
func writeFile(path string, bytes []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

// delete file at given path
func deleteFile(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func RemoveLogFiles() error {
	path, err := os.UserCacheDir()
	if err == nil {
		path = filepath.Join(path, programName, TimerFile)
		err := deleteFile(path)
		if err != nil && !os.IsNotExist(err) { // stop if not a "can't be found" error
			return err
		}
	}
	path, err = os.Executable()
	if err == nil {
		path = filepath.Join(path, programName, TimerFile)
		err := deleteFile(path)
		if err != nil && !os.IsNotExist(err) { // stop if not a "can't be found" error
			return err
		}
	}
	return nil
}

// return pointer logfile, check for valid paths in priority order
func ReturnLogFile(filename string) (*os.File, error) {
	var file *os.File

	path, err := os.UserCacheDir()
	if err != nil {
		path, err = os.Executable()
		if err != nil {
			return nil, err
		}
	}

	path = filepath.Join(path, programName, filename)
	_, err = checkFilePath(filepath.Dir(path))
	if err != nil {
		err = createDirectory(filepath.Dir(path))
		return nil, err
	}

	file, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// initialize user viewable timer logs
func InitializeTimeLog() (*History, error) {
	var h History
	file, err := ReturnLogFile(TimerFile)
	if err != nil {
		return &h, err
	}
	timeLog := log.New(file, "", log.Lshortfile)
	timeLog.SetFlags(0) // remove default flags, start from scratch
	h.logger = timeLog
	return &h, nil
}

// initialize new log for debug messages
func InitializeDebugLog() (*History, error) {
	var h History
	file, err := ReturnLogFile(DebugFile)
	if err != nil {
		return &h, err
	}
	debugFlags := log.LstdFlags | log.Lshortfile
	debugLog := log.New(file, "DEBUG: ", debugFlags)
	h.logger = debugLog
	h.enabled = EnableDebug
	return &h, nil
}

type History struct {
	logger  *log.Logger
	file    *os.File
	enabled bool
}

// call Trace to return stack information, debug.Print(t.State.debug.Trace(), "string", ...)
func (h *History) Trace() string {
	var msg string
	_, file, line, ok := runtime.Caller(1)
	path, _ := os.Getwd()
	path = path + "/"
	if ok {
		msg = fmt.Sprintf("%s:%d", strings.TrimPrefix(file, path), line)
	}
	return msg
}

// toggle History output to log
func (h *History) Enable(state bool) {
	h.enabled = state
}

func (h *History) Print(v ...interface{}) {
	if h.enabled {
		h.logger.Println(v...)
	}
}
