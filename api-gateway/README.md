# API Gateway

Prodакшн-grade API Gateway для микросервисной архитектуры pizza-app.

## Структура проекта

```
api-gateway/
├── cmd/
│   └── api-gateway/
│       └── main.go              # Entry point
├── internal/
│   ├── app/
│   │   └── app.go               # DI контейнер (uber/fx)
│   ├── client/
│   │   └── grpc_client.go       # gRPC клиенты
│   ├── config/
│   │   └── config.go            # Конфигурация из .env
│   ├── handler/
│   │   ├── user_handler.go      # User Service handlers
│   │   └── card_handler.go      # Card Service handlers
│   ├── middleware/
│   │   └── middleware.go        # CORS, Auth, Logging, RequestID
│   ├── logger/
│   │   └── logger.go            # Структурированный логирование (zap)
│   └── router/
│       └── router.go            # HTTP маршруты (gorilla/mux)
├── Dockerfile
├── docker-compose.yaml
├── go.mod
├── go.sum
├── .env                         # Конфигурация
└── README.md
```

## Архитектура

```
┌─────────────┐
│   Клиент    │ HTTP REST (POST, GET, PUT, DELETE)
└──────┬──────┘
       │
┌──────▼─────────────────────────────┐
│         API Gateway                │
│  • HTTP Server (port 8080)         │
│  • CORS + Auth Middleware          │
│  • Request logging                 │
│  • Error handling                  │
└──────┬──────────────────────────────┘
       │ gRPC
       ├──────────────────────┬──────────────────────┐
       │                      │                      │
   ┌───▼──────┐          ┌────▼──────┐         ┌────▼──────┐
   │   User   │          │   Card    │         │  Order    │
   │ Service  │          │ Service   │         │ Service   │
   │(50051)   │          │(50052)    │         │(50053)    │
   └───┬──────┘          └────┬──────┘         └────┬──────┘
       │                      │                      │
       └──────────┬───────────┴──────────┬──────────┘
                  │
             ┌────▼─────┐
             │PostgreSQL│
             └──────────┘
```

## Маршруты

### Открытые (без аутентификации)

```bash
POST   /health                    # Health check
POST   /api/v1/auth/register      # Регистрация
POST   /api/v1/auth/login         # Вход
POST   /api/v1/auth/refresh       # Обновить токен
```

### Защищённые (требуют Bearer токен)

```bash
# User Service
GET    /api/v1/user/profile       # Получить профиль
PUT    /api/v1/user/profile       # Обновить профиль
POST   /api/v1/user/logout        # Выход

# Card Service
GET    /api/v1/cards              # Все карты пользователя
POST   /api/v1/cards              # Добавить карту
GET    /api/v1/cards/balance      # Баланс карты
POST   /api/v1/cards/deposit      # Пополнение
POST   /api/v1/cards/withdraw     # Снятие
POST   /api/v1/cards/transfer     # Перевод
```

## Запуск

### Локально

```bash
cd api-gateway

# Установить зависимости
go mod download

# Запустить
go run cmd/api-gateway/main.go
```

API Gateway будет доступен на `http://localhost:8080`

### В Docker

```bash
# Из корня проекта
docker-compose up -d api-gateway

# Проверить логи
docker-compose logs -f api-gateway
```

## Примеры запросов

### Регистрация

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "phone_number": "+1234567890",
    "password": "securepassword"
  }'
```

### Получить профиль

```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <access_token>"
```

### Получить карты

```bash
curl -X GET "http://localhost:8080/api/v1/cards?user_id=123" \
  -H "Authorization: Bearer <access_token>"
```

### Пополнить карту

```bash
curl -X POST http://localhost:8080/api/v1/cards/deposit \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "card_id": 1,
    "user_id": 123,
    "amount": 100.00,
    "description": "Deposit"
  }'
```

## Middleware

1. **CORS** - разрешает запросы с других хостов
2. **RequestID** - добавляет уникальный ID для трейсирования
3. **Logging** - логирует все HTTP запросы
4. **Auth** - проверяет Bearer токен (для защищённых маршрутов)

## Конфигурация

Переменные окружения в `.env`:

```env
APP_NAME=api-gateway
APP_PORT=8080
APP_ENV=development
USER_SERVICE_HOST=localhost
USER_SERVICE_PORT=50051
CARD_SERVICE_HOST=localhost
CARD_SERVICE_PORT=50052
LOG_LEVEL=info
```

## Следующие шаги

1. **Реализовать вызовы gRPC** в обработчиках
2. **Добавить валидацию** входных данных
3. **Добавить авторизацию** (проверку роли пользователя)
4. **Добавить rate limiting** (защита от DDoS)
5. **Добавить caching** (Redis)
6. **Метрики** (Prometheus)
7. **Трейсирование** (Jaeger)

