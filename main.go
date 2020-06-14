package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	crawl "./crawler"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/crawl", crawlHandler).Methods("POST")
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./"))))
	log.Println("Starting app server...")
	log.Fatal(http.ListenAndServe(":3000", r))
}

func crawlHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	host := r.FormValue("host")
	c := crawl.NewCrawler(host)
	worker := crawl.NewCrawlerWorker(c)
	worker.AddJob("/")
	done := make(chan bool, 1)
	worker.Start(done)
	<-done // block till worker is done
	close(done)

	// write links to file
	loc, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fileLoc := path.Join(loc, "sitemap.xml")
	err = crawl.WriteXML(fileLoc, host, worker.GetLinks())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<h2>Sitemap generated, <a href='./sitemap.xml'>View</a></h2>`)
}
