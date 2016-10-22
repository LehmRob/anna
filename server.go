package main

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"rsc.io/letsencrypt"
)

const (
	upload    = "upload.html"
	result    = "result.html"
	bootstrap = "bootstrap"
)

// AnnaServer holds important structure for the server
type AnnaServer struct {
	server     *http.Server
	webrootDir string
	group      []ZipCodeGroup
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

	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(filepath.Join(anna.webrootDir,
		"static", bootstrap)))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		anna.showIndex(w, r)
	})
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		anna.upload(w, r)
	})
	mux.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		anna.showResult(w, r)
	})

	mux.HandleFunc("/export.csv", func(w http.ResponseWriter, r *http.Request) {
		anna.export(w, r)
	})

	var m letsencrypt.Manager
	if err := m.CacheFile("letsencrypt.cache"); err != nil {
		log.Fatal(err)
	}

	anna.server = &http.Server{
		Addr:    port,
		Handler: mux,
		TLSConfig: &tls.Config{
			GetCertificate: m.GetCertificate,
		},
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
		http.Redirect(w, r, "/", 302)
		return
	}

	log.Println("Get new upload")

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
	zipCodeGroups, err := analyzer.Analyze()
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Error occured")
		return
	}
	a.group = zipCodeGroups

	tmpl, err := template.ParseFiles(filepath.Join(a.webrootDir, result))
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Error occured")
		return
	}

	tmpl.Execute(w, zipCodeGroups)
}

func (a *AnnaServer) export(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Dispostion", "inline; filename=\"data.csv\"")
	w.Header().Add("Content-Type", "text/comma-separated-values")

	zipCode := r.FormValue("zip")
	log.Printf("Export file for zip code %s", zipCode)

	if a.group == nil {
		log.Printf("No data analyzed. Redirect to start")
		http.Redirect(w, r, "/", 302)
		return
	}

	for _, zipCodeGroup := range a.group {
		if zipCodeGroup.ZipCode == zipCode {
			err := zipCodeGroup.writeCsv()
			if err != nil {
				log.Println(err.Error())
				fmt.Fprintf(w, "Error occured")
				return
			}
			break
		}
	}

	f, err := os.Open(filepath.Join("/tmp", zipCode+".csv"))
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintf(w, "Error occured")
		return
	}
	defer f.Close()

	io.Copy(w, f)
}

// Run runs the server and returns only if an error occurs
func (a *AnnaServer) Run() error {
	return a.server.ListenAndServeTLS("", "")
}
