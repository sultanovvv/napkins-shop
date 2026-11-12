# Backoffice executor

Отдельный Go-модуль административного API. Структура следует public-executor: `cmd`, `configs`, `init`, `api`, `internal/entrypoints/http`, `internal/middlewares` и `internal/usecases`.

Запуск из корня репозитория:

```sh
cp backoffice-executor/config.dist.yaml backoffice-executor/config.yaml
# Указать параметры существующей локальной PostgreSQL в config.yaml.
make backoffice-run
```

Требуется Go 1.26.11, как и для остальных модулей workspace. По умолчанию сервис слушает `127.0.0.1:9889`. Секция `postgres` обязательна; HTTP-настройки и таймаут остановки имеют значения по умолчанию. Если config.yaml уже существует, добавить в него секцию postgres из config.dist.yaml, сохранив остальные настройки. При согласованных PATH и GOROOT можно также запускать `go run ./cmd` из каталога модуля.

- `GET /backoffice/v1/health` — состояние HTTP-сервера.
- `GET /backoffice/v1/swagger` — Swagger UI (ресурсы UI загружаются с CDN).
- `GET /backoffice/v1/openapi.yaml` — встроенный OpenAPI-контракт.

Все будущие административные операции размещаются под `/backoffice/v1/...`, включая `/backoffice/v1/files/{uuid}/upload` и `/backoffice/v1/files/{uuid}/download`.

Генерация HTTP-контракта из корня репозитория:

```sh
make backoffice-openapi-gen
```

Исходная спецификация разбита на файлы:

```text
api/
├── openapi.yaml                 # корневая спецификация со ссылками
├── components/
│   ├── enums/                   # перечисления
│   ├── parameters/              # переиспользуемые параметры запросов
│   ├── paths/                   # операции, сгруппированные по пути
│   ├── responses/               # переиспользуемые HTTP-ответы
│   └── schemas/                 # модели запросов и ответов
├── openapi.gen.yaml             # собранная спецификация для Swagger и codegen
└── api.gen.go                   # сгенерированные Go-типы и ServerInterface
```

Редактируются openapi.yaml и файлы в components. Перечисления из enums регистрируются в components.schemas корневой спецификации: отдельного раздела components.enums в OpenAPI нет. Ссылки в файлах относительные.

`make backoffice-openapi-gen` сначала собирает спецификацию через Redocly CLI v2.53.3, затем генерирует Go-код через oapi-codegen v2.6.0 (runtime v1.4.0). Для сборки YAML нужны Node.js и npm. Отдельная команда `make backoffice-openapi-bundle` обновляет только openapi.gen.yaml. Обе операции также выполняются через `go generate ./backoffice-executor/api` из корня workspace.

В сервис встраивается openapi.gen.yaml: Swagger получает единый документ без внешних ссылок. Оба сгенерированных файла включены в репозиторий и не редактируются вручную. Health и категории зарегистрированы через сгенерированный ServerInterface.

Модуль использует общие HTTP-конфигурацию, логгер и core-db из shared. PostgreSQL подключается и проверяется при запуске, соединение закрывается через Fx после остановки HTTP-сервера. Миграции автоматически не применяются: нужна существующая схема каталога (0001–0002). S3 пока не подключён. Аутентификация и административные права пока не подключены, поэтому сервис предназначен для локального запуска.

## Категории

| Метод | Путь | Результат |
| --- | --- | --- |
| POST | /backoffice/v1/categories | 201, созданная категория |
| GET | /backoffice/v1/categories | 200, плоский массив категорий |
| GET | /backoffice/v1/categories/{id} | 200, категория |
| PATCH | /backoffice/v1/categories/{id} | 200, обновлённая категория |
| DELETE | /backoffice/v1/categories/{id} | 204, без тела |

Категория содержит id, name, slug, parentId (null для корня) и sortOrder. Отдельного поля level нет. При создании slug необязателен; при PATCH переименование не меняет slug. Отсутствующие поля сохраняют значения, явный null для name/slug/sortOrder отклоняется. Null parentId сохраняет родителя; clearParentId=true переносит категорию в корень. Передать одновременно clearParentId=true и непустой parentId нельзя.

DELETE удаляет выбранную категорию и всех её потомков одной транзакцией. Связанные товары сохраняются, их category_id становится NULL. Это корень каталога без отдельной записи категории в БД. Товары и категории вне выбранной ветки не затрагиваются. Новая миграция не нужна: используются существующие ON DELETE CASCADE для категорий и ON DELETE SET NULL для товаров.

POST/PATCH принимают один JSON-объект с Content-Type: application/json. Неизвестные поля, некорректный JSON и ID возвращают 400 VALIDATION_ERROR. Обработчики возвращают ошибки usecase в общий error handler.

В `backoffice.http` есть последовательный сценарий CRUD, проверки slug, дерева, PATCH, ошибок и каскадного удаления трёх уровней. ID сохраняются автоматически, в конце удаляются созданные сценарием категории. Сценарий рассчитан на локальную БД.

## Ошибки API

Глобальный Echo HTTPErrorHandler подключён через Fx. HTTP-обработчик возвращает ошибку usecase напрямую; самостоятельно формировать JSON ошибки не нужно.

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Request validation failed",
    "details": [
      {
        "code": "INVALID_VALUE",
        "message": "cannot be used together with parentId",
        "target": "clearParentId"
      }
    ]
  }
}
```

`error`, `code` и `message` обязательны. `details` — необязательный массив; каждый элемент содержит code, message и target. Если ошибок отдельных полей нет, details отсутствует. Для управления интерфейсом используются code/target, а не разбор message.

Пакет `internal/customerr` предоставляет конструкторы с сигнатурой `(message string, err error)`: `NewInternal`, `NewValidation`, `NewNotFound`, `NewForbidden`, `NewConflict`, `NewUnavailable`. Сообщение предназначено для клиента, исходная ошибка сохраняется для логов и доступна через Unwrap. Конструктор определяет код ошибки; детали полей заполняются в месте формирования ошибки.

```go
return customerr.NewConflict("Slug already exists", err)
```

| Конструктор | HTTP | code |
| --- | --- | --- |
| NewInternal | 500 | INTERNAL_ERROR |
| NewValidation | 400 | VALIDATION_ERROR |
| NewNotFound | 404 | NOT_FOUND |
| NewForbidden | 403 | FORBIDDEN |
| NewConflict | 409 | CONFLICT |
| NewUnavailable | 503 | SERVICE_UNAVAILABLE |

Предметные коды, например SLUG_ALREADY_EXISTS или CATEGORY_CYCLE, передаются в details. Usecase категории формирует прикладные ошибки при валидации и преобразовании ошибок провайдера. Общий HTTP-handler знает только коды customerr, без зависимостей от usecase.

Ошибки без customerr в цепочке errors.As всегда возвращаются как 500 INTERNAL_ERROR с общим сообщением, включая echo.HTTPError. Для ожидаемых ошибок HTTP-слой должен явно использовать соответствующий конструктор customerr. Сообщения исходных ошибок клиенту не передаются; причины 5xx логируются с request ID. Для HEAD возвращается только статус, уже начатый ответ не перезаписывается.

Контракт описан в компонентах Error/ErrorBody/ErrorDetail из `api/components/schemas` и ErrorResponse из `api/components/responses`; новые операции должны ссылаться на общий ответ ошибки.
