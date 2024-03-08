package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/g026r/pocket-thumbnailizer/model"
	"github.com/g026r/pocket-thumbnailizer/util"
)

func main() {
	var args util.Args
	var err error
	if args, err = util.ParseArgs(); errors.Is(err, util.ErrInvalidArguments) {
		flag.Usage()
		os.Exit(2)
	} else if err != nil {
		log.Fatalf("parse args error: %v", err)
	}

	var datafile model.Datafile
	if len(args.Datafile) != 0 {
		// Load the datafile from disk & unmarshal it
		b, err := os.ReadFile(args.Datafile)
		if err != nil {
			log.Fatalf("read datafile: %v", err)
		}
		err = xml.Unmarshal(b, &datafile)
		if err != nil {
			log.Fatalf("unmarshal: %v", err)
		}
	} else { // No datafile means we're in single-file mode
		name := args.ImagePath
		if ext, _ := regexp.MatchString(`\.[[:alnum:]]$`, args.ImagePath); ext {
			name = name[:strings.LastIndex(name, ".")]
		}
		if sep := strings.LastIndex(name, strconv.QuoteRune(os.PathSeparator)); sep != -1 {
			name = name[sep:]
			args.ImagePath = args.ImagePath[:sep] // Remove the image from this so we can use a generic processing block
		}
		datafile = model.Datafile{
			Games: []model.Game{{
				Name: name,
				ROM: model.Rom{
					CRC32: args.CRC,
				},
			}},
		}
	}

	// Create the output dir if it doesn't exist
	err = util.MakeDir(args.OutPath)
	if err != nil {
		log.Fatalf("Unable to create output dir %s", args.OutPath)
	}

	fmt.Println(fmt.Sprintf("Found %d entries & beginning processing. This may take a while...", len(datafile.Games)))

	// For each game in the datafile:
	// 1. Determine if it's a png or a jpg (libretro-thumbnails is all pngs, but this is for my future use)
	processed := 0
	for _, g := range datafile.Games {
		// libretro uses `_` instead of `&` in file names
		img := fmt.Sprintf("%s/%s.png", args.ImagePath, strings.Replace(g.Name, "&", "_", -1))
		if _, err := os.Stat(img); errors.Is(err, os.ErrNotExist) {
			img = fmt.Sprintf("%s/%s.jpg", args.ImagePath, g.Name)
			if _, err := os.Stat(img); errors.Is(err, os.ErrNotExist) {
				if args.Verbose {
					fmt.Println(fmt.Sprintf("File for %s does not exist. Skipping.", g.Name))
				}
				continue
			}
		}
		err = util.WriteFile(g.ROM.CRC32, img, args.OutPath, args.Upscale)
		if err != nil {
			log.Fatalf("util.WriteFile: %v", err)
		}
		processed++
	}

	fmt.Println(fmt.Sprintf("Successfully processed %d game entries", processed))
}
