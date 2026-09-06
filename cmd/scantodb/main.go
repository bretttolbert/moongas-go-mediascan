package main

// Unlike the yaml scanners, this one scans everything in one cmd,
// i.e. both media files and artist.yaml files,
// and outputs a single sqlite database with multiple tables
// (currently the tables are "mediafile" and "artist")

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library

	"github.com/bretttolbert/moongas-mediascan-go/internal/shared"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Error: Invalid arguments")
		fmt.Println("Usage: go run cmd/scantodb/main.go {config yaml filepath} {output database filepath}")
		fmt.Println("Example: go run cmd/scantodb/main.go mediascan-config.yaml mediascan.db")
		os.Exit(1)
	}
	configYamlFilepath := os.Args[1]
	outputDBFilepath := os.Args[2]
	conf := shared.LoadConf(configYamlFilepath)


	outputDBAbsFilepath, err := filepath.Abs(outputDBFilepath)
	if err != nil {
		log.Fatalf("Failed to get absolute path to output file: %v", err)
	}


	var files shared.MediaFiles = shared.ScanFiles(conf)
	var artists shared.Artists = shared.ScanArtists(conf)

	os.Remove(outputDBFilepath)

	log.Printf("Creating %s...", outputDBFilepath)
	file, err := os.Create(outputDBFilepath)
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Printf("%s created", outputDBFilepath)

	sqliteDatabase, _ := sql.Open("sqlite3", outputDBFilepath)
	defer sqliteDatabase.Close()
	createTableMediaFile(sqliteDatabase, files)
	createTableArtist(sqliteDatabase, artists)

	for i := 0; i < len(files.Files); i++ {
		mf := files.Files[i]
		insertMediaFile(sqliteDatabase, mf.Path, mf.AlbumPath, mf.ArtistPath, mf.Title, mf.Artist, mf.AlbumArtist, mf.Album, mf.Genre, mf.Year)
	}

	for i := 0; i < len(artists.Artists); i++ {
		a := artists.Artists[i]
		insertArtist(sqliteDatabase, a.Path, a.ArtistData.ArtistNames[0], a.ArtistData.City, a.ArtistData.CountryCode, a.ArtistData.RegionCode, a.ArtistData.LanguageCodes[0])
	}

	// DISPLAY INSERTED RECORDS
	displayMediaFiles(sqliteDatabase)
	displayArtists(sqliteDatabase)

	log.Printf("Written to file %s", outputDBAbsFilepath)
}

func createTableMediaFile(db *sql.DB, files shared.MediaFiles) {
	sql := `CREATE TABLE mediafile (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"path" TEXT UNIQUE,
		"albumpath" TEXT,
		"artistpath" TEXT,
		"title" TEXT,
		"artist" TEXT,
		"albumartist" TEXT,
		"album" TEXT,
		"genre" TEXT,
		"year" INTEGER
	  );` // SQL Statement for Create Table
	// Keep in sync w/ below: path, albumpath, artistpath, title, artist, albumartist, album, genre, year

	log.Println("Create mediafiles table...")
	statement, err := db.Prepare(sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln(err.Error())
	}

	// In SQLite, non-unique indexes are created using the standard CREATE INDEX statement,
	// which must be executed after the CREATE TABLE statement.
	sql = `CREATE INDEX idx_artistpath ON mediafile (artistpath);`
	statement, err = db.Prepare(sql)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln(err.Error())
	}

	sql = `CREATE INDEX idx_albumpath ON mediafile (albumpath);`
	statement, err = db.Prepare(sql)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("mediafiles table created")
}

func insertMediaFile(db *sql.DB, path string, albumpath string, artistpath string, title string, artist string, albumartist string, album string, genre string, year int) {
	log.Println("Inserting mediafile record ...")
	sql := `INSERT INTO mediafile(path, albumpath, artistpath, title, artist, albumartist, album, genre, year) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(sql)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(path, albumpath, artistpath, title, artist, albumartist, album, genre, year)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func createTableArtist(db *sql.DB, files shared.Artists) {
	sql := `CREATE TABLE artist (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"path" TEXT UNIQUE,
		"name" TEXT,
		"city" TEXT,
		"countrycode" TEXT,
		"regioncode" TEXT,
		"languagecode" TEXT
	  );` // SQL Statement for Create Table

	log.Println("Create artist table...")
	statement, err := db.Prepare(sql) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("artist table created")
}

func insertArtist(db *sql.DB, path string, name string, city string, countrycode string, regioncode string, languagecode string) {
	log.Println("Inserting artist record ...")
	insertStudentSQL := `INSERT INTO artist(path, name, city, countrycode, regioncode, languagecode) VALUES (?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(path, name, city, countrycode, regioncode, languagecode)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func displayMediaFiles(db *sql.DB) {
	row, err := db.Query("SELECT * FROM mediafile ORDER BY artist")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		var id int
		var path string
		var albumpath string
		var artistpath string
		var title string
		var artist string
		var albumartist string
		var album string
		var genre string
		var year int
		row.Scan(&id, &path, &albumpath, &artistpath, &title, &artist, &albumartist, &album, &genre, &year)
		log.Println("mediafile: ", id, " ", path, " ", albumpath, " ", artistpath, " ", title, " ", artist, " ", albumartist, " ", album, " ", genre, " ", year)

	}
}

func displayArtists(db *sql.DB) {
	row, err := db.Query("SELECT * FROM artist ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var path string
		var name string
		var city string
		var countrycode string
		var regioncode string
		var languagecode string
		row.Scan(&id, &path, &name, &city, &countrycode, &regioncode, &languagecode)
		log.Println("artist: ", id, " ", path, " ", name, " ", city, " ", countrycode, " ", regioncode, " ", languagecode)
	}
}
