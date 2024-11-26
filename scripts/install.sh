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

# Check if sudo is used
if [ "$(id -u)" != 0 ]; then
  print_error 'Sorry, you need to run this script with sudo'
  exit 1
fi

# Update and upgrade the system
print_heading 'Updating system...'
apt update
apt upgrade -y

# Install required packages
print_heading 'Installing required packages...'
apt install -y \
  libmpdclient-dev \
  gcc \
  meson \
  ninja-build \
  sqlite3 \
  python3 \
  python3-pip \
  mpc \
  mpd \
  mpg123 \
  libasound2-dev \

# Prepare the directory structure for MPD
print_heading 'Setting up MPD directories...'
base_dir="/home/pi"
mpd_dir="$base_dir/.mpd"
music_dir="$base_dir/music"
playlists_dir="$mpd_dir/playlists"

mkdir -p "$mpd_dir"
mkdir -p "$music_dir"
mkdir -p "$playlists_dir"

touch "$mpd_dir/database"
touch "$mpd_dir/log"
touch "$mpd_dir/pid"

# Install gTTS library
print_heading 'Installing gTTS library...'
if pip show gTTS >/dev/null 2>&1; then
    echo "gTTS is already installed."
else
    echo "Installing gTTS..."
    pip install gTTS
fi

# Set up Go environment
print_heading 'Setting up Go environment...'
go_base_dir="$base_dir/go"
go_src_dir="$go_base_dir/src"
go_pkg_dir="$go_base_dir/pkg"
go_bin_dir="$go_base_dir/bin"

mkdir -p "$go_src_dir"
mkdir -p "$go_pkg_dir"
mkdir -p "$go_bin_dir"

# Download and install Go
print_heading 'Downloading and installing Go...'
go_tar_url="https://go.dev/dl/go1.20.5.linux-armv6l.tar.gz"

wget "$go_tar_url" && \
tar -C /usr/local -xzf go1.20.5.linux-armv6l.tar.gz && \
rm go1.20.5.linux-armv6l.tar.gz

echo "export PATH=\$PATH:/usr/local/go/bin" >> "$base_dir/.bashrc"
echo "export GOPATH=\$HOME/go" >> "$base_dir/.bashrc"
source "$base_dir/.bashrc"

# Clone the StoryBox repository
print_heading 'Cloning the StoryBox repository...'
storybox_repo="https://github.com/ozfive/StoryBox.git"

git clone "$storybox_repo"
cp "$base_dir/StoryBox/lib/mpd.conf" "$mpd_dir/mpd.conf"

# Set up phatbeat
print_heading 'Setting up phatbeat...'
cd "$base_dir" && \
shell_scripts_repo="https://github.com/ozfive/StoryBoxShellScripts.git" && \
git clone "$shell_scripts_repo" && \
cd StoryBoxShellScripts && \
chmod +x phatbeat.sh && \
./phatbeat.sh

# Install Wi-Fi configuration for wittypi3mini
print_heading 'Installing Wi-Fi configuration for wittypi3mini...'
witty_install_script="$base_dir/StoryBoxShellScripts/wittypi3mini/install.sh"

if [ -x "$witty_install_script" ]; then
  "$witty_install_script"
else
  print_warning 'Wi-Fi configuration install script not found. Skipping...'
fi

# Enable SPI interface
print_heading 'Enabling SPI interface...'
raspi-config nonint do_spi 0

# Build and install libmpdclient
print_heading 'Building and installing libmpdclient...'
libmpdclient_dir="$base_dir/StoryBox/lib/libmpdclient"

cd "$libmpdclient_dir" && \
meson . output && \
ninja -C output && \
ninja -C output install

# Build mpdcurrentsong
print_heading 'Building mpdcurrentsong...'
mpdcurrentsong_dir="$base_dir/StoryBox/lib"

cd "$mpdcurrentsong_dir" && \
gcc -o mpdcurrentsong mpdcurrentsong.c -lmpdclient

# Build mpdplaystate
print_heading 'Building mpdplaystate...'
gcc -o mpdplaystate "$mpdcurrentsong_dir/mpdplaystate.c" -lmpdclient

# Build mpdtime
print_heading 'Building mpdtime...'
gcc -o mpdtime "$mpdcurrentsong_dir/mpdtime.c" -lmpdclient

# Move the binaries to /usr/local/bin/
print_heading 'Moving binaries...'
cp "$mpdcurrentsong_dir/mpdcurrentsong" /usr/local/bin/mpdcurrentsong
cp "$mpdcurrentsong_dir/mpdplaystate" /usr/local/bin/mpdplaystate
cp "$mpdcurrentsong_dir/mpdtime" /usr/local/bin/mpdtime

# Build StoryBox binary
print_heading 'Building StoryBox binary...'
storybox_binary_dir="$go_src_dir/StoryBox"

cd "$storybox_binary_dir" && \
go build -o "$go_src_dir/StoryBox/StoryBox"

# Copy StoryBox binary
print_heading 'Copying StoryBox binary...'
cp "$go_src_dir/StoryBox/StoryBox" /usr/local/bin/StoryBox/

# Clone the StoryBox-Startup repository
print_heading 'Cloning the StoryBox-Startup repository...'
storybox_startup_repo="$go_src_dir/StoryBox-Startup"

git clone "https://github.com/ozfive/StoryBox-Startup.git" "$storybox_startup_repo"

# Build the Startup binary
print_heading 'Building the Startup binary...'
cd "$storybox_startup_repo" && \
go build -o "$go_src_dir/StoryBox-Startup/Startup"

# Copy the Startup binary
print_heading 'Copying the Startup binary...'
cp "$go_src_dir/StoryBox-Startup/Startup" /usr/local/bin

# Set up systemd service
print_heading 'Setting up systemd service...'
systemd_service_file="/lib/systemd/system/storyboxstartup.service"

cp "$storybox_startup_repo/storyboxstartup.service" "$systemd_service_file"
chmod 644 "$systemd_service_file"
systemctl enable storyboxstartup.service

# Set up sound files
print_heading 'Setting up sound files...'
sound_dir="/etc/sound"
started_mp3_file="$sound_dir/started.mp3"

if [ ! -d "$sound_dir" ]; then
    mkdir "$sound_dir"
fi

if [ ! -f "$started_mp3_file" ]; then
    print_heading 'Copying started.mp3 to /etc/sound/'
    started_mp3_source="$storybox_startup_repo/started.mp3"

    cp "$started_mp3_source" "$started_mp3_file"
    print_heading 'Copy completed...'
else
    print_warning 'File exists. Skipping...'
fi

# Copy the rest of the sound files to /etc/sound/
sound_files_source="$go_src_dir/StoryBox/sys-audio/"

cp -R "$sound_files_source" "$sound_dir/"

print_success 'Installation completed successfully!'

print_warning 'Please review the script and make any necessary adjustments specific to your environment before executing it.'

# Start the storyboxstartup.service file
systemctl start storyboxstartup.service