# README #

### What is this repository for? ###

* Quick summary

* Version

### How do I get set up? ###

* Prerequisites
		
	* lmpdclient-dev
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

		```bash
		wget https://golang.org/dl/go1.16.6.linux-armv6l.tar.gz
		
		sudo tar -C /usr/local -xzf go1.16.6.linux-armv6l.tar.gz

		rm go1.16.6.linux-armv6l.tar.gz

		sudo nano ~/.profile
		```

		add the following to the .profile file

		PATH=$PATH:/usr/local/go/bin
		GOPATH=$HOME/go

		Update the shell with the changes.
		```bash
		source ~/.profile
		```

		Check to make sure that the go version is the correct one.
		```bash
		go version
		```
		Git clone this repository
		```bash
		git clone https://cowboysteeve@bitbucket.org/cowboysteeve/storybox.git
		```
		
		Install lmpdclient-dev package for the c language applications to function.

		```bash
		sudo apt install libmpdclient-dev
		
		sudo apt install git

		git clone https://github.com/MusicPlayerDaemon/libmpdclient.git
		
		cd libmpdclient/
		
		sudo apt install meson

		sudo apt install ninja-build

		sudo meson . output
		
		sudo ninja -C output
		
		sudo ninja -C output install
		```

		The previous commands make it possible to interface with the mpd player 
		through the c applications we will now compile with gcc.




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