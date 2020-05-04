
```
git clone https://github.com/teploff/otus.git
cd otus/antibruteforce
go test -v $(go list ./... | grep -v /vendor/)
```