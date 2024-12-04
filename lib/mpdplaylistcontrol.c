#include <mpd/client.h>
#include <mpd/playlist.h>
#include <stdio.h>
#include <stdlib.h>

struct mpd_connection* connect_mpd(const char *host, int port) {
    struct mpd_connection *conn;
    conn = mpd_connection_new(host, port, 30000);
    if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS) {
        fprintf(stderr, "Connection error: %s\n", mpd_connection_get_error_message(conn));
        mpd_connection_free(conn);
        return NULL;
    }
    return conn;
}

// Create a new playlist with the given name
void create_playlist(struct mpd_connection *conn, const char *playlist_name) {
    mpd_run_save(conn, playlist_name);
}

// Clear the current playlist
void clear_playlist(struct mpd_connection *conn, const char *name) {
    mpd_run_playlist_clear(conn, name);
}

// Load a playlist
void load_playlist(struct mpd_connection *conn, const char *name) {
    mpd_run_load(conn, name);
}

// Save the current playlist
void save_playlist(struct mpd_connection *conn) {
    mpd_run_save(conn, NULL);
}

// Remove a playlist
void remove_playlist(struct mpd_connection *conn, const char *name) {
    mpd_run_rm(conn, name);
}

// Rename a playlist 
void rename_playlist(struct mpd_connection *conn, const char *oldname, const char *newname) {
    mpd_run_rename(conn, oldname, newname);
}

// Remove a song from the playlist by its song position
void remove_song_from_playlist(struct mpd_connection *conn, unsigned int songpos) {
    mpd_run_delete_id(conn, songpos);
}

// Remove a song from the playlist at the given position
void remove_song_from_playlist_at_pos(struct mpd_connection *conn, unsigned int songpos) {
    mpd_run_delete_pos(conn, songpos);
}

// Move a song to a new position in the playlist
void move_song_on_playlist_to_pos(struct mpd_connection *conn, unsigned int songpos, unsigned int newpos) {
    mpd_run_move_id(conn, songpos, newpos);
}

// Add a song to the named playlist at position "to"
bool add_song_to_playlist(struct mpd_connection *conn, const char *song_name, const char *uri, unsigned to) {
    mpd_run_playlist_add_to(conn, song_name, uri, to);
    if (mpd_connection_get_error(conn) != MPD_ERROR_SUCCESS) {
        fprintf(stderr, "Error adding song to playlist: %s\n", mpd_connection_get_error_message(conn));
        return false;
    }
    return true;
}