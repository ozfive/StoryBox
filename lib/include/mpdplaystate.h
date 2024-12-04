// mpdplaystate.h
#ifndef MPDPLAYSTATE_H
#define MPDPLAYSTATE_H

#ifdef __cplusplus
extern "C" {
#endif

unsigned int get_play_state();
const char* get_play_state_string();

#ifdef __cplusplus
}
#endif

#endif // MPDPLAYSTATE_H