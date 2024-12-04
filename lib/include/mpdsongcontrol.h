// mpdplaystate.h
#ifndef MPDPLAYSTATE_H
#define MPDPLAYSTATE_H

#ifdef __cplusplus
extern "C" {
#endif

const char* get_play_state_string();
void play_song(unsigned int songpos);
void play_song_id(unsigned int songid);
void play();
void pause();
void stop();
void previous();
void next();
void seek(unsigned int songpos, unsigned int time);
void seeksongid(unsigned int songid, unsigned int time);
void seekcurrentsong(unsigned int time);
void playsongid(unsigned int id);

#ifdef __cplusplus
}
#endif

#endif // MPDPLAYSTATE_H