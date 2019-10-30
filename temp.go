package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func main() {
	fileHandler := http.FileServer(http.Dir("./video"))

	http.Handle("/video/", http.StripPrefix("/video/", fileHandler))

	http.HandleFunc("/api/list", getFileListHandler)
	http.HandleFunc("/api/upload", uploadVlogHandler)
	http.HandleFunc("/sayhello", sayHello)

	http.ListenAndServe(":8088", nil)
}

func uploadVlogHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 15*1024*1024)
	if err := r.ParseMultipartForm(15 * 1024 * 1024); err != nil {
		http.Error(w, "文件超过15M", http.StatusInternalServerError)
		return
	}

	file, fileHeader, err := r.FormFile("uploadFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if strings.HasSuffix(fileHeader.Filename, ".mp4") == false {
		http.Error(w, "upload not video", http.StatusInternalServerError)
		return
	}

	md5Byte := md5.Sum([]byte(fileHeader.Filename + time.Now().String()))
	fileNewName := fmt.Sprintf("%x", md5Byte)
	dst, err := os.Create("./video/" + fileNewName + path.Ext(fileHeader.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	defer file.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func getFileListHandler(w http.ResponseWriter, r *http.Request) {
	files, _ := filepath.Glob("./video/*")
	var ret []string
	for _, v := range files {
		fmt.Println(v)
		ret = append(ret, "http://"+r.Host+"/video/"+filepath.Base(v))
	}
	fileJson, _ := json.Marshal(ret)
	w.Write(fileJson)
	return
}
