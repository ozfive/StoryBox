# StoryBox #

#### Summary ####

This repository contains the source code for my StoryBox project, which is a Raspberry Pi-based device that plays audio stories for children. The device uses an RFID reader to identify story/song/album cards, and then plays the corresponding audio through its built-in speaker. 

The StoryBox project includes several software components, including a startup script, a Go application for interacting with the RFID reader. The installation shell scripts reside in their own repository now. [StoryBoxShellScripts](https://github.com/ozfive/StoryBoxShellScripts)

### Images
Dev hardware I put together to work on this project:

<img src="https://github.com/ozfive/StoryBox/blob/main/github/Box-Front.jpg" alt=“Box-Front” width="415px" height="311">

<img src="https://github.com/ozfive/StoryBox/blob/main/github/Box-Internal.jpg" alt=“Box-Internal” width="415" height="553">
### How do I get set up? ###

Fork this repository.

Install git on your Rasperry Pi Zero. 
```shell
sudo apt install git
```
You must then clone your fork into /home/pi/ directory.

```shell
cd /home/pi

git clone git@github.com:[YOUR GIT USERNAME]/StoryBox.git
```
Make sure that you replace ```'[YOUR GIT USERNAME]'``` with your git username when you execute the above command.

Now move into the /home/pi/Storybox/ directory and make sure install.sh is executable.

```shell
cd Storybox/

chmod +x install.sh

./install.sh
```

When your Raspberry Pi Zero W has rebooted execute the following command to ensure the SPI interface was enabled
```shell
lsmod | grep spi
```

If you see spi_bcm2835, then you can proceed.


	
* Dependencies

REQUIRES:

	Go 1.19.2 or later
	libmpdclient-dev
	gcc
	meson
	ninja-build
	sqlite3
	python3
	mpc
	mpd
	mpg123
	libasound2-dev
	git

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
[@ozfive](https://github.com/ozfive)

* Other community or team contact

## License
This program is licensed under the [MIT License](https://opensource.org/license/mit/). See the LICENSE file for details.