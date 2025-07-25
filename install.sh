#!/bin/bash

if [ -z ${XDG_CONFIG_HOME} ]; then
	XDG_CONFIG_HOME=$HOME/.config
fi

if [ -z ${XDG_CACHE_HOME} ]; then
	XDG_CACHE_HOME=$HOME/.cache
fi

plumadoro_config=$XDG_CONFIG_HOME/plumadoro.toml
plumadoro_log=$XDG_CACHE_HOME/.plumadoro_log.csv

plumadoro_def_config=./plumadoro.toml
plumadoro_def_log=./plumadoro_log.csv

if [ ! -f "$plumadoro_config" ]; then
	cp "$plumadoro_def_config" "$plumadoro_config"
fi

if [ ! -f "$plumadoro_log" ]; then
	cp "$plumadoro_def_log" "$plumadoro_log"
fi

# TODO: make a way to check what's the default local installation path
install -D -t /usr/local/bin ./bin/plumadoro
