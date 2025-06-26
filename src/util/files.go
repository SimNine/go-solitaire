package util

import (
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func LoadEbitenImageFromFile(path string) (*ebiten.Image, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	imageDecoded, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(imageDecoded), nil
}
