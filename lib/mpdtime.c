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

		gcc -o mpdtime mpdtime.c -lmpdclient

*/
#include <mpd/client.h>
#include <stdio.h>

int main(void)
{
	unsigned time = 0;
	struct mpd_connection *conn;
	struct mpd_status *status;
	enum mpd_state state;

	conn = mpd_connection_new("localhost", 0, 0);

	mpd_run_password(conn, "yL25v21jRJGMOz6P3F");

	status = mpd_run_status(conn);

	if (!status) {
		return 0;
	}

	time = mpd_status_get_elapsed_time(status);

	state = mpd_status_get_state(status);

	mpd_status_free(status);

	mpd_connection_free(conn);

	if (state > 1) {
		printf("%u\n", time);
	}

	return 0;
}