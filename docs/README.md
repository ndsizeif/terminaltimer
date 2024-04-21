# terminalTimer

![timer run](./assets/inline.gif)

terminalTimer is a basic interval timer that outputs progress to stdout in the console. The program
allows the user to display a string to represent the current task. Some users may find this useful
as they may forget what they are supposed to be working on during a given interval. The program
provides options to output progress in a variety of ways. 

![timerbar1](./assets/bar1.gif)
![timerbar2](./assets/bar2.gif)
![timerbar3](./assets/bar3.gif)
![timerbar4](./assets/bar4.gif)
![timerbar5](./assets/bar5.gif)
![timerbar6](./assets/bar6.gif)
![timerbar7](./assets/bar7.gif)

## Installation

### Release
Download and run the latest [release binary](https://github.com/ndsizeif/terminaltimer/releases) for
your system architecture. Alternatively, build `terminalTimer` for your system from the source
tarball. Install [Go](https://go.dev/doc/install) if not already present on your system. 

### Clone

You may also clone the main branch of the project, and build the project that way. Place the resulting binary
in your `$PATH` or run locally.

<details>
    <summary>code statistics</summary>

```
===============================================================================
 Language            Files        Lines         Code     Comments       Blanks
===============================================================================
 Go                     11         1886         1686           64          136
===============================================================================
 Total                  11         1886         1686           64          136
===============================================================================
```

</details>

```
git clone depth=1 https://github.com/ndsizeif/terminaltimer
cd src/terminaltimer
go build
./terminalTimer
```

## Usage

<details>
    <summary>terminalTimer --help</summary>

```
Usage of terminalTimer:
  start
        start timer
  stop
        stop timer
  pause
        pause timer
  resume
        start timer if paused
  break
        start break
  run
        display timer inline inside terminal
  clean
        delete timer log file
  clear
        clear the string for current task
  help
        display full help

Usage of terminalTimer set:
  -alert duration
        set threshold to start alert
  -break duration
        set break interval
  -timer duration
        set timer interval

Usage of terminalTimer style:
  -bar int
        style progress bar appearance
  -icon int
        style icon appearance
  -width int
        style progress bar width

Usage of terminalTimer toggle:
  -bar
        turn progress bar on/off
  -bell
        turn terminal bell on/off
  -clock
        display time on/off
  -icon
        display icons on/off
  -notify
        turn timer notifications on/off for notify-send
  -percent
        display interval percentage on/off
  -restart
        turn automatic timer restart on/off
  -reverse
        timer displays time descending/ascending
  -tmux
        turn timer notifications on/off for tmux
```
</details>

![timerInline](./assets/inlineRun.gif)

The timer can be run inline with `terminaltimer run`. When the timer is running inline, the
following key commands can be issued, followed by enter/carriage return. The entire word will
also be accepted.

<details>
    <summary>inline commands</summary> 


| cmd | action   |
| --- | -------- |
| s   | start    |
| t   | stop     |
| q   | quit     |
| p   | pause    |
| r   | resume   |
| b   | break    |
| c   | clear    |

</details>

## Configuration

The configuration file is responsible for the appearance and behavior of the timer. The file is
located in the user config directory as `/terminalTimer/config.json`. If a configuration file is not
present when the timer is started, the timer will function using default values. These config values
can be edited manually in the json file, or changed by issuing `terminalTimer style` or
`terminalTimer toggle` commands.

<details>
    <summary>example config.json</summary>

```
{
	"barsize": 23,
	"barstyle": 1,
	"icon": 2,
	"tasklength": 20,
	"restart": false,
	"bell": false,
	"hidetime": false,
	"hidetask": false,
	"hideseconds": false,
	"hideicon": false,
	"hidebar": false,
	"reverse": true,
	"percent": false,
	"notify": true,
	"tmux": false,
	"log": false
}
```

</details>

## Notifications

By default, the program will not notify the user when an interval is complete.  This behavior can be
toggled by changing the value of `notify` or `tmux` to `true` in the configuration file. 

When `notify` is set to `true` in the configuration file, if `notify-send` is installed on the
system, the program will send a notification message when a timer interval is completed. 

![notifySend](./assets/notifySend.gif)

When `tmux` is set to `true` in the configuration file, if a tmux session is active, the tmux client
will display a notification popup when a timer interval is completed.

![notifyTmux](./assets/notifyTmux.gif)

## Logging

The program can log intervals and tasks that have been completed throughout the day.  The log file
is saved to the user cache directory as `/terminalTimer/timer.log`. This behavior is disabled by
default, and can be toggled by changing the value of `log` in the configuration file. Logging is
basic and is a work in progress. The log will write duplicate records if multiple instances of timer
are running.

## Tips

Icons require Nerd Fonts to be installed.  There is an option to suppress icons, or to use ascii
characters instead.

```
terminalTimer style -icon 3
terminalTimer toggle -icon
```

Set an alias for common intervals. Example of valid duration strings are `1h`, `60m`, `59m60s`
```
alias tt="terminalTimer set -timer 25m -break 5m && terminalTimer start && terminalTimer run"
alias tth="terminalTimer set -timer 50m -break 10m && terminalTimer start && terminalTimer run"
```

<img align="right" width="200" src="./assets/menuTmux.gif" title="example tmux menu">

Invoking `terminalTimer` with no arguments will return a static render of the timer's current
status.  This can be used with other command-line tools. [Tmux](https://github.com/tmux/tmux) users
can incorporate the timer in their status line.  If the `terminalTimer` binary is in `$PATH`, it can
be called using
[#(terminalTimer)](https://github.com/tmux/tmux/wiki/Getting-Started#embedded-commands). 

<br><br>
`set -g status-left '#(terminalTimer)'` will embed the timer into the status line.
<br><br>

Tmux users can build their own menu to quickly manage timer settings.  It is also
a convenient way to start different preset intervals.

<br clear="right"/>

## Contributing

Bug reports, or any form of constructive feedback is appreciated. Feature requests are also welcome.
Forking the project and customizing it to your liking may yield the best results.
