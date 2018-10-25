package crawler

import (
	"regexp"
	"time"
)

const maxbuffSize = 5

var (
	LinksREGX = regexp.MustCompile("<a.*href=\"([^\"]*)\"[^>]*>")
	// XML specific details
	XMLheader = []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	XMLroot   = []byte("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">")
	XMLpre    = []byte("\n\t<url>\n\t\t<loc>")
	XMLcls    = []byte("</loc>\n\t</url>")
)

// CrawlerWorker represents a crawler job controller
type CrawlerWorker struct {
	JID     int
	crawler *Crawler
	jobs    []*CrawlerJob
	pending map[int]bool    // { 1: true } // pending job ids
	links   map[string]bool // visited links
}

// NewCrawlerWorker creates a new CrawlerWorker with existing crawler struct
func NewCrawlerWorker(c *Crawler) *CrawlerWorker {
	return &CrawlerWorker{
		1,
		c,
		make([]*CrawlerJob, 0),
		make(map[int]bool),
		make(map[string]bool),
	}
}

// AddJob creates and adds a new job for a given path
func (w *CrawlerWorker) AddJob(p string) {
	if _, ok := w.links[p]; ok {
		w.links[p] = true
		return
	}
	w.links[p] = true
	job := NewJob(w.JID, p)
	w.JID++ // increment worker job id
	w.jobs = append(w.jobs, job)
	w.pending[job.ID] = true
}

// Run processes all the buffered jobs
func (w *CrawlerWorker) Run() {
	for _, j := range w.jobs {
		go w.crawler.ProcessJob(*j)
	}
	w.ClearJobs()
}

// Start does the actual crawling on each pending jobs
func (w *CrawlerWorker) Start(done chan bool) {
	go func() {
		for {
			select {
			case job := <-w.crawler.Entries:
				w.Complete(job)
				for _, link := range job.Links {
					w.AddJob(link)
				}
				// crawl if worker buffer size if maxed out
				if len(w.jobs) >= maxbuffSize {
					w.Run()
				}
				if !w.HasPendingJobs() {
					done <- true
					close(w.crawler.Entries)
					return
				}
			case <-time.After(50 * time.Millisecond):
				w.Run()
			}
		}
	}()
}

// GetLinks returns all visted links
func (w *CrawlerWorker) GetLinks() []string {
	var links []string
	for link := range w.links {
		links = append(links, link)
	}
	return links
}

// Complete removes processed job from worker queue
func (w *CrawlerWorker) Complete(job CrawlerJob) {
	if len(w.pending) < 1 {
		return
	}
	delete(w.pending, job.ID)
}

// HasPendingJobs confirms if the worker still has jobs to process
func (w *CrawlerWorker) HasPendingJobs() bool { return len(w.pending) > 0 }

// ClearJobs clears the worker's jobs buffer
func (w *CrawlerWorker) ClearJobs() {
	w.jobs = w.jobs[:0]
}
