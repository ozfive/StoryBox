# README #

### What is this repository for? ###

* Quick summary

* Version

### How do I get set up? ###

* Prerequisites
		
	* libmpdclient-dev
	* A C99 compliant compiler (e.g. gcc) - GCC Comes pre-installed on RPi-Zero-W
	* Meson 0.37
	* Ninja
	* Go 1.16 or later
	* python - Comes pre-installed on RPi-Zero-W
	* gtts-cli
	* mpg123-alsa - Comes pre-installed on RPi-Zero-W
	* mpc
	* mpd

		// Always update your system before adding something new.
		sudo apt update
		sudo apt upgrade

		Download go 1.16.6 by using wget. If by the time you are reading this  `
		there is a newer go compiler available you can try and download it from

		https://golang.org/dl/

		NOTE: It needs to be the linux arm v6 compiler to build the go applications on
		the Raspberry Pi. using 'sudo apt-get install golang' will not get you the
		latest version of the Go compiler and it may not compile this 

		```shell
		wget https://golang.org/dl/go1.16.6.linux-armv6l.tar.gz
		
		sudo tar -C /usr/local -xzf go1.16.6.linux-armv6l.tar.gz

		rm go1.16.6.linux-armv6l.tar.gz

		sudo nano ~/.bashrc
		```

		add the following to the .bashrc file

		PATH=$PATH:/usr/local/go/bin
		GOPATH=$HOME/pi/go/

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

	Go 1.15 or later
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