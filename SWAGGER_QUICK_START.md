# Быстрый старт с Swagger UI

## Запуск приложения

```bash
# Перейти в директорию проекта
cd /home/zhonkey/GolandProjects/trainer

# Запустить приложение
go run cmd/api/main.go
```

Приложение запустится на `http://localhost:8080`

## Доступ к Swagger UI

1. Откройте браузер
2. Перейдите на `http://localhost:8080/swagger`
3. Вы увидите интерактивный интерфейс со всеми доступными endpoints

## Первые шаги

### 1. Получить Access Token

В Swagger UI:
1. Найдите endpoint `POST /auth/access_token`
2. Нажмите "Try it out"
3. Введите в поле body:
```json
{
  "email": "admin@example.com",
  "password": "admin123"
}
```
4. Нажмите "Execute"
5. Скопируйте значение `access_token` из ответа

### 2. Авторизоваться в Swagger UI

1. Нажмите кнопку "Authorize" в верхней части страницы
2. В поле "value" введите: `Bearer <ваш_access_token>`
3. Нажмите "Authorize"
4. Закройте диалог

Теперь все защищенные endpoints будут автоматически использовать ваш token.

### 3. Получить список пользователей

1. Найдите endpoint `GET /api/admin/users`
2. Нажмите "Try it out"
3. Нажмите "Execute"
4. Вы увидите список всех пользователей

### 4. Создать нового пользователя

1. Найдите endpoint `PUT /api/admin/users`
2. Нажмите "Try it out"
3. Введите в поле body:
```json
{
  "email": "newuser@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "role": "user"
}
```
4. Нажмите "Execute"
5. Скопируйте ID созданного пользователя из ответа

### 5. Получить информацию о пользователе

1. Найдите endpoint `GET /api/admin/users/{id}`
2. Нажмите "Try it out"
3. В поле "id" введите ID пользователя из предыдущего шага
4. Нажмите "Execute"

### 6. Обновить пользователя

1. Найдите endpoint `POST /api/admin/users/{id}`
2. Нажмите "Try it out"
3. В поле "id" введите ID пользователя
4. Введите в поле body:
```json
{
  "email": "updated@example.com",
  "first_name": "Jane",
  "last_name": "Smith"
}
```
5. Нажмите "Execute"

### 7. Удалить пользователя

1. Найдите endpoint `DELETE /api/admin/users/{id}`
2. Нажмите "Try it out"
3. В поле "id" введите ID пользователя
4. Нажмите "Execute"

## Использование с cURL

### Получить access token
```bash
curl -X POST http://localhost:8080/auth/access_token \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }' | jq .
```

### Сохранить token в переменную
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/auth/access_token \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }' | jq -r '.access_token')

echo $TOKEN
```

### Получить список пользователей
```bash
curl -X GET http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer $TOKEN" | jq .
```

### Создать пользователя
```bash
curl -X PUT http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe",
    "role": "user"
  }' | jq .
```

## Использование с Postman

### Импортировать OpenAPI спецификацию

1. Откройте Postman
2. Нажмите "Import"
3. Выберите вкладку "Link"
4. Введите: `http://localhost:8080/swagger.json`
5. Нажмите "Continue"
6. Нажмите "Import"

Postman автоматически создаст коллекцию со всеми endpoints.

### Настроить Bearer Token

1. В Postman откройте коллекцию "Trainer API"
2. Перейдите на вкладку "Authorization"
3. Выберите тип "Bearer Token"
4. В поле "Token" введите ваш access token

Теперь все запросы в коллекции будут использовать этот token.

## Использование с VS Code REST Client

### Установить расширение

1. Откройте VS Code
2. Перейдите в Extensions (Ctrl+Shift+X)
3. Найдите "REST Client"
4. Установите расширение от Huachao Mao

### Создать файл с примерами

Создайте файл `api.http` в корне проекта:

```http
@baseUrl = http://localhost:8080
@token = 

### Получить access token
POST {{baseUrl}}/auth/access_token
Content-Type: application/json

{
  "email": "admin@example.com",
  "password": "admin123"
}

### Сохранить token (скопируйте значение access_token из ответа выше)
# @token = <скопируйте_access_token_отсюда>

### Получить список пользователей
GET {{baseUrl}}/api/admin/users
Authorization: Bearer {{token}}

### Создать пользователя
PUT {{baseUrl}}/api/admin/users
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "email": "newuser@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "role": "user"
}

### Получить пользователя по ID
GET {{baseUrl}}/api/admin/users/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer {{token}}

### Обновить пользователя
POST {{baseUrl}}/api/admin/users/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "email": "updated@example.com",
  "first_name": "Jane",
  "last_name": "Smith"
}

### Удалить пользователя
DELETE {{baseUrl}}/api/admin/users/550e8400-e29b-41d4-a716-446655440000
Authorization: Bearer {{token}}
```

Для использования:
1. Откройте файл `api.http`
2. Нажмите "Send Request" над каждым запросом
3. Результаты появятся в новой вкладке

## Структура Swagger документации

### Файлы

- `internal/interfaces/http/docs/swagger.go` - OpenAPI спецификация в JSON формате
- `internal/interfaces/http/router.go` - Обработчики для `/swagger` и `/swagger.json` endpoints

### Как обновить документацию

1. Отредактируйте JSON в `internal/interfaces/http/docs/swagger.go`
2. Обновите:
   - Описания endpoints в секции `paths`
   - Параметры в секции `parameters`
   - Примеры ответов в секции `responses`
   - Модели данных в секции `definitions`
3. Перезагрузите приложение

## Решение проблем

### Swagger UI не загружается

- Убедитесь, что приложение запущено на `http://localhost:8080`
- Проверьте консоль браузера на ошибки (F12)
- Очистите кэш браузера (Ctrl+Shift+Delete)

### Token не работает

- Убедитесь, что вы скопировали полный token
- Проверьте, что token не истек
- Используйте endpoint `/auth/refresh_token` для получения нового token

### Ошибка 403 Forbidden

- Убедитесь, что ваш пользователь имеет роль `admin`
- Проверьте, что вы авторизованы с правильным token

### Ошибка 401 Unauthorized

- Убедитесь, что вы передали token в заголовке `Authorization`
- Проверьте формат: `Bearer <token>`
- Получите новый token через `/auth/access_token`

## Дополнительные ресурсы

- [Swagger UI Documentation](https://swagger.io/tools/swagger-ui/)
- [OpenAPI 2.0 Specification](https://swagger.io/specification/v2/)
- [REST Client VS Code Extension](https://marketplace.visualstudio.com/items?itemName=humao.rest-client)
- [Postman Documentation](https://learning.postman.com/)
