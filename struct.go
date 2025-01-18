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

type LangSysObj struct {
	Date string `json:"date"`
	Hash string `json:"hash"`
}
type LangObj[V string | *image.RGBA] struct {
	Info *LangInfoObj[V] `json:"info"`
	Sys  *LangSysObj     `json:"sys"`
	Data any             `json:"data"`
}
