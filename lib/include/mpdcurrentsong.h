// mpdcurrentsong.h

#ifndef MPDCURRENTSONG_H
#define MPDCURRENTSONG_H

#include <stdlib.h>

// Define the SongInfo struct
typedef struct {
    char* error;
    int song_pos;
    int elapsed_time;
    unsigned int elapsed_ms;
    int total_time;
    int bit_rate;
    int sample_rate;
    int bits;
    int channels;
    char* artist;
    char* album;
    char* title;
    char* track;
    char* name;
    char* date;
    char* uri;
    unsigned int duration;
    unsigned int pos;
} SongInfo;

// Function declarations
SongInfo* get_current_song_info();
void free_song_info(SongInfo* info);

#endif // MPDCURRENTSONG_H