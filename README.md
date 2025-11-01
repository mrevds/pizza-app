# Pizza App 🍕

Микросервисное приложение для заказа и управления пиццей. Система состоит из нескольких сервисов, построенных на Go с использованием gRPC и REST API.

## 📋 Архитектура

```
┌─────────────────────┐
│   API Gateway       │  (HTTP REST на порту 8080)
│   (Reverse Proxy)   │
└──────────┬──────────┘
           │
      ┌────┴────┐
      │          │
┌─────▼──────┐  ┌──────▼──────┐
│User Service│  │Card Service │  (gRPC)
│(50051)     │  │(50052)      │
└─────┬──────┘  └──────┬──────┘
      │                │
┌─────▼──────┐  ┌──────▼──────┐
│ User DB    │  │  Card DB    │  (PostgreSQL)
│(54322)     │  │(5433)       │
└────────────┘  └─────────────┘
```

### Сервисы

#### 1. **API Gateway** (`./api-gateway`)
- REST API точка входа для клиентов
- Обратный прокси для маршрутизации запросов к микросервисам
- Аутентификация и авторизация
- **Порт**: 8080

#### 2. **User Service** (`./user-service`)
- Управление пользователями
- Аутентификация (регистрация, вход)
- JWT токены (access + refresh)
- gRPC интерфейс для других сервисов
- **Порт gRPC**: 50051
- **База данных**: PostgreSQL на порту 54322

#### 3. **Card Service** (`./card-service`)
- Управление платежными картами
- Хранение информации о картах пользователей
- gRPC интерфейс
- **Порт gRPC**: 50052
- **База данных**: PostgreSQL на порту 5433

---

## 🚀 Быстрый старт

### Предварительные требования

- Docker и Docker Compose
- Go 1.20+ (для локальной разработки)
- PostgreSQL client (опционально)

### Запуск приложения

```bash
# Перейди в корневую папку проекта
cd /home/denis/GolandProjects/pizza-app

# Запусти все сервисы через docker-compose
docker-compose up -d

# Проверь статус контейнеров
docker-compose ps

# Просмотри логи
docker-compose logs -f
```

### Остановка приложения

```bash
docker-compose down

# Удали также и волюмы с данными БД (если нужно сбросить данные)
docker-compose down -v
```

---

## 🔌 API Endpoints

### User Service API (через API Gateway)

**Базовый URL**: `http://localhost:8080`

#### Регистрация
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "phone_number": "+7 999 999 99 99",
  "password": "secure_password"
}
```

#### Вход
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "secure_password"
}
```

#### Получение профиля
```http
GET /api/v1/users/profile
Authorization: Bearer <ACCESS_TOKEN>
```

#### Обновление профиля
```http
PUT /api/v1/users/profile
Authorization: Bearer <ACCESS_TOKEN>
Content-Type: application/json

{
  "first_name": "Jane",
  "last_name": "Smith",
  "phone_number": "+7 888 888 88 88"
}
```

