package main

import (
	_ "embed"
	"bytes"
	"time"
	"log"

	"github.com/hajimehoshi/go-mp3"
	"github.com/ebitengine/oto/v3"
)

// TODO: support different alarm types & make them configurable
// NOTE: The audio files are embeded into the binary during the compilation (not a bad thing)

//go:embed assets/alarm.mp3
var alarm []byte

func PlayAlarm() {
	go func() {
		decoder, err := mp3.NewDecoder(bytes.NewReader(alarm))
		if err != nil {
			log.Fatal("`assets/alarm.mp3` can't be decoded.")
		}

		options := &oto.NewContextOptions{}
		options.SampleRate = 44100	//  WTH IS THIS
		options.ChannelCount = 2
		options.Format = oto.FormatSignedInt16LE // format of the source. go-mp3's format is signed 16bit integers.

		context, readyChan, err := oto.NewContext(options)
		if err != nil {
			return // XXX: User's sound didn't work for whatever reason (not my problem)
		}

		<-readyChan

		player := context.NewPlayer(decoder)

		player.Play()

		for player.IsPlaying() {
			time.Sleep(time.Millisecond)
		}
	} ()
}


