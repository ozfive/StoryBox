/*

	PREREQUISITES:
	    
	    A C99 compliant compiler (e.g. gcc)
    	Meson 0.37 and Ninja

		apt install libmpdclient-dev
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
#include <mpd/client.h>
#include <mpd/status.h>
#include <mpd/entity.h>
#include <mpd/search.h>
#include <mpd/song.h>
#include <mpd/tag.h>
#include <mpd/message.h>
#include <stdio.h>

static int
handle_error(struct mpd_connection *c)
{

	if (mpd_connection_get_error(c) != MPD_ERROR_SUCCESS) {
		fprintf(stderr, "%s\n", mpd_connection_get_error_message(c));
		mpd_connection_free(c);
		return -1;
	}

}

static void
print_tag(const struct mpd_song *song, enum mpd_tag_type type,
	  const char *label)
{
	unsigned i = 0;
	const char *value;

	while ((value = mpd_song_get_tag(song, type, i++)) != NULL)
		printf("%s\t%s\n", label, value);
}

int main(void)
{

	unsigned time = 0;
	struct mpd_connection *conn;
	struct mpd_status *status;
	enum mpd_state state;

	conn = mpd_connection_new("localhost", 0, 0);

	mpd_run_password(conn, "alraune22");

	status = mpd_run_status(conn);

	if (!status) {
		return 0;
	}

	struct mpd_song *song;
	const struct mpd_audio_format *audio_format;

	mpd_command_list_begin(conn, true);
	mpd_send_status(conn);
	mpd_send_current_song(conn);
	mpd_command_list_end(conn);

	status = mpd_recv_status(conn);
	if (status == NULL)
		return handle_error(conn);

	printf("volume\t%i\n", mpd_status_get_volume(status));
	printf("repeat\t%i\n", mpd_status_get_repeat(status));
	printf("queueversion\t%u\n", mpd_status_get_queue_version(status));
	printf("queuelength\t%i\n", mpd_status_get_queue_length(status));
	if (mpd_status_get_error(status) != NULL)
		printf("error\t%s\n", mpd_status_get_error(status));

	if (mpd_status_get_state(status) == MPD_STATE_PLAY ||
	    mpd_status_get_state(status) == MPD_STATE_PAUSE) {
		printf("song\t%i\n", mpd_status_get_song_pos(status));
		printf("elaspedTime\t%i\n",mpd_status_get_elapsed_time(status));
		printf("elasped_ms\t%u\n", mpd_status_get_elapsed_ms(status));
		printf("totalTime\t%i\n", mpd_status_get_total_time(status));
		printf("bitRate\t%i\n", mpd_status_get_kbit_rate(status));
	}

	audio_format = mpd_status_get_audio_format(status);
	if (audio_format != NULL) {
		printf("sampleRate\t%i\n", audio_format->sample_rate);
		printf("bits\t%i\n", audio_format->bits);
		printf("channels\t%i\n", audio_format->channels);
	}

	mpd_status_free(status);

	if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS)
		return handle_error(conn);

	mpd_response_next(conn);

	while ((song = mpd_recv_song(conn)) != NULL) {
		
		print_tag(song, MPD_TAG_ARTIST, "artist");
		print_tag(song, MPD_TAG_ALBUM, "album");
		print_tag(song, MPD_TAG_TITLE, "title");
		print_tag(song, MPD_TAG_TRACK, "track");
		print_tag(song, MPD_TAG_NAME, "name");
		print_tag(song, MPD_TAG_DATE, "date");
		printf("uri\t%s\n", mpd_song_get_uri(song));

		if (mpd_song_get_duration(song) > 0) {
			printf("time\t%u\n", mpd_song_get_duration(song));
		}

		printf("pos\t%u\n", mpd_song_get_pos(song));

		mpd_song_free(song);
	}

	if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS ||
	    !mpd_response_finish(conn))
		return handle_error(conn);

	mpd_status_free(status);

	mpd_connection_free(conn);

	return 0;
}