#!/bin/bash
# file: install.sh
#
# This script will install the required packages/scripts
# for a working StoryBox system.

: <<'DISCLAIMER'

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

This script is licensed under the terms of the MIT license.
Unless otherwise noted, code reproduced herein
was written for this script.

DISCLAIMER

# check if sudo is used
if [ "$(id -u)" != 0 ]; then
  echo 'Sorry, you need to run this script with sudo'
  exit 1
fi

# productname="StoryBox" # the name of the product to install
# scriptname="install"   # the name of this script

# Update and upgrade the system
apt update
apt upgrade -y

# Install required packages
apt install -y libmpdclient-dev gcc meson ninja-build sqlite3 python3 mpc mpd mpg123 libasound2-dev git

# Prepare the directory structure for MPD
mkdir ~/.mpd/
mkdir ~/music
mkdir ~/.mpd/playlists
touch ~/.mpd/database
touch ~/.mpd/log
touch ~/.mpd/pid

mv /home/pi/StoryBox/MainSystem/lib/mpd.conf ~/.mpd/mpd.conf

# Install gTTS library
pip install gTTS

# Download and install the Go compiler
wget https://golang.org/dl/go.1.19.2.linux-armv6l.tar.gz
tar -C /usr/local -xzf go.1.19.2.linux-armv6l.tar.gz
rm go.1.19.2.linux-armv6l.tar.gz


#Set Go environment variables
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
echo "export GOPATH=\$HOME/pi/go/" >> ~/.bashrc

# shellcheck source=/dev/null
source ~/.bashrc

# Clone the StoryBox repository
git clone https://github.com/ozfive/StoryBox.git

# Clone the StoryBoxShellScripts repository and execute phatbeat.sh
cd ~ || exit
git clone https://github.com/ozfive/StoryBoxShellScripts.git
cd StoryBoxShellScripts || exit
chmod +x phatbeat.sh
./phatbeat.sh

# Run the install.sh script in the wittypi3mini repository.
cd wittypi3mini || exit
chmod +x install.sh
./install.sh

# Enable SPI interface
raspi-config nonint do_spi 0

# Build and install libmpdclient
cd ~/StoryBox/lib/ || exit
git clone https://github.com/MusicPlayerDaemon/libmpdclient.git
cd libmpdclient || exit
meson . output
ninja -C output
ninja -C output install

# Build and move mpdcurrentsong, mpdplaystate, and mpdtime
cd ~/StoryBox/lib/ || exit
gcc -o mpdcurrentsong mpdcurrentsong.c -lmpdclient
mv mpdcurrentsong /usr/local/bin/mpdcurrentsong

gcc -o mpdplaystate mpdplaystate.c -lmpdclient
mv mpdplaystate /usr/local/bin/mpdplaystate

gcc -o mpdtime mpdtime.c -lmpdclient
mv mpdtime /usr/local/bin/mpdtime

# Build and move StoryBox
mkdir /home/pi/go/
mkdir /home/pi/go/src
mkdir /home/pi/go/pkg
mkdir /home/pi/go/bin

cp -r /home/pi/StoryBox /home/pi/go/src
cd /home/pi/go/src/Storybox/ || exit

go build -o StoryBox

mv StoryBox /usr/local/bin/Storybox

# Build Startup application in the Startup directory
cd /home/pi/go/src/Storybox/Startup || exit
go build -o Startup
chmod +x Startup

# Copy the Startup application to the bin directory to make it available to the system
cp Startup /usr/local/bin

# Copy storyboxstartup.service file to lib/systemd/system/
cd /home/pi/go/src/Storybox/ || exit
cp storyboxstartup.service /lib/systemd/system/storyboxstartup.service

# Copy the started.mp3 file to /etc/sound/
cp /home/pi/go/src/Storybox/Startup/started.mp3 /etc/sound/started.mp3

# Copy the rest of the sound files to /etc/sound/
cp /home/pi/go/src/Storybox/sys-audio/*.mp3 /etc/sound/

# Set the permissions for the storyboxstartup.service file
chmod 644 /lib/systemd/system/storyboxstartup.service

# Enable the storyboxstartup.service file
systemctl enable storyboxstartup.service

# Start the storyboxstartup.service file
systemctl start storyboxstartup.service

# Reboot the system
reboot now