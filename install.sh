#!/bin/bash

${XDG_CONFIG_HOME:=$HOME/.config}
${XDG_CACHE_HOME:=$HOME/.cache}

plumadoro_config=$XDG_CONFIG_HOME/.plumadoro.toml
plumadoro_log=$XDG_CACHE_HOME/.plumadoro_log.csv

plumadoro_def_config=./plumadoro.default.toml
plumadoro_def_log=./plumadoro_log.default.csv

if [! -f plumadoro_config]; then
	cp plumadoro_def_config plumadoro_config
fi

if [! -f plumadoro_log]; then
	cp plumadoro_def_log plumadoro_log
fi


