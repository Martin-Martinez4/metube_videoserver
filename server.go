package main

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func main() {

	PORT := "8081"

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc:  AllowOriginFunc,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/media/{videoname:[\\w]+}/stream/", streamInit)
	r.Get("/media/{videoname:[\\w]+}/stream/{segment:index[0-9]+.ts}", streamContinue)

	http.ListenAndServe(":"+PORT, r)

}

// Change to be more secure later
func AllowOriginFunc(r *http.Request, origin string) bool {
	// if origin == "http://example.com" {
	// 	return true
	// }
	// return false

	return true
}

func streamInit(w http.ResponseWriter, req *http.Request) {

	videoname := chi.URLParam(req, "videoname")

	mediaBase := getMediaBase(videoname)
	m3u8Name := "index.m3u8"
	mediaFile := filepath.Join(mediaBase, m3u8Name)

	w.Header().Set("Content-Type", "application/x-mpegURL")
	http.ServeFile(w, req, mediaFile)
}

func streamContinue(w http.ResponseWriter, req *http.Request) {

	segment := chi.URLParam(req, "segment")
	videoname := chi.URLParam(req, "videoname")

	mediaBase := getMediaBase(videoname)
	mediaFile := filepath.Join(mediaBase, segment)

	w.Header().Set("Content-Type", "video/MP2T")
	http.ServeFile(w, req, mediaFile)

}

func getMediaBase(videoname string) string {
	mediaRoot := "videos"
	return filepath.Join(mediaRoot, videoname)
}
