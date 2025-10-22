# OpenAPI и Swagger документация для Trainer API

## Обзор

Ваш API теперь имеет полную OpenAPI 2.0 (Swagger) документацию с интерактивным Swagger UI интерфейсом.

## Доступ к документации

### Swagger UI
Откройте в браузере: `http://localhost:8080/swagger`

Swagger UI предоставляет интерактивный интерфейс для:
- Просмотра всех доступных endpoints
- Чтения описания каждого endpoint
- Тестирования API прямо из браузера
- Просмотра примеров запросов и ответов

### OpenAPI JSON
Получить raw OpenAPI спецификацию: `http://localhost:8080/swagger.json`

## Структура документации

### Аутентификация

#### POST /auth/access_token
Получить access token для аутентификации

**Параметры:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Ответ:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### POST /auth/refresh_token
Обновить access token используя refresh token

**Параметры:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Ответ:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Управление пользователями (требуется роль admin)

#### GET /api/admin/users
Получить список всех пользователей

**Заголовки:**
```
Authorization: Bearer <access_token>
```

**Ответ:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "user",
    "created_at": "2025-10-22T10:00:00Z"
  }
]
```

#### PUT /api/admin/users
Создать нового пользователя

**Заголовки:**
```
Authorization: Bearer <access_token>
```

**Параметры:**
```json
{
  "email": "newuser@example.com",
  "password": "password123",
  "first_name": "Jane",
  "last_name": "Smith",
  "role": "user"
}
```

**Ответ:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "email": "newuser@example.com",
  "first_name": "Jane",
  "last_name": "Smith",
  "role": "user",
  "created_at": "2025-10-22T11:00:00Z"
}
```

#### GET /api/admin/users/{id}
Получить информацию о конкретном пользователе

**Заголовки:**
```
Authorization: Bearer <access_token>
```

**Параметры:**
- `id` (path): UUID пользователя

**Ответ:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "role": "user",
  "created_at": "2025-10-22T10:00:00Z"
}
```

#### POST /api/admin/users/{id}
Обновить информацию о пользователе

**Заголовки:**
```
Authorization: Bearer <access_token>
```

**Параметры:**
- `id` (path): UUID пользователя

**Тело запроса:**
```json
{
  "email": "newemail@example.com",
  "first_name": "John",
  "last_name": "Smith",
  "password": "newpassword123"
}
```

**Ответ:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "newemail@example.com",
  "first_name": "John",
  "last_name": "Smith",
  "role": "user",
  "created_at": "2025-10-22T10:00:00Z"
}
```

#### DELETE /api/admin/users/{id}
Удалить пользователя

**Заголовки:**
```
Authorization: Bearer <access_token>
```

**Параметры:**
- `id` (path): UUID пользователя

**Ответ:**
```json
{}
```

## Коды ответов

| Код | Описание |
|-----|----------|
| 200 | Успешный запрос (GET) |
| 201 | Ресурс успешно создан/обновлен |
| 400 | Неверные параметры запроса |
| 401 | Не авторизован (отсутствует или неверный token) |
| 403 | Недостаточно прав (требуется роль admin) |
| 404 | Ресурс не найден |
| 500 | Внутренняя ошибка сервера |

## Безопасность

### Bearer Token Authentication
Все защищенные endpoints требуют JWT token в заголовке `Authorization`:

```
Authorization: Bearer <your_access_token>
```

### Роли доступа
- **user**: Базовый пользователь
- **mentor**: Наставник (может иметь дополнительные права)
- **admin**: Администратор (полный доступ)

## Примеры использования

### cURL

#### Получить access token
```bash
curl -X POST http://localhost:8080/auth/access_token \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

#### Получить список пользователей
```bash
curl -X GET http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer <access_token>"
```

#### Создать нового пользователя
```bash
curl -X PUT http://localhost:8080/api/admin/users \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "password123",
    "first_name": "Jane",
    "last_name": "Smith",
    "role": "user"
  }'
```

### JavaScript/Fetch

```javascript
// Получить access token
const tokenResponse = await fetch('http://localhost:8080/auth/access_token', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'password123'
  })
});

const { access_token } = await tokenResponse.json();

// Получить список пользователей
const usersResponse = await fetch('http://localhost:8080/api/admin/users', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${access_token}`
  }
});

const users = await usersResponse.json();
console.log(users);
```

## Обновление документации

Документация находится в файле `internal/interfaces/http/docs/swagger.go`.

Для обновления документации:

1. Отредактируйте JSON спецификацию в `swagger.go`
2. Обновите описания endpoints, параметров и ответов
3. Перезагрузите приложение

## Интеграция с инструментами

### Postman
1. Откройте Postman
2. Нажмите "Import"
3. Выберите "Link"
4. Введите: `http://localhost:8080/swagger.json`
5. Нажмите "Continue"

### VS Code REST Client
Создайте файл `.http` с примерами:

```http
### Получить access token
POST http://localhost:8080/auth/access_token
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

### Получить список пользователей
GET http://localhost:8080/api/admin/users
Authorization: Bearer <access_token>
```

## Лучшие практики

1. **Всегда используйте HTTPS в production** - текущая конфигурация поддерживает оба протокола
2. **Храните tokens безопасно** - не передавайте их в URL
3. **Используйте refresh tokens** - access tokens имеют ограниченный срок действия
4. **Проверяйте роли** - убедитесь, что пользователь имеет необходимые права
5. **Обрабатывайте ошибки** - всегда проверяйте коды ответов

## Дополнительные ресурсы

- [OpenAPI 2.0 Specification](https://swagger.io/specification/v2/)
- [Swagger UI Documentation](https://swagger.io/tools/swagger-ui/)
- [JWT Authentication](https://jwt.io/)
