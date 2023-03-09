# StoryBox #

#### Summary ####

This repository contains the source code for my StoryBox project, which is a Raspberry Pi-based device that plays audio stories for children. The device uses an RFID reader to identify story/song/album cards, and then plays the corresponding audio through its built-in speaker. 

The StoryBox project includes several software components, including a startup script, a Go application for interacting with the RFID reader, and installation bash scripts to install drivers for the hardware.

### Images
Dev hardware I put together to work on this project:

![Box-Front](https://github.com/ozfive/StoryBox/blob/main/github/Box-Front.jpg)

![Box-Internal](https://github.com/ozfive/StoryBox/blob/main/github/Box-Internal.jpg)

### How do I get set up? ###

#### Prerequisites ####
		
* libmpdclient-dev
* A C99 compliant compiler (e.g. gcc)
* Meson 0.37
* Ninja
* Go 1.16 or later
* python - Comes pre-installed on RPi-Zero-W
* gtts-cli
* mpg123-alsa - Comes pre-installed on RPi-Zero-W
* mpc
* mpd

Update the system before you do anything else.

```shell
sudo apt update
sudo apt upgrade
```
---
Install the required packages.
```shell
sudo apt install libmpdclient-dev gcc meson ninja-build golang python mpc mpd 
mpg123-alsa
```
---
Using pip, install the gTTS (Google Text To Speech) library.
```shell
pip install gTTS
```
---
Note that the latest version of the Go compiler may not be available in the apt repository, so you may need to download it manually from the [official website](https://golang.org/dl/). When downloading the Go compiler, be sure to choose the linux arm v6 version to ensure that it will work on the Raspberry Pi.

Steps to install Go on the Raspberry Pi:

```shell
wget https://golang.org/dl/go1.16.6.linux-armv6l.tar.gz

sudo tar -C /usr/local -xzf go1.16.6.linux-armv6l.tar.gz

rm go1.16.6.linux-armv6l.tar.gz

sudo nano ~/.bashrc
```

add the following to the .bashrc file
```shell
PATH=$PATH:/usr/local/go/bin
GOPATH=$HOME/pi/go/
```
Update the shell with the changes.
```shell
source ~/.profile
```

Check to make sure that the go version is the correct one.
```shell
go version
```
Clone this repository
```shell
git clone https://github.com/ozfive/StoryBox.git
```

Install the phatbeat software to be able to control mpd eventually
NOTE: Be sure to run these commands as the pi user and not su
```shell
cd /home/pi/Storybox/

chmod +x phatbeat.sh

./phatbeat.sh
```
When you run this script it will ask you if you want to continue
type 'y' for yes only after you feel comfortable that the script
doesn't include malicious code that can arbitrarily execute on 
your machine.

DISCLAIMER: I have done my best to review this script. It does
reach out to other sources to download things. At any point in
the future these sources could be compromised. I will not be 
held responsible in any way, shape, or form for any damage that
is done to your systems/network due to this script. 

When asked if you want to perform a fill install you can either choose
'y' or 'N' here. It's up to you if you would like to check out the
examples for the PhatBeat pHat.

Now sit back and drink some coffee...

When the script is finished executing it will ask you if you would
like to reboot now since changes were made to the system that
require a reboot. type 'y' and hit enter here to reboot.

This will close your ssh session to the Raspberry Pi Zero W and
you will need to connect again to continue with the instructions
below.

Refer to https://www.uugear.com/product/witty-pi-3-mini-realtime-clock-and-power-management-for-raspberry-pi/ for more information.

Install witty pi 3 mini software.
```shell
wget http://www.uugear.com/repo/WittyPi3/install.sh

sudo sh install.sh

```
In case this doesn't exist at that link you can just use
the install.sh script included in lib/Wittypi3mini
NOTE: Only execute the commands in the next block if you can
no longer retrieve the install.sh script for the Witty Pi 3
Mini hardware from the link above.
```shell
cd /home/pi/StoryBox/MainSystem/lib/

sudo sh install.sh
```

Now we need to enable the SPI interface for the RFID-RC522 reader
to function properly. To do this we will open the GUI raspi-config
tool by executing the following command.
```shell
sudo raspi-config
```

Use the arrow keys to select "5 Interfacing Options". Once you 
have this option selected, press Enter.

On this next screen, you want to use your arrow keys to select 
"P4 SPI", again press Enter to select the option once it is 
highlighted.

You will now be asked if you want to enable the SPI Interface, 
select Yes with your arrow keys and press Enter to proceed. 

Once the SPI interface has been successfully enabled by the 
raspi-config tool you should see the following text appear 
on the screen, "The SPI interface is enabled".

To fully enable the SPI interface please reboot your Raspberry Pi 
Zero W
```shell
reboot -n
```

When your Raspberry Pi Zero W has rebooted execute the following command
```shell
lsmod | grep spi
```

If you see spi_bcm2835, then you can proceed.

Install Google Text To Speech CLI
```shell
pip install gTTS
```


Install MPD and MPC 

```shell

mkdir ~/.mpd/
mkdir ~/music
mkdir ~/.mpd/playlists
touch ~/.mpd/database
touch ~/.mpd/log
touch ~/.mpd/pid

mv /home/pi/StoryBox/MainSystem/lib/mpd.conf ~/.mpd/mpd.conf

sudo apt install mpd

sudo apt install mpc
```
Install libmpdclient-dev package for the c language applications to function.

```shell
sudo apt install libmpdclient-dev

sudo apt install git

git clone https://github.com/MusicPlayerDaemon/libmpdclient.git

cd libmpdclient/

sudo apt install meson

sudo apt install ninja-build

sudo meson . output

sudo ninja -C output

sudo ninja -C output install

cd ..
```

The previous commands make it possible to interface with the mpd player 
through the c applications we will now compile with gcc.

```shell
cd StoryBox/MainSystem/lib

gcc -o mpdcurrentsong mpdcurrentsong.c -lmpdclient

mv mpdcurrentsong /home/pi/StoryBox/MainSystem/bin/

gcc -o mpdplaystate mpdplaystate.c -lmpdclient

mv mpdplaystate /home/pi/StoryBox/MainSystem/bin/

gcc -o mpdtime mpdtime.c -lmpdclient

mv mpdtime /home/pi/StoryBox/MainSystem/bin/

cd /home/pi/StoryBox/MainSystem/bin/
```

We have now compiled the c applications that the project uses to 
interface with mpd through lmpdclient library.


* Summary of set up

* Configuration
	
	STEP 1: Once the storybox repo is cloned on your Raspberry Pi go into the Startup folder and run 

	```shell

	mkdir /home/pi/go/
	mkdir /home/pi/go/src
	
	cp -r /home/pi/StoryBox /home/pi/go/src

	cd /home/pi/go/src/Storybox/Startup
	
	go mod init Startup

	go mod tidy

	go build
	```

	STEP 2: Copy the binary you just built to /usr/local/bin by typing 

	```shell
	sudo cp Startup /usr/local/bin
	```

	This makes the Startup binary available to all users on the system. If for some reason you want to
	execute the Startup binary through another user it will be available system wide.

	STEP 3: Copy the storyboxstartup.service file to lib/systemd/system/ by going into the storybox 
	directory where the storyboxstartup.service file is located and typing: 

	```shell
	cd /home/pi/go/src/StoryBox/

	sudo cp storyboxstartup.service /lib/systemd/system/ 
	```

	The next command sets the permissions of the unit file to 644 

	```shell
	sudo chmod 644 /lib/systemd/system/storyboxstartup.service
	```

	Now we want to enable the storyboxstartup.service unit file
	to be started at startup. 
	
	```shell
	sudo systemctl enable storyboxstartup.service

	sudo systemctl start storyboxstartup.service
	```


	

	
* Dependencies

REQUIRES:

	Go 1.16 or later
	python
	gtts-cli
	mpg123-alsa
	mpc
	mpd

* Database configuration

	This project relies on sqlite3 for python scripts and go applications.

* How to run tests

* Deployment instructions

### Contribution guidelines ###

* Writing tests

* Code review

* Other guidelines

### Who do I talk to? ###

* Repo owner or admin

* Other community or team contact