package util

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
)

var ErrInvalidArguments = errors.New("invalid arguments")

// Args are all the possible commandline arguments
type Args struct {
	Datafile  string // Mutually exclusive with CRC
	CRC       string // Mutually exclusive with Datafile
	ImagePath string // Could be a file. Could be a dir. We'll check in verifyArgs.
	OutPath   string // It's a string rather than a File as we won't create it until after we have confirmed everything exists.
	Upscale   bool
	Verbose   bool
}

// ParseArgs parses the commandline args as well as verifies that the correct set of them is provided.
func ParseArgs() (Args, error) {
	args := Args{}
	flag.Usage = printUsage
	flag.BoolVar(&args.Upscale, "upscale", false, "Resizes images less than 170 pixels high")
	flag.BoolVar(&args.Upscale, "u", false, "Resizes images less than 170 pixels high")

	flag.BoolVar(&args.Verbose, "verbose", false, "Turns on verbose logging")
	flag.BoolVar(&args.Verbose, "v", false, "Turns on verbose logging")

	flag.StringVar(&args.CRC, "crc", "", "CRC value for single image processing")
	flag.StringVar(&args.CRC, "c", "", "CRC value for single image processing")

	flag.StringVar(&args.Datafile, "datafile", "", "no-intro dat-o-matic datafile")
	flag.StringVar(&args.Datafile, "d", "", "no-intro dat-o-matic datafile")

	flag.StringVar(&args.ImagePath, "in", "", "Image input src")
	flag.StringVar(&args.ImagePath, "i", "", "Image input src")

	// Use the current working directory as the default output dir if not specified
	// Probably should warn against this if you're using a datafile though
	out, err := os.Getwd()
	if err != nil {
		return Args{}, fmt.Errorf("os.Getwd: %w", err)
	}
	flag.StringVar(&args.OutPath, "out", out, "Output directory")
	flag.StringVar(&args.OutPath, "o", out, "Output directory")

	flag.Parse()

	args, err = verifyArgs(args)
	if errors.Is(err, ErrInvalidArguments) {
		flag.Usage()
		os.Exit(2)
	}

	return args, err
}

func verifyArgs(args Args) (Args, error) {
	if len(args.Datafile) == 0 && len(args.CRC) == 0 {
		fmt.Println("ERROR: one of --datafile or --crc must be specified")
		return args, ErrInvalidArguments
	}
	if len(args.Datafile) != 0 && len(args.CRC) != 0 {
		fmt.Println("ERROR: --datafile and --crc are mutually exclusive")
		return args, ErrInvalidArguments
	}
	if len(args.ImagePath) == 0 {
		fmt.Println("ERROR: --in must be specified")
		return args, ErrInvalidArguments
	}
	if len(args.OutPath) == 0 {
		fmt.Println("ERROR: --out must be specified")
		return args, ErrInvalidArguments
	}

	// Image source is the most complex: it can be a file or a directory
	if imgSrc, err := os.Stat(args.ImagePath); errors.Is(err, os.ErrNotExist) {
		return Args{}, fmt.Errorf("input source %s does not exist", args.ImagePath)
	} else if err != nil {
		return Args{}, fmt.Errorf("error opening input src %s: %w", args.ImagePath, err)
	} else if len(args.Datafile) != 0 && !imgSrc.IsDir() {
		fmt.Println("ERROR: --in must be a dir if --datafile is specified")
		return args, ErrInvalidArguments
	} else if len(args.CRC) != 0 && imgSrc.IsDir() {
		fmt.Println("ERROR: --in must be a file if --crc is specified")
		return args, ErrInvalidArguments
	}

	if len(args.CRC) != 0 {
		if match, err := regexp.MatchString(`[[:xdigit:]]{8}`, args.CRC); !match {
			return Args{}, fmt.Errorf("%s is not a valid crc32 hash", args.CRC)
		} else if err != nil {
			return Args{}, fmt.Errorf("regex error: %w", err)
		}
	}

	if len(args.Datafile) > 0 {
		if datafile, err := os.Stat(args.Datafile); errors.Is(err, os.ErrNotExist) {
			return Args{}, fmt.Errorf("datafile %s does not exist", args.Datafile)
		} else if err != nil {
			return Args{}, fmt.Errorf("error opening input src %s: %w", args.Datafile, err)
		} else if datafile.IsDir() {
			return Args{}, fmt.Errorf("datafile %s is a directory", args.Datafile)
		}
	}
	return args, nil
}

func printUsage() {
	fmt.Println(`
Usage of thumbnailizer:
-d, --datafile Path to no-intro dat-o-matic datafile for multi image processing mode
-c, --crc      CRC32 of game for single image processing mode
-i, --in       Path to image file (for crc mode) or image directory (for datafile mode)
-o, --out	   Output directory
-v, --verbose  Turns on verbose logging
-u, --upscale  Resizes images less than 170 pixels high
-h, --help     Prints this message`)
}
