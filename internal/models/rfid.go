package models

type RFID struct {
	ID           int    `json:"id"`
	TagID        string `json:"tagid"`
	UniqueID     string `json:"uniqueid"`
	URL          string `json:"url"`
	PlaylistName string `json:"playlistname"`
}

type Playlist struct {
	ID             int
	URLFromDB      string
	PlaylistNameDB string
	PlaylistName   string
	URL            string
}
