/* This file is generated automatically and any changes made to it may be overwritten. Please do not modify the contents of this file manually. */
import fileLangENG from './test.en.json'
import fileLangRUS from './test.ru.json'

interface LangInfoObj {date: string; hash: string; img: string; name: { short: string; english: string; native: string; }}
interface LangDataObj { arrayString1: string[]; groupString1: {a2: string; a3: string;}, string_value_1: string; }
interface LangObj {info: LangInfoObj; data: LangDataObj;}
interface LanguagesTreeObj {eng: LangObj; rus: LangObj; }

export const Languages: LanguagesTreeObj = {
	eng: {
		info: {
			img: fileLangENG.info.img,
			date: fileLangENG.info.date,
			hash: fileLangENG.info.hash,
			name: {
				short: "ENG",
				english: "English",
				native: "English",
			},
		},
		data: fileLangENG.data,
	},
	rus: {
		info: {
			img: fileLangRUS.info.img,
			date: fileLangRUS.info.date,
			hash: fileLangRUS.info.hash,
			name: {
				short: "RUS",
				english: "Russian",
				native: "Русский",
			},
		},
		data: fileLangRUS.data,
	},

}
