package main

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/codehakase/monzo_web_crawler_exercise/crawler"
)

func main() {
	s := time.Now()
	host := os.Args[1]
	log.Println("host: ", host)
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
		log.Fatal(err)
		return
	}
	fileLoc := path.Join(loc, "sitemap.xml")
	err = crawler.WriteXML(fileLoc, host, worker.GetLinks())
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("ran in: ", time.Since(s))
}
