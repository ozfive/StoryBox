# StoryBox

StoryBox is Raspberry Pi-based device that uses an RFID reader to play audio stories, songs, and playlists for children.

#### Summary

This repository contains the complete source code for StoryBox, a hardware project built around a Raspberry Pi that aims to provide children with a unique listening experience without screens. StoryBox features an RFID reader to identify story, song, album, and playlist cards, and integrated speakers to play the corresponding audio files.

In addition to the hardware components, the project includes a variety of software applications. A startup application provides feedback to the user when the hardware is ready, and a server written in Go enables communication with the RFID reader.

The installation shell scripts are housed separately in their own dedicated repository.

[StoryBoxShellScripts](https://github.com/ozfive/StoryBoxShellScripts)

## Images
The development hardware that I have assembled to undertake this project is shown in the images below:

<img src="https://github.com/ozfive/StoryBox/blob/main/github/Box-Front.jpg" alt=“Box-Front” width="415px" height="311">

<img src="https://github.com/ozfive/StoryBox/blob/main/github/Box-Internal.jpg" alt=“Box-Internal” width="415" height="553">

## How do I get set up?

Please follow these steps to set up the project on your Raspberry Pi Zero:

1. Begin by forking the repository.
2. Install git on your Raspberry Pi Zero by running the following command in your terminal:

```shell
sudo apt install git
```

3. Clone your fork into the directory /home/pi/ by executing the following commands in your terminal. Please ensure to replace `[YOUR GIT USERNAME]` with your own git username.:

```shell
cd /home/pi
git clone git@github.com:[YOUR GIT USERNAME]/StoryBox.git
```

4. Navigate to the /home/pi/Storybox/ directory and make sure that install.sh is executable by running the following commands in your terminal:

```shell
cd Storybox/

chmod +x install.sh

./install.sh
```

5. Once your Raspberry Pi Zero has rebooted, execute the following command in your terminal to ensure that the SPI interface is enabled:

```shell
lsmod | grep spi
```

If you see spi_bcm2835, then you can proceed.
	
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

	This project relies on a sqlite3 database called rfids.db.

## Who To Contact

* Repo owner or admin [@ozfive](https://github.com/ozfive)

## License
This program is licensed under the [MIT License](https://opensource.org/license/mit/). See the LICENSE file for details.