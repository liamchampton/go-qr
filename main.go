package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/divan/qrlogo"
	"github.com/gorilla/mux"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	gracefulShutdown()

	var dir string = "html"
	r := mux.NewRouter()

	r.PathPrefix("/html/").Handler(http.StripPrefix("/html/", http.FileServer(http.Dir(dir)))).Methods("GET")
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/create-qrcode", createQrCode).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	currentTime := time.Now()
	fmt.Println("QR Code Generator server is started on port :8080...", currentTime.Format("Mon 02 Jan 2006 03:04pm"))
	log.Fatal(srv.ListenAndServe())
}

func home(w http.ResponseWriter, r *http.Request) {
	filename := "./html/index.html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprint(w, err)
	}

	fmt.Fprintf(w, string(body))
}

func createQrCode(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	imgBase64 := r.FormValue("image")
	url := r.FormValue("url")

	imgBase64cleaned := imgBase64[len("data:image/png;base64,"):len(imgBase64)]

	imgBytes, _ := base64.StdEncoding.DecodeString(imgBase64cleaned)

	img, _, _ := image.Decode(bytes.NewReader(imgBytes))

	filename := RandomString(16, charset)
	currentDir, _ := os.Getwd()
	logoFile := currentDir + "/html/uploaded_logos/" + filename + ".png"
	qrCodeFile := currentDir + "/html/qr_codes/" + filename + ".png"
	imgFile, err := os.Create(logoFile)
	if err != nil {
		panic(err)
	}
	png.Encode(imgFile, img)

	file, err := os.Open(logoFile)
	errcheck(err, "Failed to open logo:")
	defer file.Close()

	logo, _, err := image.Decode(file)
	errcheck(err, "Failed to decode PNG with logo:")

	qr, err := qrlogo.Encode(url, logo, 2000)
	errcheck(err, "Failed to encode QR:")

	out, err := os.Create(qrCodeFile)
	errcheck(err, "Failed to open output file:")
	out.Write(qr.Bytes())
	out.Close()

	writeResponse(w, http.StatusCreated, map[string]string{"fname": filename + ".png"})

	fmt.Println("QR image generated: ", filename+".png")
}

func RandomString(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func writeResponse(w http.ResponseWriter, code int, message map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	var jsonData []byte
	jsonData, _ = json.Marshal(message)
	fmt.Fprint(w, string(jsonData))
}

func errcheck(err error, str string) {
	if err != nil {
		fmt.Println(err, str)
		os.Exit(1)
	}
}
func gracefulShutdown() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}
