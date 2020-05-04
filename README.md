[![Build Status](https://travis-ci.com/teploff/antibruteforce.svg?branch=master)](https://travis-ci.com/github/teploff/antibruteforce)

```
git clone https://github.com/teploff/otus.git
cd otus/antibruteforce
go test -v $(go list ./... | grep -v /vendor/)
```