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

### Usage

For multi-image mode, the following options are used:
* `--datafile` A [no-intro datomatic](https://datomatic.no-intro.org) dat file.
* `--in` The directory where the images reside.
* `--out` The output directory for the bin files.

For single image mode, the following options are used:
* `--crc` The crc to use for the filename.
* `--in` The image to process. See above for list of accepted types.
* `--out` The directory to output the file to.

By default images greater than 175 pixels in height will be resized down, while images less than that are untouched.
You can modify this behaviour with either `--no-resize` or `--upscale`. These options are mutually exclusive.

* `--no-resize` will prevent any resizing from occurring on images. Not recommended, as it results in large files.
* `--upscale` will, in addition to resizing large images down to 175px, resize any image below 175px in height up to 175px.

And one last flag is available:
* `--verbose` prints out a few extra logging messages. If you're using the multi-file mode & not seeing the image you expect generating, you may wish to try this flag.
