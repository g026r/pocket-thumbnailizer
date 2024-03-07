## pocket-thumbnailizer: a utility to generate Analogue Pocket library thumbnails

Simple golang utility hacked together so that I could generate Analogue Pocket library thumbnails without having to launch a Windows virtual environment every time.

Takes a [no-intro datomatic](https://datomatic.no-intro.org) dat file & a directory containing a set of [libretro-thumbnail](https://github.com/libretro-thumbnails) images, and from that generates a set of Analogue Pocket compliant bin files for use as library images.

Could be a bit more user friendly, but this was mostly hacked together for my own use.

It will resize box art up or down to 170px in height to best make use of the game details screen without cropping. Non box art will be scaled down to 170px if it's above that, but will be otherwise left untouched.

Currently it's hardcoded for box art as that's what I use. Change the value of `isBoxArt` in cmd/main.go to `false` if you want to use it with something else.

### Usage
`go run cmd/main.go /path/to/datomatic.dat /path/to/libretro-thumbails/files /path/to/output/dir`

The input directory should be the actual directory with the images contained in it, not the root of the libretro-thumbails copy. (i.e. point to one of Named_Boxarts, Named_Snaps, Named_Titles)

It's currently a bit chatty & will print out every entry in the dat file that it can't find an image for. It will not, however, print out when it finds an image but no datafile entry.
