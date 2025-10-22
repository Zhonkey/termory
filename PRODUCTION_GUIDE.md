# Production Deployment Guide для Trainer API

## Обзор

Этот гайд описывает как собрать и запустить Go приложение в production с использованием Docker и миграциями БД.

## Структура файлов

```
├── Dockerfile.prod              # Production Dockerfile (multi-stage build)
├── docker-compose.prod.yml      # Production docker-compose конфигурация
├── Makefile                     # Команды для сборки и деплоя
├── scripts/deploy.sh            # Bash скрипт для деплоя
├── cmd/
│   ├── api/main.go             # API сервер
│   └── migrate/main.go          # Миграционный инструмент
└── cmd/migrate/migrations/      # SQL миграции
```

## Требования

- Docker >= 20.10
- Docker Compose >= 2.0
- Bash (для скрипта deploy.sh)

## Production сборка

### Способ 1: С Nginx (рекомендуется для production)

```bash
# Полный деплой с Nginx и SSL
docker-compose -f docker-compose.prod.nginx.yml build
docker-compose -f docker-compose.prod.nginx.yml up -d
docker-compose -f docker-compose.prod.nginx.yml exec app /app/migrate up
```

**Требует:**
- Скопировать SSL сертификаты в папку `ssl/`:
  - `ssl/cert.pem` - SSL сертификат
  - `ssl/key.pem` - приватный ключ
- Отредактировать `nginx.conf` если нужны другие домены

### Способ 2: Без Nginx (для разработки/тестирования)

```bash
# Полный деплой (сборка + запуск + миграции)
make deploy-prod

# Или пошагово:
make build-prod      # Сборка образа
make up-prod         # Запуск контейнеров
make migrate-prod    # Запуск миграций
```

### Способ 3: Использование скрипта deploy.sh

```bash
# Полный деплой
./scripts/deploy.sh up

# Остановка
./scripts/deploy.sh down

# Запуск миграций
./scripts/deploy.sh migrate

# Просмотр логов
./scripts/deploy.sh logs
./scripts/deploy.sh logs-db

# Проверка здоровья
./scripts/deploy.sh health
```

### Способ 4: Прямое использование docker-compose

```bash
# Сборка
docker-compose -f docker-compose.prod.yml build

# Запуск
docker-compose -f docker-compose.prod.yml up -d

# Миграции
docker-compose -f docker-compose.prod.yml exec app /app/migrate up

# Логи
docker-compose -f docker-compose.prod.yml logs -f app
```

## Локальная сборка (без Docker)

Если нужно собрать бинарники локально:

```bash
# Сборка бинарников
make build-local

# Запуск API (требует запущенной БД)
make run-local

# Запуск миграций
make migrate-local
```

## Конфигурация

### Переменные окружения

Production конфигурация находится в `docker-compose.prod.yml`:

```yaml
environment:
  - DB_DSN=user=trainer password=trainer host=db port=5432 dbname=trainer sslmode=disable
  - DB_DRIVER=pgx
  - PORT=8080
  - ENV=production
```

Для изменения параметров БД отредактируйте `docker-compose.prod.yml`.

## Миграции БД

### Запуск миграций

```bash
# Через Makefile
make migrate-prod

# Через скрипт
./scripts/deploy.sh migrate

# Через docker-compose
docker-compose -f docker-compose.prod.yml exec app /app/migrate up
```

### Откат миграций

```bash
# Через Makefile
make migrate-down-prod

# Через скрипт
./scripts/deploy.sh migrate-down

# Через docker-compose
docker-compose -f docker-compose.prod.yml exec app /app/migrate down
```

### Проверка статуса миграций

```bash
# Через Makefile
make migrate-status-prod

# Через скрипт
./scripts/deploy.sh status

# Через docker-compose
docker-compose -f docker-compose.prod.yml exec app /app/migrate status
```

## Логирование

### Просмотр логов API

```bash
# Через Makefile
make logs-prod

# Через скрипт
./scripts/deploy.sh logs

# Через docker-compose
docker-compose -f docker-compose.prod.yml logs -f app
```

### Просмотр логов БД

```bash
# Через Makefile
make logs-db-prod

# Через скрипт
./scripts/deploy.sh logs-db

# Через docker-compose
docker-compose -f docker-compose.prod.yml logs -f db
```

## Мониторинг

### Проверка здоровья приложения

```bash
./scripts/deploy.sh health
```

Проверяет:
- Статус контейнера API
- Статус контейнера БД
- Доступность API endpoint

### Healthcheck в Docker

Production Dockerfile включает healthcheck:

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
```

## Оптимизации Production

### Dockerfile.prod использует:

1. **Multi-stage build** - уменьшает размер финального образа
2. **Alpine Linux** - минимальный базовый образ (~5MB)
3. **Статическая компиляция** - `CGO_ENABLED=0`
4. **Оптимизация бинарника** - `-ldflags="-w -s"` (удаляет debug информацию)
5. **Непривилегированный пользователь** - безопасность
6. **Healthcheck** - автоматическое восстановление

### docker-compose.prod.yml использует:

1. **Restart policy** - `unless-stopped` для автоматического восстановления
2. **Health checks** - для БД и приложения
3. **Logging** - ограничение размера логов (10MB max)
4. **Networks** - изоляция контейнеров
5. **Volumes** - персистентность данных БД

## Размеры образов

```
Development: ~1.5GB (golang:latest)
Production:  ~40MB  (alpine:latest + бинарники без ca-certificates)
Nginx:       ~10MB  (nginx:alpine)
```

## Nginx конфигурация

### Структура файлов для Nginx:

```
├── docker-compose.prod.nginx.yml  # Docker Compose с Nginx
├── nginx.conf                      # Конфигурация Nginx
└── ssl/                            # SSL сертификаты
    ├── cert.pem                    # SSL сертификат
    └── key.pem                     # Приватный ключ
