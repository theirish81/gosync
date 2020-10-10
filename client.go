package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/andreaskoch/go-fswatch"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// skipDotFilesAndFolders determines whether the path should be skipped. No skips by default
func skipDotFilesAndFolders(path string) bool {
	return false
}

// StartClient starts the client
func StartClient(config *Config) {
	log.Info("GoSync starting in client mode for directory: " + config.Fs.RootDir)
	data, _ := ioutil.ReadFile(config.Fs.PasswordFile)
	auth := make(map[string]string)
	err := yaml.Unmarshal(data, auth)
	if err != nil {
		log.Fatal("Unable to read passwd file: ", err.Error())
		os.Exit(1)
	}
	fileWatcher := fswatch.NewFolderWatcher(config.Fs.RootDir, true, skipDotFilesAndFolders, 1)
	fileWatcher.Start()
	for fileWatcher.IsRunning() {

		select {
		case <-fileWatcher.Modified():
		case <-fileWatcher.Moved():

		case changes := <-fileWatcher.ChangeDetails():

			// New files
			for _, n := range changes.New() {
				performPost(n, config, auth)
			}

			// Modified files
			for _, n := range changes.Modified() {
				performPost(n, config, auth)
			}

			// Deleted files
			for _, n := range changes.Moved() {
				performDelete(n, config, auth)
			}
		}
	}
}

// performPost executes a post call to update/create a file
func performPost(n string, config *Config, auth map[string]string) {
	relPath := extractRelative(n, config)
	data, err := ioutil.ReadFile(config.Fs.RootDir + relPath)
	if err != nil {
		log.Error("Could not read file: ", err.Error())
		return
	}
	address := config.Server.Address + ":" + strconv.Itoa(config.Server.Port) + "/files/" + relPath
	client := &http.Client{}
	req, err2 := http.NewRequest("POST", address, bytes.NewReader(data))
	if err2 != nil {
		log.Error("Could not build request: ", err.Error())
		return
	}
	req.SetBasicAuth(auth["username"], auth["password"])
	log.Debug("POST " + address)
	resp, err3 := client.Do(req)
	if resp.StatusCode != 200 {
		log.Error("Request failed")
		return
	}
	if err3 != nil {
		log.Error("Request failed: ", err3.Error())
	}
}

// performDelete performs the call to delete a file
func performDelete(n string, config *Config, auth map[string]string) {
	relPath := extractRelative(n, config)
	client := &http.Client{}
	address := config.Server.Address + ":" + strconv.Itoa(config.Server.Port) + "/files/" + relPath
	req, err := http.NewRequest("DELETE", address, nil)
	if err != nil {
		log.Error("Could not build request: ", err.Error())
		return
	}
	req.SetBasicAuth(auth["username"], auth["password"])
	log.Debug("DELETE " + address)
	resp, err2 := client.Do(req)
	if resp.StatusCode != 200 {
		log.Error("Request failed")
	}
	if err2 != nil {
		log.Error("Request failed: ", err2.Error())
	}
}

// Extracts the relative file of the file being modified
func extractRelative(path string, config *Config) string {
	return path[len(config.Fs.RootDir):]
}
