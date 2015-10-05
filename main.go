package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var (
	managerDir     string
	sourceDir      string
	destDir        string
	managerAddress string
	clientAddress  string
	dataFile       string
	userFile       string
)

func init() {

	flag.StringVar(&managerDir, "manager", "/Users/joostvanvark/www/prudonmanager/manager", "Manager dir to serve")
	flag.StringVar(&sourceDir, "source", "/Users/joostvanvark/www/prudonmanager/source", "Source dir.")
	flag.StringVar(&destDir, "dest", "/Users/joostvanvark/www/prudonmanager/dest", "Dest dir.")
	flag.StringVar(&managerAddress, "managerPort", ":8080", "Manager port.")
	flag.StringVar(&clientAddress, "clientPort", ":8081", "Client port.")
	flag.StringVar(&dataFile, "dataFile", "/Users/joostvanvark/www/prudonmanager/data/data.json", "Data file.")
	flag.StringVar(&userFile, "userFile", "/Users/joostvanvark/www/prudonmanager/data/user.json", "User file.")
	flag.Parse()

	initData()
	initUser()
	initRender()

}

func main() {

	r := mux.NewRouter()
	InitServer(r)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(managerDir)))

	go func() {
		log.Panic(http.ListenAndServe(managerAddress, r))
	}()

	c := mux.NewRouter()
	c.PathPrefix("/").Handler(http.FileServer(http.Dir(destDir)))

	go func() {
		log.Panic(http.ListenAndServe(clientAddress, c))
	}()

	select {}

}
