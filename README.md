# базовый функционал интернет-магазина

Административный API развивается в отдельном модуле [backoffice-executor](backoffice-executor/README.md), маршруты — `/backoffice/v1/...`.

## Go

Все модули и workspace используют Go **1.26.11**. Версия также указана в `.go-version` для менеджеров версий. В IDE нужно выбрать тот же SDK.

```sh
make go-version
make test
make build
make backoffice-run
```

Makefile использует `GOROOT/bin/go`, если переменная `GOROOT` задана, иначе `go` из PATH. Так бинарник Go и компилятор берутся из одного SDK. Путь можно переопределить через `make GO=/path/to/go test`; он должен соответствовать настроенному GOROOT.

При прямом запуске `go` из терминала также необходимо согласовать PATH и GOROOT: если GOROOT указан явно, добавить `$GOROOT/bin` в начало PATH. Настройки оболочки и глобальная установка Go не изменяются проектом.
