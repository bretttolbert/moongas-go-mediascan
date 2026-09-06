# moongas-mediascan-go

Simple and fast Golang command-line utility for scanning media library files and metadata files, extracting metadata (e.g. ID3v2 tags) and saving it into an SQLite database. 

- A component of the `moongas` ecosystem of media library tools.
- Suports reading metadata `moongas` metadata file formats such as `artist.yaml` files
- Output format (SQLite database) is compatible with `moongas-mediascan-py` python package

## Database Schema

The output database contains two tables:

| Table | Contents |
| ----- | -------- |
| mediafile | Track metadata from individual media files (ID3 tag information read from `mp3` files) |
| artist | Metadata from `artist.yaml` files |

## Usage

### Scan to Database

- [`cmd/scantodb/main.go`](./cmd/scantodb/main.go) - Scans media libraries for both mediafiles and moongas `artist.yaml` files, outputs an SQLite database (`.db`) file
```bash
go run cmd/scantodb/main.go mediascan-config.yaml mediascan.db
```

### Scan to Yaml (deprecated)

These commands scan media files and `artist.yaml` files and output massive combined yaml files. This was the original implementation before switching to sqlite. These may be removed in the future.

- [`cmd/scanfilesyaml/main.go`](./cmd/scanfilesyaml/main.go) - Recursively scan a directory for media files, extract metadata (including ID3v2 tags from both MP3 and M4A files), and save the output in a mediafiles YAML file (`files.yaml`). 
- Reads configuration from YAML file e.g. [mediascan-config.yaml](./mediascan-config.yaml)
- Has only two required command-line arguments: `{config yaml filepath}` and `{output yaml filepath}`
- Created specifically to run fast on a Raspberry Pi single-board computer as part of another project of mine.
- Usage:
```bash
go run cmd/scanfilesyaml/main.go mediascan-config.yaml files.yaml
```
- [`cmd/scanartistsyaml/main.go`](./cmd/scanartistsyaml/main.go) - Scans artist directories for `artist.yaml` files and aggregates them into a combined artists YAML file (`artists.yaml`).
- Usage:
```bash
go run cmd/scanartistsyaml/main.go mediascan-config.yaml artists.yaml
```

## mediascan.go YAML Configuration file Reference

The mediascan.go YAML configuration file (example: [conf.yaml](./mediascan-config.yaml)) supports the following parameters:

| Property | Description |
| -------- | ----------- |
| `mediadir` | the path to the directory to be scanned |
| `mediaexts` | the file extensions to be included in the scan |
| `excludepath` | path substrings to exclude from scan (case-insensitive) |
| `excludetitle` | id3 title substrings to exclude from scan results (case-insensitive) |
| `excludetitle` | id3 artist substrings to exclude from scan results (case-insensitive) |
| `excludealbum` | id3 album substrings to exclude from scan results (case-insensitive) |
| `excludegenre` | id3 genre substrings to exclude from scan results (case-insensitive) |
| `sortby` | (`year`, `artist`, `none`) media file sort options |
| `groupby` | (`none`, `year`) group media files into playlists using the specified tag |
| `getmp3duration` | (`true`, `false`) whether to calculate mp3 duration using `tcolgate/mp3` (**warning: slow*) |

### Dependencies
- [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) (Used for generating files.yaml)
- [dhowden/tag](https://github.com/dhowden/tag) (Used for reading ID3 tags)
    - `dhowden/tag` is a nice little Go library for ID3, MP4 and OGG/FLAC metadata parsing. I added support for [genre codes in the 148-191 range](https://en.wikipedia.org/wiki/List_of_ID3v1_genres#Extension_by_Winamp) and David was kind enough to merge my [PR](https://github.com/dhowden/tag/pull/103).
- [tcolgate/mp3](https://github.com/tcolgate/mp3) (Used for calculating MP3 duration)
    - Note: Mp3 duration calculation can be disabed by setting `getmp3duration: false` in `conf.yaml`. 
    - It's currently disabled because this library is suddenly unexpectedly slow for me.



## Performance Demo 1.5 - desktop PC with decade old i7 CPU
Much slower on first run. Why?

First run

```bash
2022/12/14 06:55:32 Successfully loaded 9124 media files
2022/12/14 06:55:32 Skipped 179 excluded media files

real	3m6.427s
user	0m4.687s
sys	0m5.567s

Second run

```bash
2022/12/14 07:09:53 Successfully loaded 9124 media files
2022/12/14 07:09:53 Skipped 179 excluded media files

real	0m3.637s
user	0m2.410s
sys	0m1.299s

```

## Performance Demo 2 - Raspberry Pi 4 model B
```bash
$ time go run mediascan.go conf.yaml files.yaml
2022/07/04 15:06:39 Successfully loaded 8376 media files

real	0m18.652s
user	0m8.251s
sys	0m11.406s
```

## Performance Demo 3 - Raspberry Pi Zero 2W
```bash
$ time go run mediascan.go conf.yaml files.yaml
2022/07/04 15:21:14 Successfully loaded 8376 media files

real	1m18.217s
user	0m18.136s
sys	0m14.305s
```

## Installing Dependencies

Install Go programming language compiler, linker, compiled stdlib and supplementary Go tools
```bash
sudo apt update
sudo apt upgrade
sudo apt install golang-go golang-src golang-doc golang-golang-x-tools
```

Edit `~/.bashrc` by adding the following lines:

```bash
export GOPATH=~/go
export GOROOT=/usr/local/go
export PATH=$PATH:$GOPATH/bin
export PATH=$PATH:$GOROOT/bin
```
Then reload it: `source ~/.bashrc`

Note: `GOROOT` directory may vary. To confirm, check the `go` symlink, e.g.:
```bash
$ which go
/usr/bin/go
$ ls -l /usr/bin/go
lrwxrwxrwx 1 root root 21 Sep 16  2020 /usr/bin/go -> ../lib/go-1.15/bin/go
```
So based on the above output on the rpi we need: `export GOROOT=/usr/lib/go-1.15`

To test it (Raspberry Pi):
```bash
$ $GOROOT/bin/go version
go version go1.15.15 linux/arm64
```

To test it (Linux desktop):
```bash
$ $GOROOT/bin/go version
go version go1.18.1 linux/amd64
```

Install Go dependency 1 of 2: `tag`
```bash
go get github.com/dhowden/tag/cmd/tag
```

Install Go dependency 2 of 2: `yaml.v3`
```bash
go get gopkg.in/yaml.v3
```

Now we can see that these packages were installed into `$GOPATH/src`:
```bash
$ ls $GOPATH/src
github.com  gopkg.in
```

Now install mediascan:
```bash
go get github.com/bretttolbert/mediascan
cd ~/go/src/github.com/bretttolbert/mediascan/ && go install
```

Edit [conf.yaml](conf.yaml) and set the `mediadir` to the desired directory path.

Now you should be all set to run mediascan.

You can run it from the mediascan directory like this:
```bash
go run mediascan.go conf.yaml files.yaml
```

Or you can run it from any directory like this:
```bash
go run github.com/bretttolbert/mediascan conf.yaml files.yaml
```

Happy scanning!
