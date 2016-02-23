# get-headers [![GoDoc](https://godoc.org/github.com/carlmjohnson/get-headers?status.svg)](https://godoc.org/github.com/carlmjohnson/get-headers) [![Go Report Card](https://goreportcard.com/badge/github.com/carlmjohnson/get-headers)](https://goreportcard.com/report/github.com/carlmjohnson/get-headers)
Simple tool to show the headers from GET-ing a URL

The problem this solves is that when you use `curl -I` it does a `HEAD` request, potentially changing the result, and when you do `curl -i` it also dumps the page HTML on you. This does a `GET` and returns those results—including any doubled headers. It also (optionally) downloads the body of the page and returns speed and timing information.

## Installation
First install [Go](http://golang.org).

If you just want to install the binary to your current directory and don't care about the source code, run

```bash
GOBIN=. GOPATH=/tmp/gobuild go get github.com/carlmjohnson/get-headers
```

## Screenshots
```bash
$ get-headers https://www.example.com
GET https://www.example.com
HTTP/1.1 200 OK

Cache-Control      max-age=604800
Content-Type       text/html
Date               Fri, 15 Jan 2016 13:40:38 GMT
Etag               "359670651+gzip"
Expires            Fri, 22 Jan 2016 13:40:38 GMT
Last-Modified      Fri, 09 Aug 2013 23:54:35 GMT
Server             ECS (iad/18CB)
Vary               Accept-Encoding
X-Cache            HIT
X-Ec-Custom-Error  1

Time            100ms 204µs
Content length  1.2 KB
Speed           12.4 KB/s
```

```bash
$ get-headers -gzip http://www.example.com
GET http://www.example.com
HTTP/1.1 200 OK

Cache-Control      max-age=604800
Content-Encoding   gzip
Content-Length     606
Content-Type       text/html
Date               Fri, 15 Jan 2016 13:43:28 GMT
Etag               "359670651+gzip"
Expires            Fri, 22 Jan 2016 13:43:28 GMT
Last-Modified      Fri, 09 Aug 2013 23:54:35 GMT
Server             ECS (iad/182A)
Vary               Accept-Encoding
X-Cache            HIT
X-Ec-Custom-Error  1

Time            9ms 709µs
Content length  606
Speed           60.9 KB/s
```

```bash
$ get-headers -h
Usage of get-headers:

get-headers [opts] <url>...
        Gets the URLs and prints their headers alphabetically.
        Repeated headers are printed with an asterisk.

  -g	Shortcut for -gzip
  -gzip
    	Enable GZIP compression
  -i	Shortcut for -ignore-body
  -ignore-body
    	Ignore body of request; close connection after gettings the headers
```
