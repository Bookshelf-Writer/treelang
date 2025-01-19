package main

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
	Info *LangInfoObj   `json:"info"`
	Sys  *LangSysObj    `json:"sys"`
	Data map[string]any `json:"data"`
}