```

### Что делает Nginx:

1. **SSL/TLS** - обработка HTTPS соединений
2. **HTTP to HTTPS redirect** - перенаправление с HTTP на HTTPS
3. **Reverse proxy** - проксирование запросов к Go приложению
4. **Rate limiting** - ограничение количества запросов
5. **Gzip compression** - сжатие ответов
6. **Security headers** - добавление заголовков безопасности
7. **Static files** - раздача статических файлов

### Получение SSL сертификатов:

**Вариант 1: Let's Encrypt (бесплатно)**
```bash
# Установить certbot
sudo apt-get install certbot python3-certbot-nginx

# Получить сертификат
sudo certbot certonly --standalone -d yourdomain.com

# Скопировать в проект
mkdir -p ssl
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem ssl/cert.pem
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem ssl/key.pem
sudo chown $USER:$USER ssl/*
```

**Вариант 2: Самоподписанный сертификат (для тестирования)**
```bash
mkdir -p ssl
openssl req -x509 -newkey rsa:4096 -keyout ssl/key.pem -out ssl/cert.pem -days 365 -nodes
```

### Редактирование nginx.conf:

Измените `server_name _;` на ваш домен:
```nginx
server_name yourdomain.com www.yourdomain.com;
```

Если используете Let's Encrypt, обновите пути сертификатов:
```nginx
ssl_certificate /etc/nginx/ssl/cert.pem;
ssl_certificate_key /etc/nginx/ssl/key.pem;
```

## Безопасность

- ✅ Непривилегированный пользователь (appuser:1000)
- ✅ Минимальный базовый образ (Alpine)
- ✅ Отсутствие исходного кода в финальном образе
- ✅ Отсутствие debug информации в бинарниках
- ✅ Изолированная сеть Docker
- ✅ SSL/TLS шифрование (Nginx)
- ✅ Security headers (HSTS, X-Frame-Options и т.д.)
- ✅ Rate limiting на уровне Nginx
- ✅ Защита от доступа к скрытым файлам

## Troubleshooting

### Контейнер не запускается

```bash
# Проверить логи
./scripts/deploy.sh logs

# Проверить статус контейнеров
docker-compose -f docker-compose.prod.yml ps
```

### Миграции не выполняются

```bash
# Проверить статус БД
./scripts/deploy.sh logs-db

# Убедиться что БД готова
docker-compose -f docker-compose.prod.yml exec db pg_isready -U trainer
```

### Ошибка подключения к БД

Проверить переменные окружения в `docker-compose.prod.yml`:
- `DB_DSN` - строка подключения
- `DB_DRIVER` - драйвер БД (pgx)

### Очистка и переустановка

```bash
# Остановить контейнеры
./scripts/deploy.sh down

# Удалить volume с данными БД (ВНИМАНИЕ: потеря данных!)
docker volume rm trainer-postgres-prod

# Заново запустить
./scripts/deploy.sh up
```

## Примеры использования

### Полный деплой на чистой машине

```bash
./scripts/deploy.sh up
```

### Обновление приложения

```bash
# Остановить старую версию
./scripts/deploy.sh down

# Собрать новую версию
make build-prod

# Запустить новую версию с миграциями
./scripts/deploy.sh up
```

### Откат на предыдущую версию

```bash
# Откатить миграции
./scripts/deploy.sh migrate-down

# Остановить контейнеры
./scripts/deploy.sh down

# Переключиться на старую версию кода (git checkout)
git checkout <old-commit>

# Запустить старую версию
./scripts/deploy.sh up
```

## Дополнительные команды

```bash
# Показать справку
make help

# Показать справку скрипта
./scripts/deploy.sh help

# Подключиться к БД
docker-compose -f docker-compose.prod.yml exec db psql -U trainer -d trainer

# Выполнить команду в контейнере API
docker-compose -f docker-compose.prod.yml exec app /bin/sh

# С Nginx:
docker-compose -f docker-compose.prod.nginx.yml logs -f nginx
docker-compose -f docker-compose.prod.nginx.yml exec nginx nginx -t  # Проверка конфига
```

## Обновление SSL сертификата

```bash
# Получить новый сертификат от Let's Encrypt
sudo certbot renew

# Скопировать обновленный сертификат
sudo cp /etc/letsencrypt/live/yourdomain.com/fullchain.pem ssl/cert.pem
sudo cp /etc/letsencrypt/live/yourdomain.com/privkey.pem ssl/key.pem
sudo chown $USER:$USER ssl/*

# Перезагрузить Nginx
docker-compose -f docker-compose.prod.nginx.yml exec nginx nginx -s reload
```

## Рекомендации

1. **Используйте переменные окружения** для чувствительных данных (пароли, ключи)
2. **Регулярно проверяйте логи** для выявления проблем
3. **Делайте бэкапы БД** перед миграциями
4. **Тестируйте миграции** на staging перед production
5. **Используйте версионирование** для контроля изменений
6. **Мониторьте ресурсы** (CPU, память, диск)

## Дополнительные ресурсы

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Go Build Documentation](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies)
