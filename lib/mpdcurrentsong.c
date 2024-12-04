/*

	PREREQUISITES:
	    
	    A C99 compliant compiler (e.g. gcc)
    	Meson 0.37 and Ninja

		sudo apt-get update -y
		sudo apt-get install -y libmpdclient-dev
		git clone https://github.com/MusicPlayerDaemon/libmpdclient.git

		sudo apt install meson
		
		apt-get install ninja-build

		meson . output

		ninja -C output
		ninja -C output install
		
	COMPILE:
		
		gcc looks in /usr/include/mpd/ for libraries so make sure mpd 
		header files are in the mpd folder before trying to compile.
		NOTE: You need to reference the lmpdclient at the end of the 
		build arguments after the c files not before.

		gcc -o mpdcurrentsong mpdcurrentsong.c -lmpdclient

*/
#include "mpdcurrentsong.h"
#include <mpd/client.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>


#define MPD_HOST "localhost"
#define MPD_PORT 0
#define MPD_PASSWORD "yL25v21jRJGMOz6P3F"

static void copy_string(char** dest, const char* src) {
	if (src) {
		*dest = strdup(src);
	} else {
		*dest = strdup("Unknown");
	}
}

SongInfo* get_current_song_info() {
	struct mpd_connection *conn;
	struct mpd_status *status;
	struct mpd_song *song;
	const struct mpd_audio_format* audio_format;
	SongInfo* info = malloc(sizeof(SongInfo));
	if (!info) return NULL;

	memset(info, 0, sizeof(SongInfo));

	conn = mpd_connection_new(MPD_HOST, MPD_PORT, 30000);
	if (!conn || mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS) {
		if (conn) mpd_connection_free(conn);
		info->error = strdup(mpd_connection_get_error_message(conn));
		return info;
	}

	mpd_run_password(conn, MPD_PASSWORD);
	status = mpd_run_status(conn);
	if (!status) {
		info->error = strdup(mpd_connection_get_error_message(conn));
		mpd_connection_free(conn);
		return info;
	}

	const char* error_msg = mpd_status_get_error(status);
	if (error_msg) {
		info->error = strdup(error_msg);
		mpd_connection_free(conn);
		mpd_status_free(status);
		return info;
	}

	if (mpd_status_get_state(status) == MPD_STATE_PLAY ||
        mpd_status_get_state(status) == MPD_STATE_PAUSE) {
        info->song_pos = mpd_status_get_song_pos(status);
        info->elapsed_time = mpd_status_get_elapsed_time(status);
        info->elapsed_ms = mpd_status_get_elapsed_ms(status);
        info->total_time = mpd_status_get_total_time(status);
        info->bit_rate = mpd_status_get_kbit_rate(status);
    }

	audio_format = mpd_status_get_audio_format(status);
	if (audio_format) {
		info->sample_rate = audio_format->sample_rate;
		info->bits = audio_format->bits;
		info->channels = audio_format->channels;
	}

	mpd_status_free(status);

	if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS) 
		goto cleanup;

	mpd_response_next(conn);

	while ((song = mpd_recv_song(conn)) != NULL) {
		copy_string(&info->artist, mpd_song_get_tag(song, MPD_TAG_ARTIST, 0));
		copy_string(&info->album, mpd_song_get_tag(song, MPD_TAG_ALBUM, 0));
		copy_string(&info->title, mpd_song_get_tag(song, MPD_TAG_TITLE, 0));
		copy_string(&info->track, mpd_song_get_tag(song, MPD_TAG_TRACK, 0));
		copy_string(&info->name, mpd_song_get_tag(song, MPD_TAG_NAME, 0));
		copy_string(&info->date, mpd_song_get_tag(song, MPD_TAG_DATE, 0));
		copy_string(&info->uri, mpd_song_get_uri(song));
		
		if (mpd_song_get_duration(song) > 0) {
			info->duration = mpd_song_get_duration(song);
		}
		
		info->pos = mpd_song_get_pos(song);
		
		mpd_song_free(song);
	}

	if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS ||
        !mpd_response_finish(conn)) {
        info->error = strdup("Response Error");
    }

cleanup:
	mpd_connection_free(conn);
	return info;
}

void free_song_info(SongInfo* info) {
	if (info) {
		free(info->error);
		free(info->artist);
		free(info->album);
		free(info->title);
		free(info->track);
		free(info->name);
		free(info->date);
		free(info->uri);
		free(info);
	}
}