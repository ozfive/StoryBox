package mpdcurrentsong

/*
#cgo CFLAGS: -I/home/chris/go/src/StoryBox/lib/include
#cgo LDFLAGS: -L/home/chris/go/src/StoryBox/lib -lmpdcurrentsong
#include "mpdcurrentsong.h"
#include <stdlib.h>
*/
import "C"

import (
	"errors"
)

// SongInfo mirrors the C SongInfo struct
type SongInfo struct {
	Error       string
	SongPos     int
	ElapsedTime int
	ElapsedMs   uint
	TotalTime   int
	BitRate     int
	SampleRate  int
	Bits        int
	Channels    int
	Artist      string
	Album       string
	Title       string
	Track       string
	Name        string
	Date        string
	URI         string
	Duration    uint
	Pos         uint
}

// GetCurrentSongInfo calls the C function and returns a Go SongInfo struct
func GetCurrentSongInfo() (*SongInfo, error) {
	cInfo := C.get_current_song_info()
	if cInfo == nil {
		return nil, errors.New("failed to retrieve song information")
	}
	defer C.free_song_info(cInfo)

	var info SongInfo

	if cInfo.error != nil {
		info.Error = C.GoString(cInfo.error)
	}

	info.SongPos = int(cInfo.song_pos)
	info.ElapsedTime = int(cInfo.elapsed_time)
	info.ElapsedMs = uint(cInfo.elapsed_ms)
	info.TotalTime = int(cInfo.total_time)
	info.BitRate = int(cInfo.bit_rate)
	info.SampleRate = int(cInfo.sample_rate)
	info.Bits = int(cInfo.bits)
	info.Channels = int(cInfo.channels)

	if cInfo.artist != nil {
		info.Artist = C.GoString(cInfo.artist)
	}
	if cInfo.album != nil {
		info.Album = C.GoString(cInfo.album)
	}
	if cInfo.title != nil {
		info.Title = C.GoString(cInfo.title)
	}
	if cInfo.track != nil {
		info.Track = C.GoString(cInfo.track)
	}
	if cInfo.name != nil {
		info.Name = C.GoString(cInfo.name)
	}
	if cInfo.date != nil {
		info.Date = C.GoString(cInfo.date)
	}
	if cInfo.uri != nil {
		info.URI = C.GoString(cInfo.uri)
	}

	info.Duration = uint(cInfo.duration)
	info.Pos = uint(cInfo.pos)

	return &info, nil
}
