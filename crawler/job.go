package crawler

import (
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// CrawlerJob represents a Job with a path to craw
type CrawlerJob struct {
	ID    int
	Path  string
	Links []string
}

type Crawler struct {
	Host    string
	Entries chan CrawlerJob
}

// NewJob creates a new Job with a set path
func NewJob(id int, path string) *CrawlerJob {
	return &CrawlerJob{ID: id, Path: path}
}

// NewCrawler creates a crawler for a given host(domain)
func NewCrawler(host string) *Crawler {
	return &Crawler{
		host,
		make(chan CrawlerJob),
	}
}

// ProcessJob processes the job path and extracts links from it
func (c *Crawler) ProcessJob(job CrawlerJob) {
	// set link
	l := c.Host + job.Path
	log.Println("hitting: ", l)
	// validate link
	res, err := http.Get(l)
	if err != nil {
		log.Printf("link access err: %v\n", err)
		return
	}
	defer res.Body.Close()
	nbyts := res.ContentLength
	job.Links = getJobLinksFromReqBody(res.Body, int(nbyts))
	c.Entries <- job
}

func getJobLinksFromReqBody(r io.Reader, cl int) []string {
	if cl > 0 {
		// read n bytes from from request body
		byts := make([]byte, cl, cl) // create bytes of size <nbyts>
		_, err := io.ReadAtLeast(r, byts, cl)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		links := getlinks(string(byts))
		ilinks := make([]string, 0)
		// select only internal links
		for _, link := range links {
			if len(link) > 0 && link[0] == '/' {
				ilinks = append(ilinks, html.UnescapeString(link))
			}
		}
		return ilinks
	} else {
		// read from string(res.Body)
		byts, err := ioutil.ReadAll(r)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		links := getlinks(string(byts))
		// select only internal links
		ilinks := make([]string, 0)
		for _, link := range links {
			if len(link) > 0 && link[0] == '/' {
				ilinks = append(ilinks, html.UnescapeString(link))
			}
		}
		return ilinks
	}

}

func getlinks(str string) []string {
	links := make([]string, 0)
	// get regex matches
	matches := LinksREGX.FindAllStringSubmatch(str, -1)
	for _, l := range matches {
		links = append(links, l[1])
	}
	return links
}
