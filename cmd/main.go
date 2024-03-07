package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/g026r/pocket-thumbnailizer/model"
	"github.com/g026r/pocket-thumbnailizer/util"
)

const isBoxArt = true

func main() {

	if len(os.Args) != 4 {
		printUsage()
	}
	dataFile := os.Args[1]
	inDir := os.Args[2]
	outDir := os.Args[3]

	var datafile model.Datafile
	b, err := ioutil.ReadFile(dataFile)
	if err != nil {
		log.Fatalf("read datafile: %v", err)
	}
	err = xml.Unmarshal(b, &datafile)
	if err != nil {
		log.Fatalf("unmarshal: %v", err)
	}

	util.MakeDir(outDir)

	processed := 0
	for _, g := range datafile.Games {
		img := fmt.Sprintf("%s/%s.png", inDir, strings.Replace(g.Name, "&", "_", -1))
		if _, err := os.Stat(img); errors.Is(err, os.ErrNotExist) {
			img = fmt.Sprintf("%s/%s.jpg", inDir, g.Name)
			if _, err := os.Stat(img); errors.Is(err, os.ErrNotExist) {
				slog.Warn("File does not exist. Skipping.", "game", g.Name)
				continue
			}
		}
		processed++
		err = util.WriteFile(g.ROM.CRC32, img, outDir, isBoxArt)
		if err != nil {
			log.Fatalf("util.WriteFile: %v", err)
		}
	}

	log.Printf("Processed %d game entries", processed)
}

func printUsage() {
	parts := strings.Split(os.Args[0], "/")
	fmt.Printf("Usage: %s <datafile> <img dir> <output dir>\n", parts[len(parts)-1])
	os.Exit(2)
}
