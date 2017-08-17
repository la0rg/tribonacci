[![Go Report Card](https://goreportcard.com/badge/github.com/la0rg/tribonacci)](https://goreportcard.com/report/github.com/la0rg/tribonacci)

# Tribonacci

Tribonacci is a RESTful web service written in Go that calculates the [Tribonacci numbers](http://oeis.org/wiki/Tribonacci_numbers).

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Docker 17.03
- Go 1.8.1

## Installing

To install the Tribonacci you need to download the sources, build docker image and then run it.

```
$ go get -d github.com/la0rg/tribonacci
$ cd $GOPATH/src/github.com/la0rg/tribonacci
$ docker build -t tribonacci .
$ docker run -it --rm -p 8080:8080 tribonacci:latest
```

## API

After the installation you can get your first tribonacci number by using
the only one available RESTful endpoint.

```
GET /tribonacci/{number}
```

If everything is fine you get response with status 200(OK) and body containing json object.

```
{
   "value": "44"
}
```

The service could also return 400 and 404.

Default request timeout is 1 minute.

## Algorithm

Calculation of the Tribonacci numbers is based on the iterative approach
with additional memorization.

This approach compared to recursive one takes less memory and could handle bigger numbers without "stackoverflow". Also, according to many perfomance tests it's faster.

Memorization in the form of caching also improves perfomance of this algorithm. Cache size could also be configured in accordance with available resources.