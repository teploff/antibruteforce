[![Build Status](https://travis-ci.com/teploff/antibruteforce.svg?branch=master)](https://travis-ci.com/github/teploff/antibruteforce)
[![codecov](https://codecov.io/gh/teploff/antibruteforce/branch/master/graph/badge.svg)](https://codecov.io/gh/teploff/antibruteforce)
[![Go Report Card](https://goreportcard.com/badge/github.com/teploff/antibruteforce)](https://goreportcard.com/report/github.com/teploff/antibruteforce)

Сборка проекта (Docker)
```
git clone https://github.com/teploff/antibruteforce.git
cd deployments/stage
docker-compose up -d --build && docker image prune -f 
```

Заверешение проекта (Docker)
```
docker-compose down
docker system prune --volumes
```