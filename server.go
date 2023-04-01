package main

import (
	"compress/flate"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func main() {

	PORT := "8081"

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	compressor := middleware.NewCompressor(flate.DefaultCompression)
	r.Use(compressor.Handler)

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc:  AllowOriginFunc,
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/media/{videoname:[\\w-]+}/stream/", streamInit)
	r.Get("/media/{videoname:[\\w-]+}/stream/{secondVideoname:[\\w-]+.m3u8}", GetPlaylist)

	r.Get("/media/{videoname:[\\w-]+}/stream/{segment:[\\w-]+.ts}", streamContinue)
	r.Get("/media/{videoname:[\\w-]+}/stream/{resolution:[\\d]+}/{segment:[\\w-]+.ts}", streamWithResolution)

	r.Get("/thumbnail/{thumbnail:[\\w]+.jpg}/", getThumbnail)
	r.Get("/profile/{profile:[\\w]+}/", getProfileImage)
	r.Get("/banner/{profile:[\\w]+}/", getBanner)

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
	m3u8Name := videoname + ".m3u8"
	mediaFile := filepath.Join(mediaBase, m3u8Name)

	println(mediaFile)

	w.Header().Set("Content-Type", "application/x-mpegURL")
	http.ServeFile(w, req, mediaFile)
}

func GetPlaylist(w http.ResponseWriter, req *http.Request) {

	videoname := chi.URLParam(req, "videoname")
	secondVideoname := chi.URLParam(req, "secondVideoname")

	mediaBase := getMediaBase(videoname)
	m3u8Name := secondVideoname
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
func streamWithResolution(w http.ResponseWriter, req *http.Request) {

	segment := chi.URLParam(req, "segment")
	videoname := chi.URLParam(req, "videoname")
	resolution := chi.URLParam(req, "resolution")

	mediaBase := getMediaBase(videoname)
	mediaFile := filepath.Join(mediaBase, resolution, segment)

	w.Header().Set("Content-Type", "video/MP2T")
	http.ServeFile(w, req, mediaFile)

}

// Assuming thumbnails are in jpeg format and in the same folder as its corresponding hls files -> videos/videotitle
func getThumbnail(w http.ResponseWriter, req *http.Request) {

	thumbnailname := chi.URLParam(req, "thumbnail")

	mediabase := getMediaBase(strings.Split(thumbnailname, ".jpg")[0])
	thumbnailFile := filepath.Join(mediabase, thumbnailname)

	w.Header().Set("Content-Type", "image/JPEG")
	http.ServeFile(w, req, thumbnailFile)
}

// Assuming profiles are in jpeg format and in a folder called profiles
func getProfileImage(w http.ResponseWriter, req *http.Request) {

	profile := chi.URLParam(req, "profile")

	profileFile := filepath.Join("profiles", profile+".jpg")

	w.Header().Set("Content-Type", "image/JPEG")
	http.ServeFile(w, req, profileFile)
}

func getBanner(w http.ResponseWriter, req *http.Request) {

	profile := chi.URLParam(req, "profile")

	profileFile := filepath.Join("banner", profile+".jpg")

	w.Header().Set("Content-Type", "image/JPEG")
	http.ServeFile(w, req, profileFile)
}

func getMediaBase(videoname string) string {
	mediaRoot := "videos"
	return filepath.Join(mediaRoot, videoname)
}
