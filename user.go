package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type user struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

func initUser() {

	if _, err := os.Stat(userFile); os.IsNotExist(err) {

		u := user{"hans", "hans"}
		writeJson(u, userFile)

	}

}

func SetUser(username, password string) {

	writeJson(user{username, password}, userFile)

}

func GetUser(username, password string) bool {

	d, err := ioutil.ReadFile(userFile)
	if err != nil {
		log.Println(err)
		return false
	}

	var result user
	err = json.Unmarshal(d, &result)
	if err != nil {
		log.Println(err)
		return false
	}

	return result.Name == username && result.Pass == password

}
