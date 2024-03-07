package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/g026r/pocket-thumbnailizer/model"
	"github.com/g026r/pocket-thumbnailizer/util"
)

func main() {
	// slog.SetLogLoggerLevel(slog.LevelDebug) // Turn this on if images are failing to load

	if len(os.Args) != 4 {
		printUsage()
	}
	dataFile := os.Args[1]
	inDir := os.Args[2]
	outDir := os.Args[3]

	// Make certain the input dir exists
	if d, err := os.Stat(inDir); err != nil || !d.IsDir() {
		log.Fatalf("Error parsing input dir %s", inDir)
	}

	// Not pleased with this way to programatically do this. Should come up with something better
	isBoxArt := false
	if strings.HasSuffix(strings.ToLower(inDir), "named_boxarts") ||
		strings.HasSuffix(strings.ToLower(inDir), "named_boxarts/") {
		slog.Debug("box art determination", "isBoxArt", isBoxArt)
		isBoxArt = true
	}

	// Load the datafile from disk & unmarshal it
	var datafile model.Datafile
	b, err := os.ReadFile(dataFile)
	if err != nil {
		log.Fatalf("read datafile: %v", err)
	}
	err = xml.Unmarshal(b, &datafile)
	if err != nil {
		log.Fatalf("unmarshal: %v", err)
	}

	// Create the output dir if it doesn't exist
	err = util.MakeDir(outDir)
	if err != nil {
		log.Fatalf("Unable to create output dir %s", outDir)
	}

	slog.Info("Processing entries. This may take a while.")

	// For each game in the datafile:
	// 1. Determine if it's a png or a jpg (libretro-thumbnails is all pngs, but this is for my future use)
	// 2.
	processed := 0
	for _, g := range datafile.Games {
		// libretro uses `_` instead of `&` in file names
		img := fmt.Sprintf("%s/%s.png", inDir, strings.Replace(g.Name, "&", "_", -1))
		if _, err := os.Stat(img); errors.Is(err, os.ErrNotExist) {
			img = fmt.Sprintf("%s/%s.jpg", inDir, g.Name)
			if _, err := os.Stat(img); errors.Is(err, os.ErrNotExist) {
				slog.Warn("File does not exist. Skipping.", "game", g.Name)
				continue
			}
		}
		err = util.WriteFile(g.ROM.CRC32, img, outDir, isBoxArt)
		if err != nil {
			log.Fatalf("util.WriteFile: %v", err)
		}
		processed++
	}

	slog.Info(fmt.Sprintf("Processed %d game entries out of %d total", processed, len(datafile.Games)))
}

func printUsage() {
	parts := strings.Split(os.Args[0], "/")
	fmt.Printf("Usage: %s <datafile> <img dir> <output dir>\n", parts[len(parts)-1])
	os.Exit(2)
}
