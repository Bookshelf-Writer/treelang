package main

import (
	"bytes"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"image"
)

// // // // // // // // // // // // // // // // // //

type LangInfoNameObj struct {
	EN  string `json:"en"`
	DEF string `json:"def"`
}

type LangInfoObj struct {
	Code string           `json:"code"`
	Name *LangInfoNameObj `json:"name"`
	Flag string           `json:"flag"`
}

type LangSysObj struct {
	Date string `json:"date"`
	Hash string `json:"hash"`
}

type LangObj struct {
	Info *LangInfoObj `json:"info"`
	Sys  *LangSysObj  `json:"sys"`
	Data any          `json:"data"`
}

// // // // //

func (obj *LangInfoObj) PNG() (*image.RGBA, error) {
	svgReader := bytes.NewReader([]byte(obj.Flag))
	icon, err := oksvg.ReadIconStream(svgReader)
	if err != nil {
		return nil, err
	}

	width := 600
	height := 400
	if icon.ViewBox.W != 0 || icon.ViewBox.H != 0 {
		width = int(icon.ViewBox.W)
		height = int(icon.ViewBox.H)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	icon.SetTarget(0, 0, float64(width), float64(height))
	dasher := rasterx.NewDasher(width, height, rasterx.NewScannerGV(width, height, img, img.Bounds()))
	icon.Draw(dasher, 1.0)

	return img, nil
}
