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
apt install -y libmpdclient-dev gcc meson ninja-build golang python3 mpc mpd mpg123 libasound2-dev git

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
mv mpdcurrentsong ~/StoryBox/bin/

gcc -o mpdplaystate mpdplaystate.c -lmpdclient
mv mpdplaystate ~/StoryBox/bin/

gcc -o mpdtime mpdtime.c -lmpdclient
mv mpdtime ~/StoryBox/bin/

# Build and move StoryBox
cd ~/StoryBox/ || exit

go build -o StoryBox

mv StoryBox ~/StoryBox/bin/
