package main

import (
	"encoding/base64"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"path/filepath"
	"os"

	qrcode "github.com/skip2/go-qrcode"
)

var urlStore = make(map[string]string)

const BASE_URL = "http://localhost:8009"

func getExecutableDir() string {
    exePath, err := os.Executable()
    if err != nil {
        log.Fatal(err)
    }
    return filepath.Dir(exePath)
}

func main() {
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/r/", handleRedirect)

	log.Println("Starting server on port 8009")
	log.Fatal(http.ListenAndServe("0.0.0.0:8009", nil))
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
	tmpl := template.Must(template.ParseFiles(htmlPath))

	if r.Method == http.MethodGet {
		tmpl.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		shortenedCode := generateShortURL(6)
		var png []byte
		png, err := qrcode.Encode(BASE_URL+"/r/"+shortenedCode, qrcode.Medium, 256)

		if err != nil {
			log.Println(err)
		}

		data := ShortenRequest{
			OriginalURL: r.FormValue("url"),
			ShortURL:    BASE_URL + "/r/" + shortenedCode,
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
