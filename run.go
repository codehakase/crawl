package main

import (
	"log"
	"os"
	"path"
	"time"

	crawl "./crawler"
)

func main() {
	s := time.Now()
	host := os.Args[1]
	if host == "" {
		log.Fatalf("URL argurment is required")
	}
	log.Println("host: ", host)
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
		log.Fatal(err)
		return
	}
	fileLoc := path.Join(loc, "sitemap.xml")
	err = crawl.WriteXML(fileLoc, host, worker.GetLinks())
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("ran in: ", time.Since(s))
}
