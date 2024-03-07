package main

import (
	"fmt"
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
	r.Get("/media/{videoname:[\\w-]+}/stream/{segment:[\\w-]+.ts}", streamContinue)
	r.Get("/thumbnail/{thumbnail:[\\w]+.jpg}/", getThumbnail)
	r.Get("/downloads/{thumbnail:[\\w]+.jpg}/", getThumbnail)
	// r.Get("/profile/{profile:^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$}/", getProfileImage)
	r.Get("/profile/{profile:[\\w]+}/", getProfileImage)
	r.Get("/banner/{profile:[\\w]+}/", getBannerImage)

	err := http.ListenAndServe(":"+PORT, r)
	fmt.Println("%v", err)

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

// Assuming thumbnails are in jpeg format and in the same folder as its corresponding hls files -> videos/videotitle
func getThumbnail(w http.ResponseWriter, req *http.Request) {

	thumbnailname := chi.URLParam(req, "thumbnail")

	mediabase := getMediaBase(strings.Split(thumbnailname, ".jpg")[0])
	thumbnailFile := filepath.Join(mediabase, thumbnailname)

	w.Header().Set("Content-Type", "image/JPEG")
	http.ServeFile(w, req, thumbnailFile)
}

func getBannerImage(w http.ResponseWriter, req *http.Request) {
	profile := chi.URLParam(req, "profile")
	profileFile := filepath.Join("banner", profile+".jpg")
	w.Header().Set("Content-Type", "image/JPEG")
	http.ServeFile(w, req, profileFile)
}

// Assuming profiles are in jpeg format and in a folder called profiles
func getProfileImage(w http.ResponseWriter, req *http.Request) {
	//try to get the profile image

	profile := chi.URLParam(req, "profile")

	profileFile := filepath.Join("profiles", profile+".jpg")
	w.Header().Set("Content-Type", "image/JPEG")
	http.ServeFile(w, req, profileFile)
}

func getMediaBase(videoname string) string {
	mediaRoot := "videos"
	return filepath.Join(mediaRoot, videoname)
}
