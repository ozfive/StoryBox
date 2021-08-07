#!/usr/bin/env python
import sqlite3
import secrets

# import RPi.GPIO as GPIO
# from mfrc522 import SimpleMFRC522

urlSafeKey = secrets.token_urlsafe()

# reader = SimpleMFRC522()
conn = sqlite3.connect('rfids.db')
c = conn.cursor()

try:

	print("Please tap the RFID Tag to read the ID.")

	ids, text = reader.read()
	
	print("Success!")
	print("Checking the database for ID.")

	# Checking the database to see if the tag is there.
	c.execute('SELECT * FROM rfid WHERE tagid=?;', (ids,))
	returnedRecord = c.fetchone()
	
	playlistName = input('Playlist Name:') # Requesting an playlist name from the command line.
	url = input('URL: ')

	# print(urlSafeKey)

	print("Place your tag to write key")

	reader.write(urlSafeKey)
	print("Written")
	
	# If the returnedRecord tagid value is blank then do an insert statement.
	# Else do an update since the record already exists.
	if returnedRecord is None:
	# Set the record with all values needed to execute the insert statement.
		record = [(None, ids, urlSafeKey, url, playlistName),]
		# Execute the insert statment.
		c.executemany('INSERT INTO rfid VALUES (?,?,?,?,?)', record)
	else:
		print("Tag Exists, updating Record.")
		c.execute("UPDATE rfid SET uniqueid = ?, url = ?, playlistname = ? WHERE tagid = ?;", (urlSafeKey, url, playlistName, ids))

	conn.commit()

	conn.close()

finally:
    GPIO.cleanup()