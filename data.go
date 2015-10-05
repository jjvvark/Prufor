package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type db struct {
	One   col `json:"one"`
	Two   col `json:"two"`
	Three col `json:"three"`
}

type col struct {
	Title string   `json:"title"`
	Paras []string `json:"paragraphs"`
}

func initData() {

	if _, err := os.Stat(dataFile); os.IsNotExist(err) {

		v := db{
			col{
				"TitleOne",
				[]string{
					"ParagraphOne",
					"ParagraphOne",
				},
			},
			col{
				"TitleTwo",
				[]string{
					"ParagraphTwo",
					"ParagraphTwo",
				},
			},
			col{
				"TitleThree",
				[]string{
					"ParagraphThree",
					"ParagraphThree",
				},
			},
		}

		writeJson(v, dataFile)

	}

}

func writeJson(v interface{}, f string) {
	dv, err := json.Marshal(v)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(f, dv, 0777)
	if err != nil {
		log.Panic(err)
	}
}

func GetDataJson() []byte {

	d, err := ioutil.ReadFile(dataFile)
	if err != nil {
		log.Panic(err)
	}

	return d

}

func GetData() db {

	d := GetDataJson()

	var result db
	err := json.Unmarshal(d, &result)
	if err != nil {
		log.Panic(err)
	}

	return result

}

func SetData(col int64, title bool, value string) error {

	d := GetData()

	if col < 1 || col > 3 {
		return errors.New("data :: SetData :: No valid col value")
	}

	if col == 1 {

		if title {

			d.One.Title = value

		} else {

			d.One.Paras = SplitParagraph(value)

		}

	} else if col == 2 {

		if title {

			d.Two.Title = value

		} else {

			d.Two.Paras = SplitParagraph(value)

		}

	} else {

		if title {

			d.Three.Title = value

		} else {

			d.Three.Paras = SplitParagraph(value)

		}

	}

	writeJson(d, dataFile)

	return nil

}

func SplitParagraph(value string) []string {

	return strings.Split(strings.TrimSpace(value), "\n")

}
