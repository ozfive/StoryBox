#!/bin/bash

# file: install.sh
#
# This script installs all required packages, builds C and Go components,
# and sets up the StoryBox system.

# Exit immediately if a command exits with a non-zero status
set -e

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

# Ensure the script is run as root
if [ "$(id -u)" != 0 ]; then
  print_error 'Sudo is required for this script to run. Exiting...'
  exit 1
fi

# Update and upgrade the system
print_heading 'Updating system...'
apt update && apt upgrade -y

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
  git \
  wget \
  unzip \
  raspi-config \
  i2c-tools

# Install gTTS library
print_heading 'Installing gTTS library...'
if pip3 show gTTS >/dev/null 2>&1; then
    echo "gTTS is already installed."
else
    echo "Installing gTTS..."
    pip3 install gTTS
fi

### 
### Modify repo_owner to your GitHub username 
### if you want to use your own repository.
### set the go_version to the version of Go 
### you want to install.
### 
repo_owner="ozfive"
go_version="1.23.3"

# Set up Go environment variables
go_tarball="go${go_version}.linux-armv6l.tar.gz"
go_tar_url="https://go.dev/dl/$go_tarball"

base_dir="/home/pi"
go_base_dir="$base_dir/go"
go_src_dir="$go_base_dir/src"
go_pkg_dir="$go_base_dir/pkg"
go_bin_dir="$go_base_dir/bin"

# URLs of repositories to clone
storybox_repo="https://github.com/$repo_owner/StoryBox.git"
shell_scripts_repo="https://github.com/$repo_owner/StoryBoxShellScripts.git"

# Path to the Witty Pi 3 Mini install script
witty_install_script="$go_src_dir/StoryBoxShellScripts/wittypi3mini/install.sh"

# Directories for MPD (Music Player Daemon) libraries and binaries
mpd_lib_dir="$go_src_dir/StoryBox/lib/libmpdclient"
libmpdplaylistcontrol_dir="$go_src_dir/StoryBox/lib/mpdplaylistcontrol"
mpdcurrentsong_dir="$go_src_dir/StoryBox/lib"

# Directory for building the StoryBox Go binary
storybox_binary_dir="$go_src_dir/StoryBox/cmd/storybox"

# Directory to clone the StoryBox-Startup repository
storybox_startup_repo="$go_src_dir/StoryBox-Startup"

# Sound file configurations
sound_dir="/etc/sound"
started_mp3_file="$sound_dir/started.mp3"

# Systemd service file path
systemd_service_file="/lib/systemd/system/storyboxstartup.service"

# Source path for the startup sound file
started_mp3_source="$storybox_startup_repo/started.mp3"

sound_files_source="$go_src_dir/StoryBox/sys-audio/"

# Set up Go environment
print_heading 'Setting up Go environment...'
mkdir -p "$go_src_dir" "$go_pkg_dir" "$go_bin_dir"

# Download and install Go
print_heading 'Downloading and installing Go...'
wget "$go_tar_url" -P /tmp
tar -C /usr/local -xzf /tmp/$go_tarball
rm /tmp/$go_tarball

echo "export PATH=\$PATH:/usr/local/go/bin:/home/pi/go/bin" >> "$base_dir/.bashrc"
echo "export GOPATH=\$HOME/go" >> "$base_dir/.bashrc"
source "$base_dir/.bashrc"

# Clone the StoryBox repository
print_heading 'Cloning the StoryBox repository...'
git clone "$storybox_repo" "$go_src_dir/StoryBox"

# Setup .mpd directory
print_heading 'Setting up .mpd directory...'
mkdir -p "$base_dir/.mpd/playlists"
touch "$base_dir/.mpd/database"
touch "$base_dir/.mpd/log"
touch "$base_dir/.mpd/pid"
touch "$base_dir/.mpd/state"

# Copy MPD configuration
cp "$go_src_dir/StoryBox/lib/mpd.conf" "$base_dir/.mpd/mpd.conf"

# Set up phatbeat
print_heading 'Setting up phatbeat...'
git clone "$shell_scripts_repo" "$go_src_dir/StoryBoxShellScripts"
cd "$go_src_dir/StoryBoxShellScripts"
chmod +x phatbeat.sh
./phatbeat.sh

# Install Wi-Fi configuration for wittypi3mini
print_heading 'Installing Wi-Fi configuration for wittypi3mini...'


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
cd "$mpd_lib_dir"
meson setup build
ninja -C build
ninja -C build install

# Build and install libmpdplaylistcontrol
print_heading 'Building and installing libmpdplaylistcontrol...'
cd "$libmpdplaylistcontrol_dir"
gcc -c mpdplaylistcontrol.c -o mpdplaylistcontrol.o
ar rcs libmpdplaylistcontrol.a mpdplaylistcontrol.o
cp libmpdplaylistcontrol.a /usr/local/lib/
ldconfig

# Build mpdcurrentsong
print_heading 'Building mpdcurrentsong...'
cd "$mpdcurrentsong_dir"
gcc -o mpdcurrentsong mpdcurrentsong.c -lmpdclient

# Build mpdplaystate
print_heading 'Building mpdplaystate...'
gcc -o mpdplaystate mpdplaystate.c -lmpdclient

# Build mpdtime
print_heading 'Building mpdtime...'
gcc -o mpdtime mpdtime.c -lmpdclient

# Move the binaries to /usr/local/bin/
print_heading 'Moving binaries...'
cp "$mpdcurrentsong_dir/mpdcurrentsong" /usr/local/bin/
cp "$mpdcurrentsong_dir/mpdplaystate" /usr/local/bin/
cp "$mpdcurrentsong_dir/mpdtime" /usr/local/bin/

# Build StoryBox binary
print_heading 'Building StoryBox binary...'
cd "$storybox_binary_dir"
go build -o "$go_src_dir/StoryBox/StoryBox"

# Copy StoryBox binary
print_heading 'Copying StoryBox binary...'
mkdir -p /usr/local/bin/StoryBox
cp "$go_src_dir/StoryBox/StoryBox" /usr/local/bin/StoryBox/

# Clone the StoryBox-Startup repository
print_heading 'Cloning the StoryBox-Startup repository...'
git clone "https://github.com/$repo_owner/StoryBox-Startup.git" "$storybox_startup_repo"

# Build the Startup binary
print_heading 'Building the Startup binary...'
cd "$storybox_startup_repo"
go build -o "$go_src_dir/StoryBox-Startup/Startup"

# Copy the Startup binary
print_heading 'Copying the Startup binary...'
cp "$go_src_dir/StoryBox-Startup/Startup" /usr/local/bin/

# Set up systemd service
print_heading 'Setting up systemd service...'
cp "$storybox_startup_repo/storyboxstartup.service" "$systemd_service_file"
chmod 644 "$systemd_service_file"
systemctl enable storyboxstartup.service

# Set up sound files
print_heading 'Setting up sound files...'
mkdir -p "$sound_dir"

if [ ! -f "$started_mp3_file" ]; then
    print_heading 'Copying started.mp3 to /etc/sound/'
    cp "$started_mp3_source" "$started_mp3_file"
    print_heading 'Copy completed...'
else
    print_warning 'started.mp3 already exists. Skipping...'
fi

# Copy the rest of the sound files to /etc/sound/
cp -R "$sound_files_source" "$sound_dir/"

# Final message
print_success 'Installation completed successfully!'

print_warning 'Please review the script and make any necessary adjustments specific to your environment before executing it.'

# Start the storyboxstartup.service
systemctl start storyboxstartup.service