package util

import "github.com/hajimehoshi/ebiten/v2"

func ScaleEbitenImage(image *ebiten.Image, dims Dims) *ebiten.Image {
	ops := &ebiten.DrawImageOptions{}
	bounds := image.Bounds().Size()
	xRatio := float64(bounds.X) / float64(dims.X)
	yRatio := float64(bounds.Y) / float64(dims.Y)
	ops.GeoM.Scale(1/xRatio, 1/yRatio)

	newImage := ebiten.NewImage(dims.X, dims.Y)
	newImage.DrawImage(image, ops)
	return newImage
}
