# Crawl
A simple sitemap scrapper written in Go. Crawls a given url `u` and writes all
links (on same domain) to a `sitemap.xml` file.

## Installation
Via `go get`

```
$ go get github.com/codehakase/crawl
```

Run http server:
```
$ go build main.go

$ ./main # navigate to http://localhost:3000 to input url to crawl
```

Run from command line:
```
$ go run main.go https://example.com
```
