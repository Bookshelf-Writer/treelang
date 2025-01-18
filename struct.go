package main

import (
	"image"
)

// // // // // // // // // // // // // // // // // //

type LangInfoNameObj struct {
	EN  string `json:"en"`
	DEF string `json:"def"`
}

type LangInfoObj[V string | *image.RGBA] struct {
	Code string           `json:"code"`
	Name *LangInfoNameObj `json:"name"`
	Flag V                `json:"flag"`
}

type LangObj[V string | *image.RGBA] struct {
	Info *LangInfoObj[V] `json:"info"`
	Data any             `json:"data"`
}
