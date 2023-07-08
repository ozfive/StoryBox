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
apt install -y libmpdclient-dev gcc meson ninja-build sqlite3 python3 python3-pip mpc mpd mpg123 libasound2-dev git

# Prepare the directory structure for MPD
mkdir /home/chris/.mpd/
mkdir /home/chris/music
mkdir /home/chris/.mpd/playlists
touch /home/chris/.mpd/database
touch /home/chris/.mpd/log
touch /home/chris/.mpd/pid

mv /home/chris/StoryBox/lib/mpd.conf /home/chris/.mpd/mpd.conf

# Install gTTS library
pip install gTTS

# Make all of the directories needed for Go.
mkdir /home/chris/go/
mkdir /home/chris/go/src
mkdir /home/chris/go/pkg
mkdir /home/chris/go/bin

# Download and install the Go compiler
wget https://go.dev/dl/go1.20.5.linux-armv6l.tar.gz
tar -C /usr/local -xzf go1.20.5.linux-armv6l.tar.gz
rm go1.20.5.linux-armv6l.tar.gz

#Set Go environment variables
echo "export PATH=$PATH:/usr/local/go/bin" >> /home/chris/.bashrc
echo "export GOPATH=\$HOME/go/" >> /home/chris/.bashrc

# shellcheck source=/dev/null
source /home/chris/.bashrc

echo "Cloning the StoryBox repo..."

# Clone the StoryBox repository
git clone https://github.com/ozfive/StoryBox.git

# Clone the StoryBoxShellScripts repository and execute phatbeat.sh
cd /home/chris/ || exit
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
cd /home/chris/StoryBox/lib/ || exit
git clone https://github.com/MusicPlayerDaemon/libmpdclient.git
cd libmpdclient || exit
meson . output
ninja -C output
ninja -C output install

# Build and move mpdcurrentsong, mpdplaystate, and mpdtime
cd /home/chris/StoryBox/lib/ || exit
gcc -o mpdcurrentsong mpdcurrentsong.c -lmpdclient
mv mpdcurrentsong /usr/local/bin/mpdcurrentsong

gcc -o mpdplaystate mpdplaystate.c -lmpdclient
mv mpdplaystate /usr/local/bin/mpdplaystate

gcc -o mpdtime mpdtime.c -lmpdclient
mv mpdtime /usr/local/bin/mpdtime

echo "Building StoryBox binary"

# Build StoryBox binary in the Startup directory
## go build -o /home/chris/go/src/StoryBox/StoryBox || exit

echo "Completed!"

echo "Copying StoryBox binary to /usr/local/bin/Storybox/"

## cp /home/chris/go/src/StoryBox/StoryBox /usr/local/bin/StoryBox/ || exit

echo "Completed!"

cd /home/chris/go/src/ || exit

echo "Cloning the StoryBox-Startup repo."

# Clone the StoryBox repository
git clone https://github.com/ozfive/StoryBox-Startup.git

echo "Building 'Startup' Binary"

# Build Startup binary in the Startup directory
## go build -o /home/chris/go/src/StoryBox-Startup/Startup || exit

echo "Completed!"

echo "Copying Startup binary to /usr/local/bin/"
# Copy the Startup binary to the bin directory to make it available to the system
## cp /home/chris/go/src/StoryBox-Startup/Startup /usr/local/bin || exit

echo "Completed!"

# Copy storyboxstartup.service file to lib/systemd/system/

echo "Copying storyboxstartup.service to /lib/systemd/system/"

cp /home/chris/go/src/StoryBox/storyboxstartup.service /lib/systemd/system/storyboxstartup.service

echo "Completed!"

mkdir /etc/sound/

# Copy the started.mp3 file to /etc/sound/
cp /home/chris/go/src/StoryBox/Startup/started.mp3 /etc/sound/started.mp3

# Copy the rest of the sound files to /etc/sound/
cp /home/chris/go/src/StoryBox/sys-audio/*.mp3 /etc/sound/

# Set the permissions for the storyboxstartup.service file
chmod 644 /lib/systemd/system/storyboxstartup.service

# Enable the storyboxstartup.service file
systemctl enable storyboxstartup.service

# Start the storyboxstartup.service file
systemctl start storyboxstartup.service