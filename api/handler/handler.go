package handler

import (
	"backend1/app/services/upload"
	"encoding/json"
	"fmt"
	"log"

	"net/http"
)

type Router struct {
	*http.ServeMux
	us *upload.UploadService

	HostAddr string
}

func NewRouter(us *upload.UploadService) *Router {
	r := &Router{
		ServeMux: http.NewServeMux(),
		us:       us,
	}
	r.HandleFunc("/upload", http.HandlerFunc(r.uploadFile))
	r.HandleFunc("/all_files", http.HandlerFunc(r.getAllFileInfo))
	r.HandleFunc("/all_files_ext", http.HandlerFunc(r.getFilesWithExt))
	return r
}

func (rt *Router) getFilesWithExt(w http.ResponseWriter, r *http.Request) {
	ext := r.URL.Query().Get("ext")
	if ext == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	files, err := rt.us.AllFilesFilterExt(ext)
	if err != nil {
		http.Error(w, "error when reading", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(files)
}

func (rt *Router) getAllFileInfo(w http.ResponseWriter, r *http.Request) {
	files, err := rt.us.AllFiles()
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(files)
}

func (rt *Router) uploadFile(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filePath := rt.us.UploadDir + "/" + header.Filename

	err = rt.us.WriteFile(filePath, file)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileLink := rt.HostAddr + "/" + header.Filename

	fmt.Fprintln(w, fileLink)
}
