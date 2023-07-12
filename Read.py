#!/usr/bin/env python

import json
import requests
import socket
import RPi.GPIO as GPIO
from mfrc522 import SimpleMFRC522

# Initialize the RFID reader
reader = SimpleMFRC522()

def get_outbound_ip():
    connection = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    connection.connect(("8.8.8.8", 80))

    local_ip = connection.getsockname()[0]
    connection.close()

    return local_ip


# Get the outbound IP dynamically
outbound_ip = get_outbound_ip()
port = '3001'
# Construct the API base URL with the dynamic IP
api_url_base = f"http://{outbound_ip}:{port}/"

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

        # Print the RFID tag information
        print(f'Tag ID: {ids}')
        print(f'Unique ID: {text}')

        # Start the playlist based on RFID tag information
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