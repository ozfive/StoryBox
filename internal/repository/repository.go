package repository

import (
	"StoryBox/internal/models"
	"database/sql"
	"fmt"
)

type RFIDRepository interface {
	Create(rfid *models.RFID) error
	GetByTagAndUniqueID(tagID, uniqueID string) (*models.RFID, error)
}

type PlaylistRepository interface {
	Create(url, playlistName string) error
	Delete(url, playlistName string) error
	Get(url, playlistName string) (*models.Playlist, error)
}

type repository struct {
	db *sql.DB
}

func NewRFIDRepository(db *sql.DB) RFIDRepository {
	return &repository{db: db}
}

func NewPlaylistRepository(db *sql.DB) PlaylistRepository {
	return &repository{db: db}
}

func ConnectDatabase(databasePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open SQLite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("unable to connect to SQLite database: %w", err)
	}

	// Initialize tables
	if err := initializeTables(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func initializeTables(db *sql.DB) error {
	rfidTable := `
    CREATE TABLE IF NOT EXISTS rfid (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        tagid TEXT,
        uniqueid TEXT,
        url TEXT,
        playlistname TEXT
    );`

	playlistTable := `
    CREATE TABLE IF NOT EXISTS playlist (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        url TEXT,
        playlistname TEXT
    );`

	_, err := db.Exec(rfidTable)
	if err != nil {
		return fmt.Errorf("failed to create rfid table: %w", err)
	}

	_, err = db.Exec(playlistTable)
	if err != nil {
		return fmt.Errorf("failed to create playlist table: %w", err)
	}

	return nil
}

// Implement RFIDRepository methods
func (r *repository) Create(rfid *models.RFID) error {
	stmt, err := r.db.Prepare("INSERT INTO rfid (tagid, uniqueid, url, playlistname) VALUES (?, ?, ?, ?);")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(rfid.TagID, rfid.UniqueID, rfid.URL, rfid.PlaylistName)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	rfid.ID = int(id)
	return nil
}

func (r *repository) GetByTagAndUniqueID(tagID, uniqueID string) (*models.RFID, error) {
	stmt, err := r.db.Prepare("SELECT id, tagid, uniqueid, url, playlistname FROM rfid WHERE tagid = ? AND uniqueid = ?;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rfid := &models.RFID{}
	err = stmt.QueryRow(tagID, uniqueID).Scan(&rfid.ID, &rfid.TagID, &rfid.UniqueID, &rfid.URL, &rfid.PlaylistName)
	if err != nil {
		return nil, err
	}

	return rfid, nil
}

// Implement PlaylistRepository methods similarly...
