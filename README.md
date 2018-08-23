<p align="center"><img src="https://habrastorage.org/webt/te/ys/9a/teys9a3vkkbf7watmp8vh-2ar4s.png"></p>

[![GoDoc](https://godoc.org/github.com/furdarius/date?status.svg)](https://godoc.org/github.com/furdarius/date)
[![Build Status](https://travis-ci.org/furdarius/date.svg?branch=master)](https://travis-ci.org/furdarius/date)
[![Go Report Card](https://goreportcard.com/badge/github.com/furdarius/date)](https://goreportcard.com/report/github.com/furdarius/date)

# Date it easy

Small, user-friendly library to help you work with dates in Go.

## Install
```
go get github.com/furdarius/date
```

### Adding as dependency by "go dep"
```
$ dep ensure -add github.com/furdarius/date
```

## Usage
```go

start, _ := date.Parse("2018-10-15")
end, _ := date.Parse("2018-10-20")

interval1 := date.Range{start, end}
interval2 := date.Range{start, end.AddDays(-1)}

base := []date.Range{
	date.Range{date.Date{2018, 10, 1}, date.Date{2018, 10, 5}},
	date.Range{date.Date{2018, 10, 8}, date.Date{2018, 10, 22}},
}

ranges := date.RangeSet(base).Sub(interval).Impose(interval2).ExtendEnd().List()
```

## Contributing

Pull requests are very much welcomed. Make sure a test or example is included that covers your change and
your commits represent coherent changes that include a reason for the change.

Use `gometalinter` to check code with linters:
```
gometalinter -t --vendor
```
