package crawler

import "testing"

func TestJobIsAddedToWorker(t *testing.T) {
	c := NewCrawler("http://example.com")
	worker := NewCrawlerWorker(c)
	worker.AddJob("/")
	if len(worker.jobs) < 1 {
		t.Errorf("Expected at least 1 job added to worker got %v", len(worker.jobs))
	}
}
