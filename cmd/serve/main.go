package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	defaultStorage    = "./storage"
	defaultListenAddr = "0.0.0.0:8080"
)

var (
	storage    = defaultStorage
	listenAddr = defaultListenAddr
	name       string
)

func main() {

	argCount := len(os.Args)
	log.Infof("Have %d args", argCount)

	name = os.Getenv("HOSTNAME")
	log.Infof("Setting worker name to: %s", name)

	if argCount >= 2 {
		storage = os.Args[1]
		log.Infof("Have argument for storage. Setting to: %s", storage)
	}

	if argCount >= 3 {
		listenAddr = os.Args[2]
		log.Infof("Have argument for listen addr. Setting to: %s", listenAddr)
	}

	r := mux.NewRouter()
	r.HandleFunc("/write", writeHandler).Methods("POST")

	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         listenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("Listening on: %s", listenAddr)
	log.Fatal(srv.ListenAndServe())
	log.Infof("Done")
}

func writeHandler(response http.ResponseWriter, request *http.Request) {

	log.Info("Starting handling write request.")
	writeFile()
	response.WriteHeader(http.StatusAccepted)
	log.Info("Done handling write request.")
}

func writeFile() {

	checkStorage()

	path := fmt.Sprintf("%s/%d.txt", storage, time.Now().UnixNano())
	d1 := []byte(fmt.Sprintln(name))
	if err := ioutil.WriteFile(path, d1, 0644); err != nil {
		log.Errorf("Error trying to write file: '%s': %s", path, err.Error())
		return
	}

	log.Infof("Successfully wrote file: '%s'", path)
}

func checkStorage() error {

	fd, err := os.Open(storage)
	defer fd.Close()

	if err != nil {
		if os.IsNotExist(err) {
			log.Infof("Directory does not exist. Creating: [%s]", storage)
			if err = os.MkdirAll(storage, 0700); err != nil {
				log.Errorf("Directory [%s] could not be created: %+v", storage, err)
				return err
			}
		} else {
			log.Errorf("Directory [%s] could not be read: %+v", storage, err)
			return err
		}
	}

	return nil
}
