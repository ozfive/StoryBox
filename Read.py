#!/usr/bin/env python

import json
import requests
import RPi.GPIO as GPIO
from mfrc522 import SimpleMFRC522

reader = SimpleMFRC522()

api_url_base = 'http://192.168.1.111:3001/'

def start_playlist(tag, unique):

    api_url = '{0}rfid/'.format(api_url_base)
    print(api_url)
    print(tag)
    print(unique)
    payload = {'tagid':tag, 'uniqueid':unique}
    response = requests.post(api_url, data = json.dumps(payload))

    if response.status_code == 200:
        return json.loads(response.read().decode('utf-8'))
    else:
        return False


try:
    while True:
        ids, text = reader.read()

        playlist_info = start_playlist(ids, text)

        if playlist_info is not False:

            print("RFID was read.")

    for k, v in playlist_info.items():
        print('{0}:{1}'.format(k, v))

    else:
        print('[!] Request Failed')

finally:
    GPIO.cleanup()