#### Обновление токена
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "<REFRESH_TOKEN>"
}
```

---

## 🗄️ База данных

### Подключение к БД

#### User Service DB
```bash
psql -h localhost -p 54322 -U user_db_user -d user_db
# Пароль: user_db_password
```

#### Card Service DB
```bash
psql -h localhost -p 5433 -U card_db_user -d card_db
# Пароль: card_db_password
```

### Миграции

Миграции выполняются автоматически при запуске сервиса через Goose.

**User Service миграции** (`user-service/migrations/`):
- `01_user.sql` - Таблица пользователей
- `02_add_phone_number.sql` - Добавление поля телефона
- `03_refresh_tokens.sql` - Таблица refresh токенов

**Card Service миграции** (`card-service/migrations/`):
- `01_card.sql` - Таблица платежных карт

---

## 🔐 Аутентификация

Приложение использует JWT токены для аутентификации:

- **Access Token**: Краткосрочный токен для доступа к защищённым ресурсам (15 минут)
- **Refresh Token**: Долгосрочный токен для получения нового access токена (7 дней)

Токены передаются в заголовке `Authorization` со схемой `Bearer`:
```
Authorization: Bearer <ACCESS_TOKEN>
```

---

## 🏗️ Структура проекта

```
pizza-app/
├── docker-compose.yaml          # Конфигурация контейнеров
├── README.md                    # Этот файл
│
├── api-gateway/                 # REST API Gateway
│   ├── cmd/api-gateway/main.go
│   ├── internal/
│   │   ├── app/
│   │   ├── handler/
│   │   ├── middleware/
│   │   └── router/
│   ├── Dockerfile
│   ├── go.mod
│   └── Makefile
│
├── user-service/                # User Service (gRPC)
│   ├── cmd/user-service/main.go
│   ├── api/user-service_v1/     # Proto definitions
│   ├── internal/
│   │   ├── app/
│   │   ├── config/
│   │   ├── entity/
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── repository/
│   │   ├── service/
│   │   └── utils/
│   ├── migrations/              # SQL миграции
│   ├── client/                  # DB клиент
│   ├── Dockerfile
│   ├── go.mod
│   └── Makefile
│
└── card-service/                # Card Service (gRPC)
    ├── api/user-card_v1/        # Proto definitions
    ├── internal/
    │   ├── config/
    │   ├── entity/
    │   ├── repository/
    │   └── migrations/
    ├── client/                  # DB клиент
    ├── Dockerfile
    ├── go.mod
    ├── config.yaml
    └── Makefile
```

---

## 🛠️ Разработка

### Локальный запуск сервиса

```bash
# User Service
cd user-service
make run

# API Gateway
cd api-gateway
make run

# Card Service
cd card-service
make run
```

### Сборка Proto файлов

```bash
# Сгенерировать Go код из .proto
cd user-service
make proto

cd ../card-service
make proto
```

### Просмотр логов

```bash
# Все сервисы
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f user-service
docker-compose logs -f card-service
docker-compose logs -f api-gateway
```

---

## 🔧 Конфигурация

### Environment переменные

Переменные устанавливаются в `docker-compose.yaml`:

#### User Service
```yaml
DB_HOST=user-db
DB_PORT=5432
DB_NAME=user_db
DB_USER=user_db_user
DB_PASSWORD=user_db_password
```

#### Card Service
```yaml
DB_HOST=card-db
DB_PORT=5432
DB_NAME=card_db
DB_USER=card_db_user
DB_PASSWORD=card_db_password
```

#### API Gateway
```yaml
APP_PORT=8080
APP_ENV=production
USER_SERVICE_HOST=user-service
USER_SERVICE_PORT=50051
CARD_SERVICE_HOST=card-service
CARD_SERVICE_PORT=50052
LOG_LEVEL=info
```

### YAML конфиги (опционально)

Некоторые сервисы могут использовать YAML конфиги:
- `user-service/config.yaml`
- `card-service/config.yaml`

---

## 🧪 Тестирование

### Используя Postman или Curl

**1. Регистрация пользователя:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "phone_number": "+7 999 999 99 99",
    "password": "password123"
  }'
```

**2. Вход в систему:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

**3. Получение профиля (требует токен):**
```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

---

## 🐛 Troubleshooting

### Ошибка: "failed to read dockerfile"
Убедись, что все Dockerfile файлы находятся в соответствующих директориях сервисов.

### Ошибка подключения к БД
Проверь, что все контейнеры запустились:
```bash
docker-compose ps
```

Проверь логи контейнера БД:
```bash
docker-compose logs user-db
docker-compose logs card-db
```

### Ошибка при миграции БД
Убедись, что миграции находятся в правильной директории и БД полностью инициализирована. Проверь логи:
```bash
docker-compose logs user-service
```

---

## 📚 Дополнительные ресурсы

- [gRPC Documentation](https://grpc.io/)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [Go Documentation](https://golang.org/doc/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Uber/FX Dependency Injection](https://pkg.go.dev/go.uber.org/fx)

---

## 👤 Автор

Проект разрабатывается в процессе обучения микросервисной архитектуре на Go.

---

## 📝 Лицензия

MIT License - используй как угодно для личных и коммерческих проектов.

