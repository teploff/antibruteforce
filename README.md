# Anti brute-force gRPC web service

<img align="right" width="160" src="static/img/gopher.png">

[![Build Status](https://travis-ci.com/teploff/antibruteforce.svg?branch=master)](https://travis-ci.com/github/teploff/antibruteforce)
[![codecov](https://codecov.io/gh/teploff/antibruteforce/branch/master/graph/badge.svg)](https://codecov.io/gh/teploff/antibruteforce)
[![Go Report Card](https://goreportcard.com/badge/github.com/teploff/antibruteforce)](https://goreportcard.com/report/github.com/teploff/antibruteforce)

Сервис предназначен для борьбы с подбором паролей при авторизации в какой-либо системе, вызывается перед авторизацией пользователя и может либо разрешить, либо заблокировать попытку.
Предполагается, что сервис используется только для server-server, т.е. скрыт от конечного пользователя.

## Содержание

1. [ Описание. ](#desc)
2. [ Конфигурирование Rate Limiter-а. ](#usage)
3. [ Сборка и запуск проекта. ](#build)
    1. [ Docker. ](#build-docker)
    2. [ Makefile. ](#build-makefile)

<a name="desc"></a>
## 1. Описание
Сервис ограничивает частоту попыток авторизации для различных комбинаций параметров, а именно:
- не более N попыток в T1 единицу времени для данного логина;
- не более M попыток в T2 единицу времени для данного пароля (защита от обратного brute-force);
- не более K попыток в T3 единицу времени для данного IP.

Для подсчета и ограничения частоты запросов, использовался алгоритм [leaky bucket](https://en.wikipedia.org/wiki/Rate_limiting). Реализовано множество bucket-ов, по одному на каждый логин/пароль/ip. Bucket-ы храниться в памяти, для каждого из типов bucket-ов (логин/пароль/ip) предусмотрено время протухания expire_time. Реализован Bucket GC, который с интервалом времени GCTime отслеживает протухшие бакеты и удаляет их из памяти, что позволяет избежать утечек памяти.

Разработан command-line интерфейс для ручного администрирования сервиса, через который существует возможность вызвать сброс бакета и управлять whitelist/blacklist-ами. CLI работает через HTTP интерфейс.

White/black листы содержат списки адресов сетей, которые обрабатываются более простым способом:
Если входящий ip в whitelist, то сервис безусловно разрешает авторизацию;
Если - в blacklist, то отклоняет.

<a name="usage"></a>
## 2. Конфигурирование Rate Limiter-а
Для того, чтобы ограничить частоту попыток авторизации для конкретного типа bucket'a, например, логина, необходимо перейти в раздел конфигурации rate_limiter -> login. 
- **rate** - количество попыток
- **interval** - интервал времени, на котором действует ограничение. Величина соотвествует типу time.Duration
- **expire_time** - интервал времени, через который bucket считается протухшим. Величина так же соотвествует типу time.Duration

Конфигурация для типов bucket'ов пароля и ip представляет идентичную стркутуру.

<a name="build"></a>
## 3. Сборка и запуск проекта

<a name="build-docker"></a>
## 3.1 Docker
Тут описание про докер
<a name="build-makefile"></a>
## 3.1 Makefile
Тут описание про makefile
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

<kbd>
    <p align="center">
      <img src="static/img/cli_help.png">
    </p>
</kbd>

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

