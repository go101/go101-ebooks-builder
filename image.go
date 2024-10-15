package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"

	//"golang.org/x/image/math/fixed"
	"github.com/golang/freetype"
	"golang.org/x/image/font/gofont/goregular"
	//"github.com/golang/freetype/truetype"
)

const ExternalLinkPNG = `iVBORw0KGgoAAAANSUhEUgAAABgAAAAYCAYAAADgdz34AAAABGdBTUEAALGPC/xhBQAAAAFzUkdCAK7OHOkAAAAgY0hSTQAAeiYAAICEAAD6AAAAgOgAAHUwAADqYAAAOpgAABdwnLpRPAAAAAZiS0dEAAAAAAAA+UO7fwAAAAlwSFlzAAAASAAAAEgARslrPgAAAKpJREFUSMftlTsOg0AMRF8QNd1KuVCkhPvmIlDQk+3ScYFQsJEC7MdmoYkYyZ1nnm2tAE4pVQMW+ETqK0nPSq+EMRuQbPDIAM1RAAO0ztNK/BrA7+QdcN0T4AsX+SWAUPgugFi4SDFAdngMIA3fdCLN5GpAxfydG+0FyoRhAJ6u7wa8yZRvxYvbZKtf16AFFLkrp7QE2MUk2oLpkx/UA9k/IVQ9cD/6Kn+mESDFiPdj8h9+AAAAJXRFWHRkYXRlOmNyZWF0ZQAyMDE5LTA0LTA2VDIxOjAzOjUzKzAwOjAwwCzL2wAAACV0RVh0ZGF0ZTptb2RpZnkAMjAxOS0wNC0wNlQyMTowMzo1MyswMDowMLFxc2cAAAAodEVYdHN2ZzpiYXNlLXVyaQBmaWxlOi8vL3RtcC9tYWdpY2stbkRub0JiNjHXYHRtAAAAAElFTkSuQmCC`

const CoverImageFilename = "101-front-cover-1400x.png"
const CoverImageTempFilePattern = "101-front-cover-*.png"

func buildCoverImage(bookProjectDir, bookVersion string) string {

	revison := ""
	hash := runShellCommand2(true, bookProjectDir, "git", "rev-parse", bookVersion)
	hash = bytes.TrimSpace(hash)
	if len(hash) > 7 {
		hash = hash[:7]
		revison = string(hash)
	}
	
	// git log -1 --pretty='%ad' --date=format:'%Y/%m/%d' v1.16.a
	// 2021/02/18

	var versionText string
	if revison != "" {
		versionText = "-= " + bookVersion + "-" + revison + " =-"
	} else {
		versionText = "-= " + bookVersion + " =-"
	}

	// Load cover image
	inFile, err := os.Open(filepath.Join(bookProjectDir, "pages", ArticlesFolder, "res", CoverImageFilename))
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	pngImage, err := png.Decode(inFile)
	if err != nil {
		log.Fatal(err)
	}

	// Draw cover image
	output := image.NewRGBA(image.Rect(0, 0, pngImage.Bounds().Max.X, pngImage.Bounds().Max.Y))
	draw.Draw(output, output.Bounds(), pngImage, image.ZP, draw.Src)

	// Load font
	utf8Font, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	// Draw text
	dpi := float64(72)
	fontsize := float64(29.0)
	//spacing := float64(1.5)

	ctx := new(freetype.Context)
	ctx = freetype.NewContext()
	ctx.SetDPI(dpi)
	ctx.SetFont(utf8Font)
	ctx.SetFontSize(fontsize)
	ctx.SetClip(output.Bounds())
	ctx.SetDst(output)

	pt := freetype.Pt(0, int(ctx.PointToFixed(fontsize)>>6))
	ctx.SetSrc(image.NewUniform(color.RGBA{0, 0, 0, 0}))
	extent, err := ctx.DrawString(versionText, pt)
	if err != nil {
		log.Fatal(err)
	}

	pt = freetype.Pt(output.Bounds().Max.X/2, 469)
	pt.X -= extent.X / 2

	ctx.SetSrc(image.NewUniform(color.RGBA{0, 0, 0, 255}))
	_, err = ctx.DrawString(versionText, pt)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(extent)
	//pt.Y += ctx.PointToFixed(fontsize * spacing)

	// Save new cover image
	pngFilename := mustCreateTempFile(CoverImageTempFilePattern, nil)

	pngFile, err := os.Create(pngFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer pngFile.Close()

	err = png.Encode(pngFile, output)
	if err != nil {
		log.Fatal(err)
	}

	err = pngFile.Sync()
	if err != nil {
		log.Fatal(err)
	}

	return pngFilename
}
