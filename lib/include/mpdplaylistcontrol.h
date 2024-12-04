#ifndef MPDPLAYLISTCONTROL_H
#define MPDPLAYLISTCONTROL_H

#include <mpd/client.h>

#ifdef __cplusplus
extern "C" {
#endif

/**
 * connect_mpd Establishes a connection to the MPD server.
 *
 * @param host The hostname of the MPD server. Use NULL for localhost.
 * @param port The port number of the MPD server. Use 0 for the default port (6600).
 * @return A pointer to the mpd_connection structure if successful, NULL otherwise.
 */
struct mpd_connection* connect_mpd(const char *host, int port);

/**
 * create_playlist creates a new stored playlist with the specified songs.
 *
 * @param conn           The active MPD connection.
 * @param name           The name of the playlist to create.
 */
void create_playlist(struct mpd_connection *conn, const char *name);

/**
 * clear_playlist clears all songs from a specified stored playlist.
 *
 * @param conn           The active MPD connection.
 * @param name           The name of the playlist to clear.
 */
void clear_playlist(struct mpd_connection *conn, const char *name);

/**
 * load_playlist loads a specified stored playlist into the current playback queue.
 *
 * @param conn           The active MPD connection.
 * @param name           The name of the playlist to load.
 */
void load_playlist(struct mpd_connection *conn, const char *name);

/**
 * save_playlist saves the current playback queue as a stored playlist with the specified name.
 *
 * @param conn           The active MPD connection.
 * @param name           The name of the playlist to save.
 */
void save_playlist(struct mpd_connection *conn);

/**
 * remove_playlist deletes a specified stored playlist.
 *
 * @param conn           The active MPD connection.
 * @param name           The name of the playlist to delete.
 */

void remove_playlist(struct mpd_connection *conn, const char *name);

/**
 * rename_playlist changes the name of a specified stored playlist.
 * 
 * @param conn           The active MPD connection.
 * @param oldname        The current name of the playlist.
 * @param newname        The new name for the playlist.
 */
void rename_playlist(struct mpd_connection *conn, const char *oldname, const char *newname);

/**
 * remove_song_from_playlist removes a song from the current playlist by its ID.
 *
 * @param conn    The active MPD connection.
 * @param songpos The ID of the song to remove.
 */
void remove_song_from_playlist(struct mpd_connection *conn, unsigned int songpos);

/**
 * remove_song_from_playlist_at_pos removes a song from the current playlist at the specified position.
 *
 * @param conn    The active MPD connection.
 * @param songpos The zero-based position of the song to remove.
 */
void remove_song_from_playlist_at_pos(struct mpd_connection *conn, unsigned int songpos);

/**
 * move_song_on_playlist_to_pos moves a song to a new position in the current playlist.
 *
 * @param conn    The active MPD connection.
 * @param songpos The ID of the song to move.
 * @param newpos  The new position for the song.
 */
void move_song_on_playlist_to_pos(struct mpd_connection *conn, unsigned int songpos, unsigned int newpos);

/**
 * add_song_to_playlist adds a song to the current playlist at the specified position.
 * 
 * @param conn     The active MPD connection.
 * @param song_name The name of the playlist to add the song to.
 * @param uri       The URI of the song to add.
 * @param to        The position in the playlist to add the song.
 * @return true if the song was added successfully, false otherwise.
 */
void add_song_to_playlist(struct mpd_connection *conn, const char *song_name, const char *uri, unsigned to);

#ifdef __cplusplus
}
#endif

#endif // MPDPLAYLISTCONTROL_H