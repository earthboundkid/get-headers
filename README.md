# get-headers
Simple tool to show the headers from GET-ing a URL

The problem this solves is that when you use `curl -I` it does a `HEAD` request, potentially changing the result, and when you do `curl -i` it also dumps the page HTML on you. This does a `GET` and returns those results—including any doubled headers.
