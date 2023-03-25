# Introducing StoryBox

StoryBox is an engaging and entertaining device designed for children. Built on a Raspberry Pi foundation, StoryBox incorporates an RFID reader to play audio stories, songs, and playlists, offering a delightful listening experience without screens.

## Overview

The StoryBox repository encompasses the entire source code for this inventive hardware project centered on a Raspberry Pi. Designed to captivate children's imagination, StoryBox utilizes an RFID reader to recognize story, song, album, and playlist cards, and it features integrated speakers to deliver the associated audio content.

Alongside the hardware components, this project incorporates various software applications, including a startup application to inform users when the hardware is ready, and a Go-based server that facilitates communication with the RFID reader.

The hardware installation shell scripts reside in a separate repository: [StoryBoxShellScripts](https://github.com/ozfive/StoryBoxShellScripts)

## Gallery
Below are images of the development hardware assembled for this project:

<img src="https://github.com/ozfive/StoryBox/blob/main/github/Box-Front.jpg" alt=“Box-Front” width="415px" height="311">

<img src="https://github.com/ozfive/StoryBox/blob/main/github/Box-Internal.jpg" alt=“Box-Internal” width="415" height="553">

## Setup Guide

To set up the project on your Raspberry Pi Zero, please follow these steps:

1. Fork the repository.
2. Install git on your Raspberry Pi Zero by executing this command in your terminal:

```shell
sudo apt install git
```

3. Clone your fork into the directory /home/pi/ by executing the following commands in your terminal. Please make sure to replace `[GIT_USER]` with your own git username.:

```shell
cd /home/pi
git clone git@github.com:[GIT_USER]/StoryBox.git
```

4. Access the /home/pi/Storybox/ directory and ensure that install.sh is executable by running these commands in your terminal:

```shell
cd Storybox/

chmod +x install.sh

./install.sh
```

5. After your Raspberry Pi Zero reboots, run this command to verify that the SPI interface is enabled:

```shell
lsmod | grep spi
```

If you see spi_bcm2835, you can proceed.
	
## Dependencies

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

## Database configuration

	The project depends on a SQLite3 database named rfids.db.

## Who To Contact

* For any inquiries, please reach out to the repo owner [@ozfive](https://github.com/ozfive)

## License
This program is licensed under the [MIT License](https://opensource.org/license/mit/). See the LICENSE file for details.