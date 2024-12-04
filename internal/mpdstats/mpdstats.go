package mpdstats

/*
#cgo LDFLAGS: -L. -lmpdstats
#include <stdlib.h>

// Function declarations
unsigned int get_play_state();
const char* get_current_song_info();
unsigned int get_elapsed_time();
unsigned int get_total_time();
*/
import "C"

import (
	"errors"
	"unsafe"
)

// GetPlayState retrieves the current play state from MPD
func GetPlayState() (uint, error) {
	state := C.get_play_state()
	if state == 0 {
		return 0, errors.New("failed to retrieve play state")
	}
	return uint(state), nil
}

// GetCurrentSongInfo calls the C function to get the current song info
func GetCurrentSongInfo() (string, error) {
	currentSongInfo := C.get_current_song_info()
	if currentSongInfo == nil {
		return "", errors.New("failed to retrieve current song")
	}
	defer C.free(unsafe.Pointer(currentSongInfo))
	return C.GoString(currentSongInfo), nil
}

// GetElapsedTime calls the C function to get the elapsed time
func GetElapsedTime() (uint, error) {
	elapsed := C.get_elapsed_time()
	if elapsed == 0 {
		return 0, errors.New("failed to retrieve elapsed time")
	}
	return uint(elapsed), nil
}

// GetTotalTime calls the C function to get the total time
func GetTotalTime() (uint, error) {
	total := C.get_total_time()
	if total == 0 {
		return 0, errors.New("failed to retrieve total time")
	}
	return uint(total), nil
}
