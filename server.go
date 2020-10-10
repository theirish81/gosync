package main

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/hoisie/web"
	log "github.com/sirupsen/logrus"
	"github.com/tg123/go-htpasswd"
)

// Server the server
type Server struct {
	config *Config
}

// auth authenticates a web context against the password file
func (s *Server) auth(ctx *web.Context) bool {
	username, password, err := ctx.GetBasicAuth()
	if err != nil {
		return false
	}
	passwd, err2 := htpasswd.New(s.config.Fs.PasswordFile, htpasswd.DefaultSystems, nil)
	if err2 != nil {
		return false
	}
	return passwd.Match(username, password)
}

// handleUpload web handler for the upload operation
func (s *Server) handleUpload(ctx *web.Context, val string) {
	if strings.Index(val, "..") > -1 {
		log.Error("Invalid path")
		ctx.BadRequest()
		return
	}
	if !s.auth(ctx) {
		log.Error("No auth")
		ctx.Forbidden()
		return
	}
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Error("Error reading body", err)
		ctx.BadRequest()
		return
	}
	dir, _ := path.Split(s.config.Fs.RootDir + val)
	err2 := os.MkdirAll(dir, os.ModePerm)
	if err2 != nil {
		log.Error("Error making dirs", err2)
		ctx.BadRequest()
		return
	}
	err3 := ioutil.WriteFile(s.config.Fs.RootDir+val, data, 0644)
	if err3 != nil {
		log.Error("Error writing file", err3)
		ctx.BadRequest()
		return
	}
	ctx.WriteString("OK")
}

// handleDelete handles the deletion of a file
func (s *Server) handleDelete(ctx *web.Context, val string) {
	if !s.auth(ctx) {
		log.Error("No auth")
		ctx.Forbidden()
		return
	}
	err := os.RemoveAll(s.config.Fs.RootDir + val)
	if err != nil {
		ctx.WriteString(err.Error())
		ctx.BadRequest()
		return
	}
	ctx.WriteString("OK")
}

// StartServer Start the server
func StartServer(config *Config) {
	log.Info("GoSync starting in server mode for directory: " + config.Fs.RootDir)
	server := Server{config}
	web.Post("/files/(.*)", server.handleUpload)
	web.Delete("/files/(.*)", server.handleDelete)
	web.Run("0.0.0.0:" + strconv.Itoa(config.Server.Port))
}
