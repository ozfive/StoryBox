#!/usr/bin/env python

import json
import requests
import RPi.GPIO as GPIO
from mfrc522 import SimpleMFRC522

# Initialize the RFID reader
reader = SimpleMFRC522()

# API endpoint URL
api_url_base = 'http://192.168.1.111:3001/'

# Function to start the playlist based on RFID tag information
def start_playlist(tag, unique):
    """
    Starts the playlist based on the RFID tag information.

    Args:
        tag: The RFID tag ID.
        unique: The unique ID associated with the tag.

    Returns:
        If the request is successful, the response JSON is returned.
        Otherwise, False is returned.
    """
    api_url = f'{api_url_base}rfid/'
    print(api_url)
    print(tag)
    print(unique)
    payload = {'tagid': tag, 'uniqueid': unique}

    try:
        headers = {'Content-Type': 'application/json'}
        response = requests.post(api_url, json=payload, headers=headers)

        if response.ok:
            return response.json()
        else:
            return False
    except requests.RequestException as e:
        print(f'Request failed: {e}')
        return False

try:
    while True:
        # Read RFID tag information
        ids, text = reader.read()
        playlist_info = start_playlist(ids, text)

        if playlist_info is not False:
            print("RFID was read.")
            for k, v in playlist_info.items():
                print(f'{k}: {v}')
        else:
            print('[!] Request Failed')

except KeyboardInterrupt:
    print("Program interrupted by user.")

except Exception as e:
    print(f'An error occurred: {e}')

finally:
    GPIO.cleanup()