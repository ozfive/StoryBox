#include "mpdplaystate.h"
#include <mpd/client.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

#define MPD_HOST "localhost"
#define MPD_PORT 0
#define MPD_PASSWORD "yL25v21jRJGMOz6P3F"

// play_song Begins playing the playlist at song number SONGPOS
void play_song(unsigned int songpos) {
	struct mpd_connection *conn;
	conn = mpd_connection_new(MPD_HOST, MPD_PORT, 30000);
	if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS) {
		mpd_connection_free(conn);
		return 0;
	}

	mpd_run_password(conn, MPD_PASSWORD);
	mpd_run_play_pos(conn, songpos);

	mpd_connection_free(conn);
}

// play_song_id Begins playing the playlist at song SONGID
void play_song_id(unsigned int songid) {
	struct mpd_connection *conn;
	conn = mpd_connection_new(MPD_HOST, MPD_PORT, 30000);
	if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS) {
		mpd_connection_free(conn);
		return 0;
	}

	mpd_run_password(conn, MPD_PASSWORD);
	mpd_run_play_id(conn, songid);

	mpd_connection_free(conn);
}

// play: play the current song
void play() {
	struct mpd_connection *conn;
	conn = mpd_connection_new(MPD_HOST, MPD_PORT, 30000);
	if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS) {
		mpd_connection_free(conn);
		return 0;
	}

	mpd_run_password(conn, MPD_PASSWORD);
	mpd_run_play(conn);

	mpd_connection_free(conn);
}

// pause: pause the current song
void pause() {
	struct mpd_connection *conn;
	conn = mpd_connection_new(MPD_HOST, MPD_PORT, 30000);
	if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS) {
		mpd_connection_free(conn);
		return 0;
	}

	mpd_run_password(conn, MPD_PASSWORD);
	mpd_run_pause(conn, true);

	mpd_connection_free(conn);
}

// stop: stop the current song
void stop() {
	struct mpd_connection *conn;
	conn = mpd_connection_new(MPD_HOST, MPD_PORT, 30000);
	if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS) {
		mpd_connection_free(conn);
		return 0;
	}

	mpd_run_password(conn, MPD_PASSWORD);
	mpd_run_stop(conn);

	mpd_connection_free(conn);
}

// prev: play the previous song
void previous() {
	struct mpd_connection *conn;
	conn = mpd_connection_new(MPD_HOST, MPD_PORT, 30000);
	if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS) {
		mpd_connection_free(conn);
		return 0;
	}

	mpd_run_password(conn, MPD_PASSWORD);
	mpd_run_previous(conn);

	mpd_connection_free(conn);
}

// next: play the next song
void next() {
	struct mpd_connection *conn;
	conn = mpd_connection_new(MPD_HOST, MPD_PORT, 30000);
	if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS) {
		mpd_connection_free(conn);
		return 0;
	}

	mpd_run_password(conn, MPD_PASSWORD);
	mpd_run_next(conn);

	mpd_connection_free(conn);
}

// seek: seeks to the position TIME (in seconds; fractions allowed) of entry 
// SONGPOS in the playlist
void seek(unsigned int songpos, unsigned int time) {
	struct mpd_connection *conn;
	conn = mpd_connection_new(MPD_HOST, MPD_PORT, 30000);
	if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS) {
		mpd_connection_free(conn);
		return 0;
	}

	mpd_run_password(conn, MPD_PASSWORD);
	mpd_run_seek_pos_id(conn, songpos, time);

	mpd_connection_free(conn);
}

// seeksongid: seeks to the position TIME (in seconds; fractions allowed) of 
// song SONGID
void seeksongid(unsigned int songid, unsigned int time) {
	struct mpd_connection *conn;
	conn = mpd_connection_new(MPD_HOST, MPD_PORT, 30000);
	if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS) {
		mpd_connection_free(conn);
		return 0;
	}

	mpd_run_password(conn, MPD_PASSWORD);
	mpd_run_seek_id(conn, songid, time);

	mpd_connection_free(conn);
}

// seekcurrentsong: Seeks to the position TIME (in seconds; fractions allowed) 
// within the current song. If prefixed by + or -, then the time is relative 
// to the current playing position.
void seekcurrentsong(unsigned int time) {
	struct mpd_connection *conn;
	conn = mpd_connection_new(MPD_HOST, MPD_PORT, 30000);
	if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS) {
		mpd_connection_free(conn);
		return 0;
	}

	mpd_run_password(conn, MPD_PASSWORD);
	mpd_run_seek_id(conn, mpd_run_current_id(conn), time);

	mpd_connection_free(conn);
}

// playsongid: play a song by id
void playsongid(unsigned int id) {
	struct mpd_connection *conn;
	conn = mpd_connection_new(MPD_HOST, MPD_PORT, 30000);
	if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS) {
		mpd_connection_free(conn);
		return 0;
	}

	mpd_run_password(conn, MPD_PASSWORD);
	mpd_run_play_id(conn, id);

	mpd_connection_free(conn);
}