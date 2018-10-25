package crawler

import (
	"net/http"
	"reflect"
	"testing"
)

func TestGetJobLinksFromReqBody(t *testing.T) {
	expected := []string{}
	res, err := http.Get("http://example.com")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	links := getJobLinksFromReqBody(res.Body, int(res.ContentLength))
	if !reflect.DeepEqual(links, expected) {
		t.Errorf("links: <%T> %+v not equals <%T> %+v", links, links, expected, expected)
	}
}

func TestGetLinksFromHTMLString(t *testing.T) {
	html := `
	<html>
	<div><a href="/blog">Blog</a></div>
	<div><a href="/contact">Contact us</a></div>
	</html>
	`
	expected := []string{"/blog", "/contact"}
	links := getlinks(html)

	if !reflect.DeepEqual(links, expected) {
		t.Errorf("expected %v not matched", links)
	}
}
