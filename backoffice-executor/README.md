# Backoffice executor

Отдельный Go-модуль административного API. Структура следует public-executor: `cmd`, `configs`, `init`, `api`, `internal/entrypoints/http`, `internal/middlewares` и `internal/usecases`.

Запуск из корня репозитория:

```sh
make backoffice-run
```

Требуется Go 1.26.11, как и для остальных модулей workspace. По умолчанию сервис слушает `127.0.0.1:9889`. Для изменения настроек скопировать `config.dist.yaml` в `backoffice-executor/config.yaml`. Без этого файла используются локальные значения по умолчанию. При согласованных PATH и GOROOT можно также запускать `go run ./cmd` из каталога модуля.

- `GET /backoffice/v1/health` — состояние HTTP-сервера.
- `GET /backoffice/v1/swagger` — Swagger UI (ресурсы UI загружаются с CDN).
- `GET /backoffice/v1/openapi.yaml` — встроенный OpenAPI-контракт.

Все будущие административные операции размещаются под `/backoffice/v1/...`, включая `/backoffice/v1/files/{uuid}/upload` и `/backoffice/v1/files/{uuid}/download`.

Генерация HTTP-контракта из корня репозитория:

```sh
make backoffice-openapi-gen
```

Пока спецификация хранится одним файлом и не требует bundling. Генератор закреплён на той же версии v2.4.1, что и сгенерированный код public-executor. Единственный инфраструктурный healthcheck зарегистрирован напрямую в Echo; текущий каркас не зависит от наличия api.gen.go. При добавлении прикладных операций подключается сгенерированный ServerInterface.

Модуль использует общие HTTP-конфигурацию и логгер из shared. Подключения PostgreSQL/S3 и прикладные сценарии добавляются по мере реализации фич; текущий каркас запускается без внешних сервисов. Аутентификация и административные права пока не подключены, поэтому сервис предназначен для локального запуска.
