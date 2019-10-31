package controllers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// MainController create Main Controller
type MainController struct {
	Router *mux.Router
}

// Initialize Initialize controller
func (controller *MainController) Initialize() {
	controller.Router = mux.NewRouter()
	controller.initRoutes()
}

// initRoutes Initialize routes
func (controller *MainController) initRoutes() {
	controller.Router.HandleFunc("/", http.HandlerFunc(controller.showIndex))
	controller.Router.Handle("/video/{video_name}", http.HandlerFunc(controller.getVideo)).Methods("GET")
	controller.Router.HandleFunc("/api/upload", http.HandlerFunc(controller.uploadVideo))
	controller.Router.HandleFunc("/api/list", http.HandlerFunc(controller.getVideoList))
}

func (controller *MainController) getVideo(w http.ResponseWriter, r *http.Request) {
	fileHandler := http.FileServer(http.Dir("./videos/"))
	http.StripPrefix("/video/", fileHandler).ServeHTTP(w, r)
}

// uploadVideo upload a video to server
func (controller *MainController) uploadVideo(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
	if err := r.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		http.Error(w, "file size too large, limit 10M", http.StatusInternalServerError)
		return
	}
	f, fHeader, err := r.FormFile("upload_video")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fName := fHeader.Filename
	if strings.HasSuffix(fName, ".mp4") == false {
		http.Error(w, "file not mp4 file, please check it", http.StatusForbidden)
		return
	}
	defer f.Close()

	fileNewName := getNewName(fName + time.Now().String())
	newFile, err := os.Create("./videos/" + fileNewName + path.Ext(fName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	if _, err := io.Copy(newFile, f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

// getVideoList get all video
func (controller *MainController) getVideoList(w http.ResponseWriter, r *http.Request) {
	files, err := filepath.Glob("./videos/*")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var ret []string
	for _, v := range files {
		ret = append(ret, "http://"+r.Host+"/video/"+filepath.Base(v))
	}

	retJson, err := json.Marshal(ret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(retJson)
	return
}

func (controller *MainController) showIndex(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("./view/index.html")
	if err != nil {
		logs.Error("open Index failed")
		return
	}
	defer file.Close()
	buffer := make([]byte, 15*1024*1024)
	file.Read(buffer)
	w.Write(buffer)
}

// getNewName Util function
func getNewName(name string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(name)))
}
