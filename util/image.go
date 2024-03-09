package util

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image/png"
	"os"
	"strings"

	"github.com/disintegration/imaging"

	"github.com/g026r/pocket-thumbnailizer/model"
)

// maxImgSize is the maximum height of the image on the game details screen.
// Much larger than this & it starts getting cropped.
const maxImgSize = 175

// imgHeader = " IPA"
// Why this? I dunno. But it's what's necessary.
var imgHeader = []byte{0x20, 0x49, 0x50, 0x41}

// ProcessGames takes a list of games & transforms any images it finds for them into Pocket-compatible bin files
func ProcessGames(args Args, games []model.Game) (int, error) {
	// For each game in the datafile:
	// 1. Determine if it's a png or a jpg (libretro-thumbnails is all pngs, but this is for my future use)
	processed := 0
	for _, g := range games {
		img := args.ImagePath
		// If we're in multifile mode we need to build the file from the game name
		if len(args.Datafile) != 0 {
			//libretro uses '_' instead of '&' in file names
			img = fmt.Sprintf("%s/%s.png", img, strings.Replace(g.Name, "&", "_", -1))

			// Check if it's a PNG or a JPEG
			if _, err := os.Stat(img); errors.Is(err, os.ErrNotExist) {
				img = fmt.Sprintf("%s/%s.jpg", args.ImagePath, g.Name)
				if _, err := os.Stat(img); errors.Is(err, os.ErrNotExist) {
					if args.Verbose {
						fmt.Printf("%s: image file does not exist. Skipping.\n", g.Name)
					}
					continue
				}
			}
		}

		err := writeFile(g.ROM.CRC32, img, args.OutPath, args.Upscale)
		if err != nil {
			return processed, fmt.Errorf("util.WriteFile error: %w", err)
		}
		processed++
	}

	return processed, nil
}

// WriteFile does what it says: write the image file out to disk
// hash is the crc32 that will be used for the filename
// src is the full path to the image being processed
// outDir is the directory to write the file to; it will be created if it doesn't exist
// upscale will cause it to stretch images less than maxImgSize in height to maxImgSize
func writeFile(hash, src string, outDir string, upscale bool) error {
	img, err := imaging.Open(src)
	if err != nil {
		// A bunch of libretro-thumbnail images have invalid checksums & are simply un-openable via Go's strict image loader
		// Let the user know so that it can be fixed.
		if errors.Is(err, png.FormatError("invalid checksum")) {
			fmt.Printf("Image %s has an invalid checksum. Try opening & re-saving it in an image editor.", src)
		}
		return fmt.Errorf("imaging.Open: %s: %w", src, err)
	}

	// rotate 90 degrees.
	// Necessary for the Pocket for some reason, though I've been told the Duo doesn't need it.
	// I'll worry about that if ever I get a Duo, I guess.
	rotated := imaging.Rotate90(img) // return type: image.NRGBA

	if rotated.Rect.Max.X > maxImgSize { // Only resize non box art if the image is too big.
		rotated = imaging.Resize(rotated, maxImgSize, 0, imaging.Lanczos)
	} else if upscale {
		// We scale here. Should I use a different scaling algo?
		rotated = imaging.Resize(rotated, maxImgSize, 0, imaging.Lanczos)
	}
	width := rotated.Rect.Max.X
	height := rotated.Rect.Max.Y

	// Why use BGRA, Analogue? Why?
	bgra := make([]byte, len(rotated.Pix))
	for i := 0; i < len(rotated.Pix); i = i + 4 {
		bgra[i] = rotated.Pix[i+2]
		bgra[i+1] = rotated.Pix[i+1]
		bgra[i+2] = rotated.Pix[i]
		bgra[i+3] = rotated.Pix[i+3]
	}

	err = MakeDir(outDir) // This occurs here & not earlier why again? I forget.
	if err != nil {
		return fmt.Errorf("WriteFile MakeDir: %w", err)
	}
	outSrc := fmt.Sprintf("%s/%s.bin", outDir, hash)

	outFile, err := os.Create(outSrc)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}
	defer outFile.Close()

	// Now we write the header.
	_, err = outFile.Write(imgHeader)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	// Specification requires the height then the width to be written out as little endian bytes
	err = binary.Write(outFile, binary.LittleEndian, (int16)(height))
	if err != nil {
		return fmt.Errorf("height: %w", err)
	}
	err = binary.Write(outFile, binary.LittleEndian, (int16)(width))
	if err != nil {
		return fmt.Errorf("width: %w", err)
	}

	// NB: I'm just ignoring stride & hoping it works. Fucking stride...
	_, err = outFile.Write(bgra)
	if err != nil {
		return fmt.Errorf("img data: %w", err)
	}

	return nil
}
