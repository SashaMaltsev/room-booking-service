# Room Booking Service

Backend-приложение для бронирования переговорок. Сервис поддерживает создание переговорок и расписаний администратором, просмотр доступных слотов, создание и отмену броней пользователем, а также выдачу JWT через `/dummyLogin`.

API описано в [api.yaml](./api.yaml). Базовый адрес сервиса после запуска: `http://localhost:8080`.

## Что реализовано

Основные возможности:

- `POST /dummyLogin` с JWT, содержащим `user_id` и `role`
- `GET /_info` и `GET /`
- создание переговорок только для `admin`
- создание единственного расписания переговорки только для `admin`
- просмотр списка переговорок для `admin` и `user`
- просмотр доступных слотов переговорки по дате
- создание брони только для `user`
- идемпотентная отмена своей брони только для `user`
- просмотр всех броней для `admin` с пагинацией
- просмотр будущих броней текущего пользователя
- Docker Compose для запуска приложения и PostgreSQL
- unit-тесты и E2E-тесты на сценарии бронирования и отмены

## Стек

- Go 1.22
- PostgreSQL 16
- стандартный `net/http`
- JWT
- Docker / Docker Compose

## Архитектура

Структура проекта разделена по слоям:

- `cmd/api` - точка входа
- `internal/app` - сборка приложения и wiring зависимостей
- `internal/transport/http` - HTTP-обработчики, middleware, DTO
- `internal/service` - прикладная бизнес-логика
- `internal/domain` - доменные сущности, ошибки и интерфейсы
- `internal/repository/postgres` - работа с PostgreSQL
- `migrations` - SQL-миграции
- `tests/e2e` - end-to-end сценарии

## Ключевые решения

### Генерация слотов

Выбран ленивый подход: слоты не создаются заранее на большое окно дат, а генерируются при обращении к конкретной дате.

Как это работает:

- администратор один раз создаёт расписание переговорки
- при запросе `GET /rooms/{roomId}/slots/list?date=YYYY-MM-DD` сервис строит ожидаемые слоты для этой даты на основе расписания
- перед сохранением используются ограничения БД, поэтому повторная генерация не создаёт дубликаты
- если расписания для переговорки нет, слоты не создаются и переговорка считается недоступной

Почему выбран именно такой вариант:

- это хорошо подходит для сценария чтения на ближайшие 7 дней
- не нужно держать фоновый планировщик и заранее генерировать много лишних записей
- сохраняются стабильные UUID слотов, потому что после первой генерации они лежат в БД
- высокая нагрузка на чтение упирается в конкретную комнату и дату, а не в полный пересчёт календаря

### Защита инвариантов

Часть бизнес-ограничений enforced на уровне базы данных:

- уникальный email пользователя
- только одно расписание на переговорку
- только одна активная бронь на слот
- запрет пересекающихся слотов одной переговорки через `EXCLUDE USING GIST`
- фиксированная длительность слота `30 minutes`

Такой подход уменьшает риск гонок между конкурентными запросами.

### Авторизация

В базовом сценарии используется `POST /dummyLogin`.

Фиксированные пользователи:

- `admin` -> `00000000-0000-0000-0000-000000000001`
- `user` -> `00000000-0000-0000-0000-000000000002`

JWT содержит минимум:

- `user_id`
- `role`

## Бизнес-правила

- все даты и время хранятся и отдаются в UTC
- администратор не может создавать брони
- пользователь не может отменять чужую бронь
- повторная отмена уже отменённой брони возвращает `200 OK`
- бронирование слота в прошлом запрещено
- `GET /bookings/my` возвращает только будущие брони
- если у комнаты нет расписания, доступных слотов у неё нет

## Запуск

### Вариант 1. Через Docker Compose

`docker-compose` автоматически читает `.env` из корня проекта.

```bash
cp .env.example .env
docker-compose up -d --build
```

После запуска будут подняты:

- PostgreSQL
- миграции
- API на `localhost:8080`

Проверка:

```bash
curl http://localhost:8080/_info
```

Остановка:

```bash
docker-compose down -v
```

Если у вас используется новый синтаксис Docker, те же команды можно запускать как `docker compose ...`.

### Вариант 2. Локально без Docker для API

Нужен запущенный PostgreSQL и применённые миграции.

```bash
cp .env.example .env
docker-compose up -d db
docker-compose run --rm migrate
go run ./cmd/api
```

## Переменные окружения

Пример находится в [.env.example](./.env.example).

Основные переменные:

- `HTTP_PORT=8080`
- `HTTP_SHUTDOWN_TIMEOUT=10s`
- `POSTGRES_HOST=localhost`
- `POSTGRES_PORT=5432`
- `POSTGRES_USER=postgres`
- `POSTGRES_PASSWORD=postgres`
- `POSTGRES_DB=app`
- `POSTGRES_SSLMODE=disable`
- `JWT_SECRET=dev-secret-change-me`
- `JWT_TTL=12h`

Также можно целиком переопределить DSN через `DATABASE_DSN`.

## Примеры запросов

Получить токен администратора:

```bash
curl -X POST http://localhost:8080/dummyLogin \
  -H 'Content-Type: application/json' \
  -d '{"role":"admin"}'
```

Получить токен пользователя:

```bash
curl -X POST http://localhost:8080/dummyLogin \
  -H 'Content-Type: application/json' \
  -d '{"role":"user"}'
```

Создать переговорку:

```bash
curl -X POST http://localhost:8080/rooms/create \
  -H "Authorization: Bearer <admin-token>" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Blue Room","description":"2 floor","capacity":6}'
```

Получить слоты комнаты на дату:

```bash
curl "http://localhost:8080/rooms/<room-id>/slots/list?date=2026-03-26" \
  -H "Authorization: Bearer <user-token>"
```

Создать бронь:

```bash
curl -X POST http://localhost:8080/bookings/create \
  -H "Authorization: Bearer <user-token>" \
  -H 'Content-Type: application/json' \
  -d '{"slotId":"<slot-id>"}'
```

## Тесты

Все тесты:

```bash
go test ./...
```

Сборка:

```bash
go build ./...
```

Покрытие:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

E2E-тесты против живого сервиса:

```bash
RUN_E2E=1 BASE_URL=http://localhost:8080 go test ./tests/e2e -v
```

## Ограничения и упрощения

- ссылка на конференцию в брони пока не интегрирована с внешним сервисом
- `Makefile`, Swagger-генерация, CI и нагрузочное тестирование не добавлены
- `POST /register` и `POST /login` пока не подключены как HTTP-эндпоинты
- в текущей версии основной сценарий авторизации рассчитан на `/dummyLogin`

