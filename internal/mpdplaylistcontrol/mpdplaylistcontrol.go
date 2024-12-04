// mpdplaylistcontrol.go
package mpdplaylistcontrol

/*
#cgo CFLAGS: -I/home/chris/go/src/StoryBox/lib/include
#cgo LDFLAGS: -L/home/chris/go/src/StoryBox/lib -lmpdplaylistcontrol
#include "mpdplaylistcontrol.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// RenamePlaylist changes the name of a specified stored playlist.
func RenamePlaylist(conn *C.struct_mpd_connection, oldName, newName string) error {
	cOldName := C.CString(oldName)
	defer C.free(unsafe.Pointer(cOldName))

	cNewName := C.CString(newName)
	defer C.free(unsafe.Pointer(cNewName))

	C.rename_playlist(conn, cOldName, cNewName)

	if C.mpd_connection_get_error(conn) != C.MPD_ERROR_SUCCESS {
		errMsg := C.mpd_connection_get_error_message(conn)
		return fmt.Errorf("rename_playlist error: %s", C.GoString(errMsg))
	}

	return nil
}

// RemovePlaylist deletes a specified stored playlist.
func RemovePlaylist(conn *C.struct_mpd_connection, name string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	C.remove_playlist(conn, cName)

	if C.mpd_connection_get_error(conn) != C.MPD_ERROR_SUCCESS {
		errMsg := C.mpd_connection_get_error_message(conn)
		return fmt.Errorf("remove_playlist error: %s", C.GoString(errMsg))
	}

	return nil
}

// RemoveSongFromPlaylist removes a song from the current playlist by its ID.
func RemoveSongFromPlaylist(conn *C.struct_mpd_connection, songID uint) error {
	C.remove_song_from_playlist(conn, C.unsigned(songID))

	if C.mpd_connection_get_error(conn) != C.MPD_ERROR_SUCCESS {
		errMsg := C.mpd_connection_get_error_message(conn)
		return fmt.Errorf("remove_song_from_playlist error: %s", C.GoString(errMsg))
	}

	return nil
}

// RemoveSongFromPlaylistAtPos removes a song from the current playlist at the specified position.
func RemoveSongFromPlaylistAtPos(conn *C.struct_mpd_connection, position uint) error {
	C.remove_song_from_playlist_at_pos(conn, C.unsigned(position))

	if C.mpd_connection_get_error(conn) != C.MPD_ERROR_SUCCESS {
		errMsg := C.mpd_connection_get_error_message(conn)
		return fmt.Errorf("remove_song_from_playlist_at_pos error: %s", C.GoString(errMsg))
	}

	return nil
}

// MoveSongOnPlaylistToPos moves a song to a new position in the current playlist.
func MoveSongOnPlaylistToPos(conn *C.struct_mpd_connection, songID, newPos uint) error {
	C.move_song_on_playlist_to_pos(conn, C.unsigned(songID), C.unsigned(newPos))

	if C.mpd_connection_get_error(conn) != C.MPD_ERROR_SUCCESS {
		errMsg := C.mpd_connection_get_error_message(conn)
		return fmt.Errorf("move_song_on_playlist_to_pos error: %s", C.GoString(errMsg))
	}

	return nil
}

// AddSongToPlaylist adds a song to the current playlist at the specified position.
func AddSongToPlaylist(conn *C.struct_mpd_connection, songName, uri string, position uint) error {
	cSongName := C.CString(songName)
	defer C.free(unsafe.Pointer(cSongName))

	cURI := C.CString(uri)
	defer C.free(unsafe.Pointer(cURI))

	C.add_song_to_playlist(conn, cSongName, cURI, C.unsigned(position))

	if C.mpd_connection_get_error(conn) != C.MPD_ERROR_SUCCESS {
		errMsg := C.mpd_connection_get_error_message(conn)
		return fmt.Errorf("add_song_to_playlist error: %s", C.GoString(errMsg))
	}

	return nil
}
