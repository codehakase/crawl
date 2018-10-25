package main

import (
	"log"
	"net/http"
	"os"
	"path"

	"github.com/codehakase/monzo_web_crawler_exercise/crawler"
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
	c := crawler.NewCrawler(host)
	worker := crawler.NewCrawlerWorker(c)
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
	err = crawler.WriteXML(fileLoc, host, worker.GetLinks())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sitemap generated, view => <a href='./sitemap.xml'>sitemap.xml</a>"))
}
