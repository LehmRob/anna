package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	upload = "upload.html"
	result = "result.html"

	maxMemory = 6 * 1024 * 1024
)

// AnnaServer holds important structure for the server
type AnnaServer struct {
	server     *http.Server
	webrootDir string
}

// NewServer creates a new instance of AnnaServer
func NewServer(port string) (*AnnaServer, error) {
	anna := new(AnnaServer)

	curDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	anna.webrootDir = filepath.Join(curDir, "webroot")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		anna.showIndex(w, r)
	})
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		anna.upload(w, r)
	})
	mux.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		anna.showResult(w, r)
	})

	anna.server = &http.Server{
		Addr:    port,
		Handler: mux,
	}

	return anna, nil
}

func (a *AnnaServer) showIndex(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(filepath.Join(a.webrootDir, upload))
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Error occured")
		return
	}
	defer file.Close()

	io.Copy(w, file)

}

func (a *AnnaServer) upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, "Please goto /")
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Error occured")
		return
	}
	defer file.Close()

	tmpFile, err := os.Create("/tmp/data.csv")
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Error occured")
		return
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, file)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Error occured")
		return
	}

	log.Printf("Successfully uploaded %s", header.Filename)
	http.Redirect(w, r, "/result", 302)
}

func (a *AnnaServer) showResult(w http.ResponseWriter, r *http.Request) {
	analyzer := NewCsvAnalyzer("/tmp/data.csv")
	results, err := analyzer.Analyze()
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Error occured")
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join(a.webrootDir, result))
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Error occured")
		return
	}

	tmpl.Execute(w, results)
}

// Run runs the server and returns only if an error occurs
func (a *AnnaServer) Run() error {
	return a.server.ListenAndServe()
}
