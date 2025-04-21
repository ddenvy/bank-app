# Банковский REST API

REST API сервис для банковского приложения с поддержкой основных банковских операций.

## Функциональные возможности

- Регистрация и аутентификация пользователей
- Управление банковскими счетами
- Операции с картами (выпуск, просмотр)
- Денежные переводы
- Кредитные операции
- Финансовая аналитика
- Интеграция с ЦБ РФ и SMTP-сервисом

## Требования

- Go 1.23+
- Docker и Docker Compose
- Make (опционально)

## Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/your-username/bank-app.git
cd bank-app
```

2. Создайте файл .env в корневой директории:
```env
SERVER_ADDRESS=:8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/bank_app?sslmode=disable
JWT_SECRET=your-secret-key
SMTP_HOST=smtp.example.com
SMTP_USERNAME=your-username
SMTP_PASSWORD=your-password
```

3. Запустите базу данных в Docker:
```bash
make docker-up
```
или
```bash
docker-compose up -d
```

4. Примените миграции:
```bash
make migrate
```
или
```bash
docker exec -i bank_app_db psql -U postgres -d bank_app < migrations/001_init.sql
```

5. Установите зависимости и запустите сервер:
```bash
make init
make run
```
или
```bash
go mod download
go run cmd/api/main.go
```

## API Endpoints

### Публичные эндпоинты

#### Регистрация пользователя
```http
POST /api/v1/register
Content-Type: application/json

{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
}
```

#### Аутентификация пользователя
```http
POST /api/v1/login
Content-Type: application/json

{
    "email": "test@example.com",
    "password": "password123"
}

Response:
{
    "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Защищенные эндпоинты (требуют JWT-токен)

#### Счета
- `POST /api/v1/accounts` - Создание счета
```http
POST /api/v1/accounts
Authorization: Bearer <token>

Response:
{
    "id": 1,
    "number": "40817810099910004312",
    "balance": 0,
    "currency": "RUB"
}
```

- `GET /api/v1/accounts` - Получение списка счетов
- `GET /api/v1/accounts/{id}` - Получение информации о счете

#### Карты
- `POST /api/v1/cards` - Выпуск карты
```http
POST /api/v1/cards
Authorization: Bearer <token>
Content-Type: application/json

{
    "account_id": 1
}

Response:
{
    "id": 1,
    "number": "4276 1234 5678 9012",
    "expiry_date": "2028-03-14T00:00:00Z",
    "status": "active"
}
```

- `GET /api/v1/cards` - Получение списка карт
- `GET /api/v1/cards/{id}` - Получение информации о карте

#### Переводы
- `POST /api/v1/transfers` - Создание перевода
```http
POST /api/v1/transfers
Authorization: Bearer <token>
Content-Type: application/json

{
    "from_account": 1,
    "to_account": 2,
    "amount": 1000.00
}
```

#### Кредиты
- `POST /api/v1/credits` - Оформление кредита
```http
POST /api/v1/credits
Authorization: Bearer <token>
Content-Type: application/json

{
    "account_id": 1,
    "amount": 100000.00,
    "term": 12
}
```

- `GET /api/v1/credits/{id}/schedule` - Получение графика платежей

#### Аналитика
- `GET /api/v1/analytics` - Получение финансовой аналитики
- `GET /api/v1/accounts/{id}/predict` - Прогноз баланса

## Тестирование

### Unit-тесты

Для запуска всех тестов:
```bash
make test
```
или
```bash
go test ./...
```

Для запуска тестов конкретного пакета:
```bash
go test ./internal/service -v
```

### Моки

В проекте используется библиотека `testify/mock` для создания моков. Примеры тестов:

#### CreditService

```go
func TestCreditService_Create(t *testing.T) {
    // Подготовка
    mockCreditRepo := new(MockCreditRepository)
    mockAccountRepo := new(MockAccountRepository)
    cfg := &config.Config{}
    service := NewCreditService(mockCreditRepo, mockAccountRepo, cfg)

    // Настройка мока для проверки кредитного лимита
    existingCredits := []*model.Credit{
        {
            UserID: 1,
            Amount: 1000000,
            Status: "active",
        },
    }
    mockCreditRepo.On("GetByUserID", ctx, int64(1)).Return(existingCredits, nil)

    // Действие
    err := service.Create(ctx, 1, 1, 100000, 12)

    // Проверка
    assert.Error(t, err)
    assert.Equal(t, "credit load limit exceeded", err.Error())
}

func TestCreditService_calculateMonthlyPayment(t *testing.T) {
    tests := []struct {
        name     string
        amount   float64
        term     int
        rate     float64
        expected float64
    }{
        {
            name:     "кредит 500000 на 24 месяца под 12%",
            amount:   500000,
            term:     24,
            rate:     12,
            expected: 23536.74,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            payment := service.calculateMonthlyPayment(tt.amount, tt.term, tt.rate)
            assert.InDelta(t, tt.expected, payment, 0.01)
        })
    }
}
```

### Лучшие практики тестирования

1. Изоляция тестов
   - Каждый тест использует свои собственные моки
   - Моки создаются заново для каждого подтеста
   - Состояние не передается между тестами

2. Проверка граничных условий
   - Проверка лимитов (например, кредитный лимит)
   - Обработка ошибок (несуществующие счета)
   - Валидация входных данных

3. Точность вычислений
   - Правильное округление денежных значений
   - Корректные формулы расчета (аннуитетные платежи)
   - Использование float64 с округлением до копеек

4. Документирование тестов
   - Понятные названия тестов
   - Комментарии к сложным проверкам
   - Структура "подготовка-действие-проверка"

### Ручное тестирование

Для ручного тестирования API можно использовать Postman или curl:

1. Регистрация:
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
```

2. Вход:
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

3. Создание счета (с токеном):
```bash
curl -X POST http://localhost:8080/api/v1/accounts \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Разработка

### Структура проекта

```
bank-app/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   └── handlers.go
│   ├── model/
│   │   └── models.go
│   ├── repository/
│   │   ├── interfaces.go
│   │   ├── postgres.go
│   │   ├── user_repository.go
│   │   ├── account_repository.go
│   │   ├── card_repository.go
│   │   ├── credit_repository.go
│   │   ├── transfer_repository.go
│   │   └── analytics_repository.go
│   └── service/
│       ├── interfaces.go
│       ├── service.go
│       ├── user_service.go
│       ├── account_service.go
│       ├── card_service.go
│       ├── credit_service.go
│       ├── transfer_service.go
│       └── analytics_service.go
├── migrations/
│   └── 001_init.sql
├── docker-compose.yml
├── Makefile
├── .env
├── go.mod
└── README.md
```

### Команды Make

- `make init` - Инициализация проекта (установка зависимостей, запуск БД, миграции)
- `make build` - Сборка приложения
- `make run` - Запуск приложения
- `make test` - Запуск тестов
- `make docker-up` - Запуск Docker контейнеров
- `make docker-down` - Остановка Docker контейнеров
- `make migrate` - Применение миграций
- `make lint` - Запуск линтера

## Безопасность

- Все пароли хешируются с использованием bcrypt
- Данные карт шифруются с помощью PGP
- Используется HMAC для проверки целостности данных
- Все критические операции требуют JWT-аутентификации

## Лицензия

MIT 