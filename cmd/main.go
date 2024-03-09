package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
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
		fmt.Println("Unable to parse arguments:", err)
		os.Exit(1)
	}

	var datafile model.Datafile
	if len(args.Datafile) != 0 {
		// Load the datafile from disk & unmarshal it
		b, err := os.ReadFile(args.Datafile)
		if err != nil {
			fmt.Println("Unable to load datafile:", err)
			os.Exit(1)
		}
		err = xml.Unmarshal(b, &datafile)
		if err != nil {
			fmt.Println("Unable to parse datafile:", err)
			os.Exit(1)
		}
	} else { // No datafile means we're in single-file mode
		name := args.ImagePath
		if ext, _ := regexp.MatchString(`\.[[:alnum:]]+$`, args.ImagePath); ext {
			name = name[:strings.LastIndex(name, ".")]
		}
		if sep := strings.LastIndex(name, string(os.PathSeparator)); sep != -1 {
			name = name[sep+1:]
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
		fmt.Println("Unable to create output dir", args.OutPath)
		os.Exit(1)
	}

	fmt.Printf("Found %d entries & beginning processing. This may take a while...\n", len(datafile.Games))

	processed, err := util.ProcessGames(args, datafile.Games)
	if err != nil {
		fmt.Println("Error processing games:", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully processed %d game entries\n", processed)
	os.Exit(0)
}
