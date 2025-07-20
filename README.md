
# SubscriptionService

Сервис для управления подписками пользователей. Является тестовым заданием, если вы ищите какое-то готовое решение то скорее всего это не оно

## Требования

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Сборка и запуск с помощью Docker Compose

1. Клонируйте репозиторий:
   ```sh
   git clone https://github.com/alexputin/subscriptions.git
   cd subscriptions
   ```

2. Создайте файл окружения `test.env` (или используйте уже существующий):
   ```env
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_HOSTNAME=db
   DB_PORT=5432
   DB_NAME=postgres
   SERVER_ADDRESS=0.0.0.0:3000
   ENVIRONMENT=test
   ```

3. Запустите сервисы:
   ```sh
   docker compose up --build
   ```

4. После запуска сервис будет доступен по адресу: [http://localhost:1337](http://localhost:1337)

5. Swagger-документация доступна по адресу: [http://localhost:1337/swagger/index.html](http://localhost:1337/swagger/index.html)

## Миграции

Миграции применяются автоматически при запуске контейнера.

## Остановка

Для остановки и удаления контейнеров выполните:
```sh
docker compose down
```

## Локальная сборка (без Docker)

1. Установите Go 1.24+ и PostgreSQL.
2. Создайте файл `.env` с переменными окружения (см. `test.env`).
3. Выполните миграции:
   ```sh
   make migrate-up
   ```
4. Соберите и запустите приложение:
   ```sh
   make build
   ./tmp/subscriptions
   ```

## Контакты

Автор: Александр Путин
