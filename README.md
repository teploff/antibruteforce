# Anti brute-force gRPC web service

<img align="right" width="160" src="static/img/gopher.png">

[![Build Status](https://travis-ci.com/teploff/antibruteforce.svg?branch=master)](https://travis-ci.com/github/teploff/antibruteforce)
[![codecov](https://codecov.io/gh/teploff/antibruteforce/branch/master/graph/badge.svg)](https://codecov.io/gh/teploff/antibruteforce)
[![Go Report Card](https://goreportcard.com/badge/github.com/teploff/antibruteforce)](https://goreportcard.com/report/github.com/teploff/antibruteforce)

## Предназначение
Сервис предназначен для борьбы с подбором паролей при авторизации в какой-либо системе.
Сервис вызывается перед авторизацией пользователя и может либо разрешить, либо заблокировать попытку.
Предполагается, что сервис используется только для server-server, т.е. скрыт от конечного пользователя.

Сервис ограничивает частоту попыток авторизации для различных комбинаций параметров, например:

не более N попыток/минуту для данного логина.
не более M попыток/минуту для данного пароля (защита от обратного brute-force).
не более K попыток/минуту для данного IP.
Для подсчета и ограничения частоты запросов, использовался алгоритм [leaky bucket](https://en.wikipedia.org/wiki/Rate_limiting).

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

CLI-админка для Anti Brute-Force сервиса

Сборка и запуск опции help
```shell script
make build_cli
./cli -help
```
![](static/img/cli_help.png)

Варианты использования

Если необходимо указать явно destination адрес brute-force сервиса, необходимо запускать binary-файл с флагом --dest
```shell script
./cli --dest 192.168.130.132:80
```

Запуск команд

Команду можно запускать по полному имени, например:
```shell script
./cli reset_bucket_by_login teploff
```
Так и по ее alias:
```shell script
./cli rbl teploff
```

Полное описание команд можно увидеть выше или в разделе "COMMANDS", вызвав опцию --help или -h при запуске binary-файла.

