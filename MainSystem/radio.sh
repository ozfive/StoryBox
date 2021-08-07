#!/bin/bash

# Script to generate an mpc playlist containing up-to-date BBC stream
# locations.  This can be run every hour or so as a cronjob to keep
# the non-static BBC streams working.  My version differs from the
# ones I've seen online in that each stream is added to a single
# playlist, rather than individual playlist files. This makes it
# easier to use with MPoD and via command line with "mpc play 4" for
# Radio 4.

# The cronjob can be set up using crontab -e:

# 00 */2 * * * /home/pi/scripts/radio.sh >> /tmp/radio.log

# This runs at the top of the hour every two hours, with the output going to a tmp file to check things are working

# Elements borrowed from:
# http://www.codedefied.co.uk/2011/12/24/playing-bbc-radio-streams-with-mpd/
# http://thenated0g.wordpress.com/2013/06/06/raspberry-pi-add-bbc1-6-radio-streams-and-mpc-play-command/
# http://www.gebbl.net/2013/10/playing-internet-radio-streams-mpdmpc-little-bash-python/
# http://jigfoot.com/hatworshipblog/?p=60
# http://www.raspberrypi.org/forums/viewtopic.php?t=50501&p=391258

# Set file paths and names
playlistdir=/var/lib/mpd/playlists
filename=playlist.m3u
playlistname=playlist
filepath=${playlistdir}/${filename}

# Array of BBC higher quality streams

declare -a streams=(http://www.bbc.co.uk/radio/listen/live/r1_aaclca.pls
http://www.bbc.co.uk/radio/listen/live/r2_aaclca.pls
http://www.bbc.co.uk/radio/listen/live/r3_aaclca.pls
http://www.bbc.co.uk/radio/listen/live/r4_aaclca.pls
http://www.bbc.co.uk/radio/listen/live/r5l_aaclca.pls
http://www.bbc.co.uk/radio/listen/live/r6_aaclca.pls
http://www.bbc.co.uk/radio/listen/live/r1x_aaclca.pls
http://www.bbc.co.uk/radio/listen/live/r4x_aaclca.pls
http://www.bbc.co.uk/radio/listen/live/r5lsp_aaclca.pls)

# Array of stream names for the above list - make sure they are in the same order!

declare -a names=("Abinayam FM
Star FM - Tamil Radio
TamilOne Radio CH
Nesaganam Online Tamil Radio
TamilSun FM
Southradios.com - IR Radio
A9Radio - Tamil Music Only
BTC Tamil FM
BBC Radio - Radio 5 Live"
"BBC Radio - Radio 1"
"RBBC Radio - Radio 2"
"BBC Radio - Radio 3"
"BBC Radio - Radio 4"
"BBC Radio - RBBC World Service News"
"BBC Radio - Radio 4 Extra"
)

echo "Updating playlist"

# Remove previous playlist

rm $filepath

echo "#EXTM3U" >> "$filepath"

# Iterate over the streams / names arrays to get the latest stream location

start=0
length=${#streams[@]}
end=$((length - 1))

for i in $(eval echo "{$start..$end}")
do 
#echo ${names[$i]}
# Places stream name in playlist
echo "#EXTINF:-1, ${names[$i]}" >> "$filepath"
# Places stream location in playlist
#curl -s ${streams[$i]} | grep "File1=" | sed 's/File1=//g' >> $filepath
wget -qO - ${streams[$i]} | grep "File1=" | sed 's/File1=//g' >> $filepath
done

# Adds in Magic105.4 static stream
echo "# #EXTINF:-1, Abinayam FM
# #EXTINF:-1, Star FM - Tamil Radio
# #EXTINF:-1, TamilOne Radio CH
# #EXTINF:-1, Nesaganam Online Tamil Radio
# #EXTINF:-1, TamilSun FM
# #EXTINF:-1, Southradios.com - IR Radio
# #EXTINF:-1, A9Radio - Tamil Music Only
# #EXTINF:-1, BTC Tamil FM
# #EXTINF:-1, Radio City Tamil
# #EXTINF:-1, Star Radio Tamil
# #EXTINF:-1, BBC Radio - Radio 5 Live
# #EXTINF:-1, BBC Radio - Radio 1
# #EXTINF:-1, BBC Radio - Radio 2
# #EXTINF:-1, BBC Radio - Radio 3
# #EXTINF:-1, BBC Radio - Radio 4
# #EXTINF:-1, BBC Radio - RBBC World Service News
# #EXTINF:-1, BBC Radio - Radio 4 Extra
# #EXTINF:-1, KUOW - NPR News High Quality
# #EXTINF:-1, NPR Program Stream
" >> "$filepath"
# echo "http://icy-e-02.sharp-stream.com:80/magic1054.mp3" >> $filepath

# Adds in 'static' intl BBC streams in case the cron job fails to run - always have a backup BBC stream that should work
# echo "http://bbcmedia.ic.llnwd.net/stream/bbcmedia_intl_lc_radio1_p?s=1365376033
# http://bbcmedia.ic.llnwd.net/stream/bbcmedia_intl_lc_radio2_p?s=1365376067
# http://bbcmedia.ic.llnwd.net/stream/bbcmedia_intl_lc_radio3_p?s=1365376123
# http://bbcmedia.ic.llnwd.net/stream/bbcmedia_intl_lc_radio4_p?s=1365376126
# http://bbcmedia.ic.llnwd.net/stream/bbcmedia_intl_lc_5live_p?s=1365376271
# http://bbcmedia.ic.llnwd.net/stream/bbcmedia_intl_lc_6music_p?s=1365376386" >> $filepath

echo "http://192.99.170.8:5756/
http://s1.voscast.com:8734/stream/1/
http://www.tamilone.ch:8000/stream_128
http://192.95.39.65:5206/
http://192.99.4.210:3596/
http://212.83.138.48:8324/stream/1/
http://195.154.217.103:8175/stream/1/
http://prclive1.listenon.in:9948/
http://stream11.shoutcastsolutions.com:9277/stream
http://bbcmedia.ic.llnwd.net/stream/bbcmedia_radio5live_mf_q
http://bbcmedia.ic.llnwd.net/stream/bbcmedia_radio1_mf_q
http://bbcmedia.ic.llnwd.net/stream/bbcmedia_radio2_mf_q
http://bbcmedia.ic.llnwd.net/stream/bbcmedia_radio3_mf_q
http://bbcmedia.ic.llnwd.net/stream/bbcmedia_radio4fm_mf_q
http://bbcwssc.ic.llnwd.net/stream/bbcwssc_mp1_ws-einws_backup 
http://bbcmedia.ic.llnwd.net/stream/bbcmedia_radio4extra_mf_q
http://18693.live.streamtheworld.com:3690/KUOWFM_HIGH_MP3_SC
https://npr-ice.streamguys1.com/live.mp3
https://18193.live.streamtheworld.com/SAM02AAC287.mp3" >> $filepath

# Clear existing mpc playlist and reload with the just generated playlist
mpc --host alraune22@localhost clear
mpc --host alraune22@localhost load ${playlistname}

echo "BBC streams updated at $(date)"
