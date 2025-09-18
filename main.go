package main

import (
	"encoding/base64"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	qrcode "github.com/skip2/go-qrcode"
)

var urlStore = make(map[string]string)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getExecutableDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}
	return filepath.Dir(exe)
}

func getBaseURL() string {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost"
	}
	return baseURL
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	return ":" + port
}

func main() {
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/r/", handleRedirect)

	port := getPort()
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

type ShortenRequest struct {
	OriginalURL string
	ShortURL    string
	QRCode      string
}

func generateShortURL(length int) string {
	set := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortCode := make([]byte, length)
	for i := range shortCode {
		shortCode[i] = set[rand.Intn(len(set))]
	}
	return string(shortCode)
}

// func generateHashedURL(length int, url string) string {
// 	hash := sha256.Sum256([]byte(url))
// 	return hex.EncodeToString(hash[:])[:6]
// }

func handleForm(w http.ResponseWriter, r *http.Request) {
	htmlPath := filepath.Join(getExecutableDir(), "index.html")
	tmpl, err := template.ParseFiles(htmlPath)
	if err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		shortenedCode := generateShortURL(6)
		baseURL := getBaseURL()
		var png []byte
		png, err := qrcode.Encode(baseURL+"/r/"+shortenedCode, qrcode.Medium, 256)

		if err != nil {
			log.Println(err)
		}

		data := ShortenRequest{
			OriginalURL: r.FormValue("url"),
			ShortURL:    baseURL + "/r/" + shortenedCode,
			QRCode:      base64.StdEncoding.EncodeToString(png),
		}
		urlStore[shortenedCode] = data.OriginalURL

		tmpl.Execute(w, data)
		return
	}
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/r/")

	checkURL, ok := urlStore[shortCode]
	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, checkURL, http.StatusFound)
}
