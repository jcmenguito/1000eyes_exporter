# 1000eyes_exporter

Prometheus exporter ot export test metrics and alerts from [ThousandEyes](https://www.thousandeyes.com/).
The port 9350 was chosen because someone already [reserved](https://github.com/prometheus/prometheus/wiki/Default-port-allocations) it for a ThousandEyes exporter that was supposed to be coming soon but has not been published up to November 2018.

## Building and Installing
###Build Pre-Req:
- Golang installed with `$GOBIN` and `$GOPATH` set
- [Gox](https://github.com/mitchellh/gox) -- Uses Gox for cross-platform builds
```shell
$ go get github.com/mitchellh/gox
...
$ gox -h
...
```

Makefile Targets:

| Target | Description|
|---|---|
|deps | Download and Install any missing dependecies|
|build  | Install missing dependencies. Builds binaries for linux and darwin in ./dist |
|tidy |                    Verifies and downloads all required dependencies|
|fmt   |                   Runs gofmt on all source files|
|clean |                   Removes build, dist and report dirs|
|debug  |                  Print make env information|

# Environment Settings
Mandatory 

- `ENV VAR " THOUSANDEYES_BEARER_TOKEN"` 

or 

- `ENV VAR "THOUSANDEYES_BASIC_AUTH_USER"` &&
- `ENV VAR "THOUSANDEYES_BASIC_AUTH_TOKEN"`

set to a valid ThousandEyes token to be able to query.

# Program Arguments

- `-GetBGP=true [true|false (default)]` if you want BGP test data collected
- `-GetHTTP=true [true|false (default)]` if you want HTTP request test data collected (false is default if not set)
- `-GetHttpMetrics=true [true|false (default)]` if you want HTTP routing test data collected (false is default if not set)
- `-GetNetPathViz=true [true|false (default)]` if you want Network Visualization Tests (false is default if not set)

<p>HINT: please be aware of the API request limit per minute .. if you have many tests and collect all details it's 
pretty sure that you're going to it.</p> 

- Just for debugging purpose: `-RetrospectionPeriod` You can set the period of time it queries into the past, e.g. `-RetrospectionPeriod 12h`. Large values do not make much sense, because we do not get data about when they started or ended. Just that they existed.