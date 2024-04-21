package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// initialize and/or load config data structure
func InitializeConfig() (*Config, error) {
	var c Config
	history, err := InitializeDebugLog()
	if err != nil {
		c.Debug.Enable(false)
		return &c, err
	}
	c.Debug = history
	err = c.Load() // on any error, use default config
	if err != nil {
		c.UseDefaults()
		c.Debug.Print("Using Default Config")
		return &c, err
	}
	return &c, nil
}

type Config struct {
	BarSize     int      `json:"barsize"`
	BarStyle    int      `json:"barstyle"`
	Icon        int      `json:"icon"`
	TaskLength  int      `json:"tasklength"`
	Restart     bool     `json:"restart"`
	Bell        bool     `json:"bell"`
	HideTime    bool     `json:"hidetime"`
	HideTask    bool     `json:"hidetask"`
	HideSeconds bool     `json:"hideseconds"`
	HideIcon    bool     `json:"hideicon"`
	HideBar     bool     `json:"hidebar"`
	ReverseTime bool     `json:"reverse"`
	Percent     bool     `json:"percent"`
	Notify      bool     `json:"notify"`
	NotifyTmux  bool     `json:"tmux"`
	Log         bool     `json:"log"`
	Debug       *History `json:"-"`
}

func (c *Config) SetRestart(state bool)     { c.Restart = state }
func (c *Config) SetBell(state bool)        { c.Bell = state }
func (c *Config) SetPercent(state bool)     { c.Percent = state }
func (c *Config) SetReverse(state bool)     { c.ReverseTime = state }
func (c *Config) SetHideTime(state bool)    { c.HideTime = state }
func (c *Config) SetHideTask(state bool)    { c.HideTask = state }
func (c *Config) SetHideSeconds(state bool) { c.HideSeconds = state }
func (c *Config) SetHideBar(state bool)     { c.HideBar = state }
func (c *Config) SetHideIcon(state bool)    { c.HideIcon = state }
func (c *Config) SetBarSize(v int)          { c.BarSize = v }
func (c *Config) SetBarStyle(v int)         { c.BarStyle = v }
func (c *Config) SetIcon(v int)             { c.Icon = v }
func (c *Config) SetNotify(state bool)      { c.Notify = state }
func (c *Config) SetNotifyTmux(state bool)  { c.NotifyTmux = state }

// convert bytes (from file) to configuration struct
func (c *Config) Unmarshal(bytes []byte) error {
	err := json.Unmarshal(bytes, c)
	if err != nil {
		return err
	}
	return nil
}

// convert configuration to bytes
func (c Config) Marshal() ([]byte, error) {
	json, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return nil, err
	}
	return json, nil
}

// save configuration to file using config standard path, create if necessary
func (c *Config) Save() error {
	bytes, err := c.Marshal()
	if err != nil {
		return err
	}

	path, err := os.UserConfigDir()
	if err != nil {
		c.Debug.Print("UserConfigDir() failed, using binary path", err)
		path, err = os.Executable()
		if err != nil {
			return err
		}
	}
	
	configuration := filepath.Join(path, programName, ConfigFile)	
	_, err = checkFilePath(filepath.Dir(configuration)) // make sure directory is present
	if err != nil {
		err = createDirectory(filepath.Dir(configuration))
		if err != nil {
			return err
		}
	}

	err = writeFile(configuration, bytes)
	if err != nil {
		return err
	}
	return nil
}

// load configuration from file, use default configuration on io error
func (c *Config) Load() error {
	var path, loadFile string

	path, err := os.UserConfigDir()
	if err != nil {
		c.Debug.Print("UserConfigDir() failed, using binary path", err)
		path, err = os.Executable()
		if err != nil {
			c.Debug.Print("os.Executable() failed to get binary path", err)
			return err
		}
		path = filepath.Dir(path)
	}

	loadFile = filepath.Join(path, programName, ConfigFile)
	bytes, err := readFile(loadFile)
	if err != nil {
		return err
	}
	err = c.Unmarshal(bytes)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) UseDefaults() {
	c.BarSize = 10
	c.ReverseTime = true
	c.TaskLength = 20
}
