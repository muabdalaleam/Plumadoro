<p align="center">
  <img src="https://github.com/muabdalaleam/Plumadoro/blob/main/imgs/logo.svg?raw=true" alt="Plumadoro"/>
</p>

<h3 align="center">Minimal but featurefull Pomodoro TUI</h3>

---

## Description
A pomodoro terminal app (TUI) written in Go using the bubbletea framework, It supports `TOML` 
configuration throw a dotfile and a `CSV` log file

## Installation
You download the repo along with the binary from the releases page (link) to install it do the following:<br>
```
tar -xf plumadoro_<version>.tar.gz

cd plumadoro_<version>

sudo ./install.sh
```

And here you are. this installation process DOESN'T require you to have Go pre-installed

## Configuration
The default config `plumadoro.toml` file should exist in $XDG_CONFIG_HOME or in $HOME/.config if 
your XDG_* variables are not definded, for linux the config file should be: `~/.config/plumadoro.toml`
inside it you would find comments on how to tweak it to your own liking.

## Logging
By default the log file `plumadoro_log.csv` is inside $XDG_CACHE_HOME or in $HOME/.cache if your
XDG_* variables are not definded `~/.cache/plumadoro_log.csv.toml`, modifying this file can **break**
the logging and state restoring features in plumadoro, this file exist so you can analyze with any
software you'd like but be careful not to modify it.

## Features
- Customization throw a TOML file
- CSV log file to analyze your progress
- Alarming sound in the end of phases
- The ability to set a maximum pause time per phase or disable it
- Minimal, sleek interface
- Popup error system

## TODOs
- [ ] Support system notifications
- [ ] Add a help view
- [ ] Make sound effects configurable
- [ ] Add different sound effects
- [ ] Support gradient filled progress bar
- [ ] Make the key bindings configurable
- [ ] Modify the configuration tags' names to make more sense

## Contributing
This project is still in beta so it still have a lot of problems any contribution is much appreciated :)<br>
How to contribute to this project:
- Make an issue with your feature request to discuss how to do it
- Do a pull request with your approaved feature request, or bug fix 
- Thanks ⭐⭐
*(you don't have to do an issue first if your contribution is a bug fix)*

## Thanks for
- [Charmbraclet](https://github.com/charmbracelet), it's contributers
- [Burnt sushi](https://github.com/BurntSushi/toml) for the TOML parser
- [Hajime Hoshi](https://github.com/hajimehoshi) for the MP3 parser
- YOU for being awesome ❤️
