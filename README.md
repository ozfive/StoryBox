# README #

### What is this repository for? ###

* Quick summary

* Version

### How do I get set up? ###

* Prerequisites
		
	* libmpdclient-dev
	* A C99 compliant compiler (e.g. gcc)
	* Meson 0.37
	* Ninja
	* Go 1.16 or later
	* python
	* gtts-cli
	* mpg123-alsa
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

		sudo nano ~/.profile
		```

		add the following to the .profile file

		PATH=$PATH:/usr/local/go/bin
		GOPATH=$HOME/go

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


		Install MPD and MPC

		```shell

		mkdir ~/.mpd/

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

		We have now compiled the c applications that the project uses to interface with mpd through l

* Summary of set up

* Configuration
	
	STEP 1: Once the storybox repo is cloned on your Raspberry Pi go into the Startup folder and run 

	```
	go build
	```

	STEP 2: Move the binary you just built to /usr/local/bin by typing 

	```
	sudo mv Startup /usr/local/bin
	```

	This makes the Startup binary available to all users on the system. If for some reason you want to
	execute the Startup binary through another user it will be available system wide.

	STEP 3: Move the storyboxstartup.service file to lib/systemd/system/ by going into the storybox 
	directory where the storyboxstartup.service file is located and typing 

	```
	sudo mv storyboxstartup.service /lib/systemd/system/ 
	```
	This should be followed by:

	```
	sudo chmod 644 /lib/systemd/system/storyboxstartup.service
	```

	What this does is sets the permissions of the unit file to '644'

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