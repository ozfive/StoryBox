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

productname="StoryBox" # the name of the product to install
scriptname="install"   # the name of this script

sudo apt update
sudo apt upgrade
sudo apt install libmpdclient-dev gcc meson ninja-build golang python mpc mpd mpg123-alsa
pip install gTTS

wget https://go.dev/dl/go1.19.2.linux-arm64.tar.gz

# Untar the archive
tar -C /usr/local -xzf go1.19.2.linux-armv6l.tar.gz

# Remove the tar.gz file.
rm go1.19.2.linux-armv6l.tar.gz

# Set the PATH environment variable.
echo "export PATH=$PATH:/usr/local/go/bin" >>~/.profile

# Set the GOPATH environment variable.
echo "export GOPATH=$HOME/go" >>~/.profile

goVersion=$(go version | {
  read -r _ _ v _
  echo "${v#go}"
})

if $goVersion -lt "1.16.6"; then
  echo "You have chosen to install an older version of go which won't work with this project."
  exit 1
fi

git clone https://github.com/ozfive/StoryBox.git
