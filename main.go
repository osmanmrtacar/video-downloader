package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const videoDir = "videos"

var DESCRIPTION_URLS = []string{"twitter.com", "x.com", "instagram.com"}

func main() {
	if _, err := os.Stat(videoDir); os.IsNotExist(err) {
		os.Mkdir(videoDir, 0755)
	}

	var proxyManager *ProxyManager
	auth := os.Getenv("PROXY_AUTH")
	if auth != "" {
		proxyManager = NewProxyManager(auth)
	}
	videoManager := NewVideoManager(proxyManager, videoDir)

	http.HandleFunc("/videos", func(w http.ResponseWriter, r *http.Request) {
		handleVideos(w, r, videoManager)
	})
	http.HandleFunc("/videos/", handleVideoByName)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleVideos(w http.ResponseWriter, r *http.Request, vm *VideoManager) {
	switch r.Method {
	case "POST":
		var req struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid request body"))
			return
		}
		filename, description, err := vm.Download(req.URL, DESCRIPTION_URLS)
		if err != nil {
			log.Printf("POST /videos download error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		resp := map[string]string{"filename": filename, "description": description}
		json.NewEncoder(w).Encode(resp)
	case "GET":
		files, err := os.ReadDir(videoDir)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var videos []map[string]interface{}
		for _, f := range files {
			if !f.IsDir() {
				info, err := f.Info()
				var created string
				if err == nil {
					created = info.ModTime().Format("2006-01-02 15:04:05")
				}
				videos = append(videos, map[string]interface{}{
					"filename": f.Name(),
					"created":  created,
				})
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(videos)
	case "DELETE":
		// Expect timestamp in query param: /videos?before=UNIX_TIMESTAMP
		before := r.URL.Query().Get("before")
		if before == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "missing 'before' query param"})
			return
		}
		unixTs, err := strconv.ParseInt(before, 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid timestamp format"})
			return
		}
		t := time.Unix(unixTs, 0)
		files, err := os.ReadDir(videoDir)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var deleted []string
		for _, f := range files {
			if !f.IsDir() {
				info, err := f.Info()
				if err == nil && info.ModTime().Before(t) {
					err := os.Remove(filepath.Join(videoDir, f.Name()))
					if err == nil {
						deleted = append(deleted, f.Name())
					}
				}
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"deleted": deleted})
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "not found"})

	}
}

func handleVideoByName(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/videos/")
	if filename == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filePath := filepath.Join(videoDir, filename)
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, filePath)
	case "DELETE":
		if err := os.Remove(filePath); err != nil {
			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(http.StatusNotFound)

			json.NewEncoder(w).Encode(map[string]string{"error": "not found"})

			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
