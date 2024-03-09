## pocket-thumbnailizer: a utility to generate Analogue Pocket library thumbnails

Simple commandline utility hacked together so that I could generate Analogue Pocket library thumbnails without having to launch a Windows virtual environment every time.

It can work in multi-image or single image processing mode, allowing generation of an entire system's library or a single game. (i.e. if a game was missing from the library or it had the wrong art.)

Supports the following image formats:
* JPG
* GIF
* PNG
* WEBP
* BMP
* TIFF

It currently assumes you know how to build & run a go application. I'll try to get some binaries up eventually.

### Usage

For multi-image mode, the following options are used:
* `--datafile` A [no-intro datomatic](https://datomatic.no-intro.org) dat file.
* `--in` The directory where the images reside.
* `--out` The output directory for the bin files.

For single image mode, the following options are used:
* `--crc` The crc to use for the filename.
* `--in` The image to process. Can handle jpg, png, and gif are supported.
* `--out` The directory to output the file to.

You can also pass it `--upscale` or `--verbose` in both modes.

* `--upscale` will resize any image below 175 pixels in height up to 175 pixels. (Images will always be resized down to 175 pixels in height if they're larger to avoid cropping.)
* `--verbose` prints out a few extra logging messages. If you're using the multi-file mode & not seeing the image you expect generating, you may wish to try this flag.
