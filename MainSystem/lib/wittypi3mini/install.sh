[ -z $BASH ] && { exec bash "$0" "$@" || exit; }
#!/bin/bash
# file: installWittyPi.sh
#
# This script will install required software for Witty Pi.
# It is recommended to run it in your home directory.
#

# check if sudo is used
if [ "$(id -u)" != 0 ]; then
  echo 'Sorry, you need to run this script with sudo'
  exit 1
fi

# target directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/wittypi"

# error counter
ERR=0

echo '================================================================================'
echo '|                                                                              |'
echo '|                   Witty Pi Software Installation Script                      |'
echo '|                                                                              |'
echo '================================================================================'

# enable I2C on Raspberry Pi
echo '>>> Enable I2C'
if grep -q 'i2c-bcm2708' /etc/modules; then
  echo 'Seems i2c-bcm2708 module already exists, skip this step.'
else
  echo 'i2c-bcm2708' >> /etc/modules
fi
if grep -q 'i2c-dev' /etc/modules; then
  echo 'Seems i2c-dev module already exists, skip this step.'
else
  echo 'i2c-dev' >> /etc/modules
fi

i2c1=$(grep 'dtparam=i2c1=on' /boot/config.txt)
i2c1=$(echo -e "$i2c1" | sed -e 's/^[[:space:]]*//')
if [[ -z "$i2c1" || "$i2c1" == "#"* ]]; then
  echo 'dtparam=i2c1=on' >> /boot/config.txt
else
  echo 'Seems i2c1 parameter already set, skip this step.'
fi

i2c_arm=$(grep 'dtparam=i2c_arm=on' /boot/config.txt)
i2c_arm=$(echo -e "$i2c_arm" | sed -e 's/^[[:space:]]*//')
if [[ -z "$i2c_arm" || "$i2c_arm" == "#"* ]]; then
  echo 'dtparam=i2c_arm=on' >> /boot/config.txt
else
  echo 'Seems i2c_arm parameter already set, skip this step.'
fi

miniuart=$(grep 'dtoverlay=pi3-miniuart-bt' /boot/config.txt)
miniuart=$(echo -e "$miniuart" | sed -e 's/^[[:space:]]*//')
if [[ -z "$miniuart" || "$miniuart" == "#"* ]]; then
  echo 'dtoverlay=pi3-miniuart-bt' >> /boot/config.txt
else
  echo 'Seems setting Pi3/4 Bluetooth to use mini-UART is done already, skip this step.'
fi

core_freq=$(grep 'core_freq=250' /boot/config.txt)
core_freq=$(echo -e "$core_freq" | sed -e 's/^[[:space:]]*//')
if [[ -z "$core_freq" || "$core_freq" == "#"* ]]; then
  echo 'core_freq=250' >> /boot/config.txt
else
  echo 'Seems the frequency of GPU processor core is set to 250MHz already, skip this step.'
fi

if [ -f /etc/modprobe.d/raspi-blacklist.conf ]; then
  sed -i 's/^blacklist spi-bcm2708/#blacklist spi-bcm2708/' /etc/modprobe.d/raspi-blacklist.conf
  sed -i 's/^blacklist i2c-bcm2708/#blacklist i2c-bcm2708/' /etc/modprobe.d/raspi-blacklist.conf
else
  echo 'File raspi-blacklist.conf does not exist, skip this step.'
fi

# install i2c-tools
echo '>>> Install i2c-tools'
if hash i2cget 2>/dev/null; then
  echo 'Seems i2c-tools is installed already, skip this step.'
else
  apt-get install -y i2c-tools || ((ERR++))
fi

# make sure en_GB.UTF-8 locale is installed
echo '>>> Make sure en_GB.UTF-8 locale is installed'
locale_commentout=$(sed -n 's/\(#\).*en_GB.UTF-8 UTF-8/1/p' /etc/locale.gen)
if [[ $locale_commentout -ne 1 ]]; then
	echo 'Seems en_GB.UTF-8 locale has been installed, skip this step.'
else
	sed -i.bak 's/^.*\(en_GB.UTF-8[[:blank:]]\+UTF-8\)/\1/' /etc/locale.gen
	locale-gen
fi

# check if it is Raspberry Pi 4
isRpi4=$(cat /proc/device-tree/model | sed -n 's/.*\(Raspberry Pi 4\).*/1/p')

# install wiringPi
if [ $ERR -eq 0 ]; then
  echo '>>> Install wiringPi'
  ver=0;
  if hash gpio 2>/dev/null; then
  	ver=$(gpio -v | sed -n '1 s/.*\([0-9]\+\.[0-9]\+\).*/\1/p')
  	echo "wiringPi version: $ver"
 	else
 		apt-get -y install wiringpi
 		ver=$(gpio -v | sed -n '1 s/.*\([0-9]\+\.[0-9]\+\).*/\1/p')
  fi
	if [[ $isRpi4 -eq 1 ]] && (( $(awk "BEGIN {print ($ver < 2.52)}") )); then
 		wget https://project-downloads.drogon.net/wiringpi-latest.deb || ((ERR++))
		dpkg -i wiringpi-latest.deb || ((ERR++))
		rm wiringpi-latest.deb
  fi
fi

# install wittyPi
if [ $ERR -eq 0 ]; then
  echo '>>> Install wittypi'
  if [ -d "wittypi" ]; then
    echo 'Seems wittypi is installed already, skip this step.'
  else
    wget https://www.uugear.com/repo/WittyPi3/LATEST -O wittyPi.zip || ((ERR++))
    unzip wittyPi.zip -d wittypi || ((ERR++))
    cd wittypi
    chmod +x wittyPi.sh
    chmod +x daemon.sh
    chmod +x syncTime.sh
    chmod +x runScript.sh
    chmod +x afterStartup.sh
    chmod +x beforeShutdown.sh
    sed -e "s#/home/pi/wittypi#$DIR#g" init.sh >/etc/init.d/wittypi
    chmod +x /etc/init.d/wittypi
    update-rc.d wittypi defaults || ((ERR++))
    touch wittyPi.log
    touch schedule.log
    cd ..
    chown -R $SUDO_USER:$(id -g -n $SUDO_USER) wittypi || ((ERR++))
    sleep 2
    rm wittyPi.zip
  fi
fi

# install UUGear Web Interface
curl https://www.uugear.com/repo/UWI/installUWI.sh | bash

echo
if [ $ERR -eq 0 ]; then
  echo '>>> All done. Please reboot your Pi :-)'
else
  echo '>>> Something went wrong. Please check the messages above :-('
fi
