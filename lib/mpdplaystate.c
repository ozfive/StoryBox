/*

	PREREQUISITES:
	    
	    A C99 compliant compiler (e.g. gcc)
    	Meson 0.37 and Ninja

		sudo apt-get update -y
		sudo apt-get install -y libmpdclient-dev
		git clone https://github.com/MusicPlayerDaemon/libmpdclient.git

		apt-get install ninja-build

		meson . output

		ninja -C output
		ninja -C output install
		
	COMPILE:
		
		gcc looks in /usr/include/mpd/ for libraries so make sure mpd 
		header files are in the mpd folder before trying to compile.
		NOTE: You need to reference the lmpdclient at the end of the 
		build arguments after the c files not before.

		gcc -o mpdplaystate mpdplaystate.c -lmpdclient

*/
// mpdplaystate.c
#include "mpdplaystate.h"
#include <mpd/client.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

#define MPD_HOST "localhost"
#define MPD_PORT 0
#define MPD_PASSWORD "yL25v21jRJGMOz6P3F"

unsigned int get_play_state() {
    struct mpd_connection *conn;
    struct mpd_status *status;
    enum mpd_state state = 0;

    conn = mpd_connection_new(MPD_HOST, MPD_PORT, 30000);
    if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS) {
        mpd_connection_free(conn);
        return 0;
    }

    mpd_run_password(conn, MPD_PASSWORD);
    status = mpd_run_status(conn);
    if (status != NULL) {
		// 0 = State unknown, 1 = STATE STOP, 2 = STATE PLAY, 3 = STATE PAUSE,
        state = mpd_status_get_state(status);
        mpd_status_free(status);
    }

    mpd_connection_free(conn);
    return state;
}

const char* get_play_state_string() {
	unsigned int state = get_play_state();
	switch (state) {
		case 1:
			return "STOP";
		case 2:
			return "PLAY";
		case 3:
			return "PAUSE";
		default:
			return "UNKNOWN";
	}
}