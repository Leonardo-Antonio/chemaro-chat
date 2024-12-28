package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// obtener imagen que llega en formulario "file"
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	name := handler.Filename
	contentType := handler.Header.Get("Content-Type")
	size := handler.Size
	ext := filepath.Ext(name)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	pathFiles := path.Join("assets", "temp")
	fullpath := path.Join(pathFiles, filename)
	tempFile, err := os.Create(fullpath)
	if err != nil {
		log.Println(err)
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "ok",
		"data": map[string]any{
			"url":  fullpath,
			"name": name,
			"type": contentType,
			"size": size,
		},
	})
}
