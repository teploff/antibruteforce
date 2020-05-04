[![Build Status](https://travis-ci.org/github/teploff/antibruteforce.png?branch=master)](https://travis-ci.org/github/teploff/antibruteforce)

```
git clone https://github.com/teploff/otus.git
cd otus/antibruteforce
go test -v $(go list ./... | grep -v /vendor/)
```