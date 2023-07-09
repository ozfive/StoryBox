#!/bin/bash
# file: install.sh
#
# This script will install the required packages/scripts
# for a working StoryBox system.

# Color codes for fancy output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No color

# Function to print a fancy heading
print_heading() {
  printf "${GREEN}\n%s${NC}\n" "$1"
}

# Function to print a fancy success message
print_success() {
  printf "${GREEN}%s${NC}\n" "$1"
}

# Function to print a fancy error message
print_error() {
  printf "${RED}%s${NC}\n" "$1"
}

# Function to print a fancy warning message
print_warning() {
  printf "${YELLOW}%s${NC}\n" "$1"
}

# check if sudo is used
if [ "$(id -u)" != 0 ]; then
  print_error 'Sorry, you need to run this script with sudo'
  exit 1
fi

print_heading 'Updating system...'
# Update and upgrade the system
apt update
apt upgrade -y

print_heading 'Installing required packages...'
# Install required packages
apt install -y libmpdclient-dev gcc meson ninja-build sqlite3 python3 python3-pip mpc mpd mpg123 libasound2-dev git

print_heading 'Setting up MPD directories...'
# Prepare the directory structure for MPD
mkdir -p /home/chris/.mpd/
mkdir -p /home/chris/music/
mkdir -p /home/chris/.mpd/playlists/

touch /home/chris/.mpd/database
touch /home/chris/.mpd/log
touch /home/chris/.mpd/pid

mv /home/chris/StoryBox/lib/mpd.conf /home/chris/.mpd/mpd.conf

print_heading 'Installing gTTS library...'
# Install gTTS library
pip install gTTS

print_heading 'Setting up Go environment...'
# Make all of the directories needed for Go.
mkdir -p /home/chris/go/src
mkdir -p /home/chris/go/pkg
mkdir -p /home/chris/go/bin

print_heading 'Downloading and installing Go...'
# Download the Go compiler.
wget https://go.dev/dl/go1.20.5.linux-armv6l.tar.gz

# Untar and gunzip the file.
tar -C /usr/local -xzf go1.20.5.linux-armv6l.tar.gz

# Remove the Go tar.gz file.
rm go1.20.5.linux-armv6l.tar.gz

# Set Go environment variables
echo 'export PATH=$PATH:/usr/local/go/bin' >> /home/chris/.bashrc
echo 'export GOPATH=$HOME/go' >> /home/chris/.bashrc

# shellcheck source=/dev/null
source /home/chris/.bashrc

print_heading 'Cloning the StoryBox repository...'
# Clone the StoryBox repository
git clone https://github.com/ozfive/StoryBox.git

print_heading 'Setting up phatbeat...'
# Move into /home/chris/ to git clone the StoryBoxShellScripts repo
cd /home/chris/ || exit

# Git clone the StoryBoxShellScripts repo.
git clone https://github.com/ozfive/StoryBoxShellScripts.git

# Move into StoryboxShellScripts/
cd StoryBoxShellScripts || exit

# Make the phatbeat.sh file executable.
chmod +x phatbeat.sh

# Execute phatbeat.sh
./phatbeat.sh

print_heading 'Installing Wi-Fi configuration for wittypi3mini...'
# Run the install.sh script in the wittypi3mini repository.
cd wittypi3mini || exit

# Make the install.sh script executable.
chmod +x install.sh
./install.sh

print_heading 'Enabling SPI interface...'
# Enable SPI interface
raspi-config nonint do_spi 0

print_heading 'Building and installing libmpdclient...'
# Build and install libmpdclient
cd /home/chris/StoryBox/lib/ || exit

# Git clone the libmpdclient repo in /home/chris/StoryBox/lib/
git clone https://github.com/MusicPlayerDaemon/libmpdclient.git

# Move into libmpdclient/
cd libmpdclient || exit

# use Meson and Ninja to build and install the libmpdclient
meson . output
ninja -C output
ninja -C output install

print_heading 'Building mpdcurrentsong...'
# Move into /home/chris/StoryBox/lib/
cd /home/chris/StoryBox/lib/ || exit

# Build and move mpdcurrentsong, mpdplaystate, and mpdtime
gcc -o mpdcurrentsong mpdcurrentsong.c -lmpdclient

print_heading 'Building mpdplaystate...'
gcc -o mpdplaystate mpdplaystate.c -lmpdclient

```bash
print_heading 'Building mpdtime...'
gcc -o mpdtime mpdtime.c -lmpdclient

# Move the binaries to /usr/local/bin/
mv mpdcurrentsong /usr/local/bin/mpdcurrentsong
mv mpdplaystate /usr/local/bin/mpdplaystate
mv mpdtime /usr/local/bin/mpdtime

print_heading 'Building StoryBox binary...'
# Build StoryBox binary in the Startup directory
## go build -o /home/chris/go/src/StoryBox/StoryBox || exit

print_heading 'Copying StoryBox binary...'
## cp /home/chris/go/src/StoryBox/StoryBox /usr/local/bin/StoryBox/ || exit

print_heading 'Cloning the StoryBox-Startup repository...'
cd /home/chris/go/src/ || exit

# Clone the StoryBox-Startup repository
git clone https://github.com/ozfive/StoryBox-Startup.git

print_heading 'Building the Startup binary...'
# Build Startup binary in the StoryBox-Startup directory
## go build -o /home/chris/go/src/StoryBox-Startup/Startup || exit

print_heading 'Copying the Startup binary...'
# Copy the Startup binary to the bin directory to make it available to the system
## cp /home/chris/go/src/StoryBox-Startup/Startup /usr/local/bin || exit

print_heading 'Setting up systemd service...'
# Copy storyboxstartup.service file to lib/systemd/system/
cp /home/chris/go/src/StoryBox/storyboxstartup.service /lib/systemd/system/storyboxstartup.service

# Set the permissions for the storyboxstartup.service file
chmod 644 /lib/systemd/system/storyboxstartup.service

# Enable the storyboxstartup.service file
systemctl enable storyboxstartup.service

print_heading 'Setting up sound files...'
mkdir /etc/sound/

# Copy the started.mp3 file to /etc/sound/
cp /home/chris/go/src/StoryBox-Startup/started.mp3 /etc/sound/started.mp3

# Copy the rest of the sound files to /etc/sound/
cp /home/chris/go/src/StoryBox/sys-audio/*.mp3 /etc/sound/

print_success 'Installation completed successfully!'

print_warning 'Please review the script and make any necessary adjustments specific to your environment before executing it.'

# Start the storyboxstartup.service file
systemctl start storyboxstartup.service
