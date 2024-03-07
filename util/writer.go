package util

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/disintegration/imaging"
)

var ImgHeader = []byte{0x20, 0x49, 0x50, 0x41}

// WriteFile does what it says: write the image file out to disk
func WriteFile(hash, src string, outDir string, boxArt bool) error {
	imgFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer imgFile.Close()

	img, err := imaging.Open(src) // Weirdly `resize` doesn't work correctly if you use `Decode` instead of `Open`
	if err != nil {
		return fmt.Errorf("imgconv.Open: %w", err)
	}

	rotated := imaging.Rotate90(img) // image.NRGBA

	if boxArt {
		// We scale here. Not sure on the algo though.
		rotated = imaging.Resize(rotated, 175, 0, imaging.Lanczos)
	} else if rotated.Rect.Max.X > 175 { // Only resize non box art if the image is too big.
		rotated = imaging.Resize(rotated, 175, 0, imaging.Lanczos)
	}
	width := rotated.Rect.Max.X
	height := rotated.Rect.Max.Y

	// Hooray! We get to convert from RGBA to BGRA format.
	bgra := make([]byte, len(rotated.Pix))
	for i := 0; i < len(rotated.Pix); i = i + 4 {
		bgra[i] = rotated.Pix[i+2]
		bgra[i+1] = rotated.Pix[i+1]
		bgra[i+2] = rotated.Pix[i]
		bgra[i+3] = rotated.Pix[i+3]
	}

	MakeDir(outDir)
	outSrc := fmt.Sprintf("%s/%s.bin", outDir, hash)

	outFile, err := os.Create(outSrc)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}
	defer outFile.Close()

	// Now we write the header.
	_, err = outFile.Write(ImgHeader)
	if err != nil {
		return fmt.Errorf("header: %w", err)
	}

	err = binary.Write(outFile, binary.LittleEndian, (int16)(height))
	if err != nil {
		return fmt.Errorf("dim 1: %w", err)
	}
	err = binary.Write(outFile, binary.LittleEndian, (int16)(width))
	if err != nil {
		return fmt.Errorf("dim 2: %w", err)
	}

	// NB: I'm just ignoring stride & hoping it works. Fucking stride...
	_, err = outFile.Write(bgra)
	if err != nil {
		return fmt.Errorf("img data: %w", err)
	}

	return nil
}
