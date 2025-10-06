# Subscription Service

![Go](https://img.shields.io/badge/Go-1.24-blue)
![Postgres](https://img.shields.io/badge/PostgreSQL-16-%234169E1)
![Docker](https://img.shields.io/badge/Docker-✔-2496ED)
![License](https://img.shields.io/badge/license-MIT-green)

Микросервис для управления подписками с REST API интерфейсом.

## Функциональность

- **CRUD операции** над подписками:
    - Создание подписки
    - Получение информации о подписке    
    - Получение всех подписок    
    - Обновление подписки    
    - Удаление подписки    
- **Расчет стоимости** подписок за период с фильтрацией
- **Валидация данных** и обработка ошибок
- **Документация API** через Swagger
## Технологии

- **Go 1.24** - основной язык программирования 
- **PostgreSQL** - база данных
- **Docker & Docker Compose** - контейнеризация
- **Gorilla Mux** - HTTP роутинг
- **pgx** - драйвер PostgreSQL
- **Migrate** - миграции базы данных

##  API Endpoints

### Подписки

- `POST /subscriptions` - Создать новую подписку
- `GET /subscriptions` - Получить все подписки
- `GET /subscriptions/{id}` - Получить подписку по ID
- `PUT /subscriptions/{id}` - Обновить подписки
- `DELETE /subscriptions/{id}` - Удалить подписку
### Расчеты

- `GET /subscriptions/sum` - Подсчет суммарной стоимости подписок
#### Параметры для расчета суммы:

- `id` (опционально) - фильтр по пользователю
- `service_name` (опционально) - фильтр по сервису
- `from` (обязательно) - начало периода (формат: MM-YYYY)  
- `to` (обязательно) - конец периода (формат: MM-YYYY)

## Структура проекта

```text
subscription/
├── cmd/
│   └── main.go              # Точка входа
├── internal/
│   ├── database/            # Подключение к БД
│   ├── handlers/            # HTTP обработчики
│   ├── model/               # Модели данных
│   ├── dto/                 # Data Transfer Objects
│   └── server/              # HTTP сервер
├── migrations/              # Миграции БД
├── docker-compose.yml       # Docker Compose
├── Dockerfile              # Docker образ
└── .env           # Пример переменных окружения
```

##  Docker контейнеры

Сервис состоит из трех контейнеров:
1. **postgres** - База данных PostgreSQL
2. **migrate** - Применение миграций БД
3. **app** - Основное приложение

## Разработка

### Локальный запуск без Docker

1. Установите зависимости:

```bash
go mod download
```

2. Запустите PostgreSQL локально
3. Запустите приложение:
```bash
go run cmd/main.go
```

### Миграции базы данных

Миграции автоматически применяются при запуске через Docker Compose. Для ручного применения:

```bash
docker-compose run migrate
```

##  Модель данных

### Subscription

```go

type Subscription struct {
    ID          string      `json:"id"`
    ServiceName string      `json:"service_name"`
    Price       int         `json:"price"`
    UserId      string      `json:"user_id"`
    StartDate   CustomDate  `json:"start_date"`
    EndDate     *CustomDate `json:"end_date,omitempty"`
}
```

### CustomDate

Кастомный тип для работы с датами в формате "MM-YYYY"

##  Примеры использования

### Создание подписки

```bash

curl -X POST http://localhost:9091/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Netflix",
    "price": 599,
    "user_id": "user123",
    "start_date": "01-2024"
  }'
```

### Расчет суммы подписок

```bash

curl "http://localhost:9091/subscriptions/sum?from=01-2024&to=12-2024"
```

##  Логирование

Приложение логирует ключевые события:

- Успешные операции
- Ошибки валидации
- Ошибки базы данных
- HTTP запросы

Логи выводятся в stdout контейнера.

##  Планы по развитию

- Добавить аутентификацию и авторизацию
- Реализовать пагинацию для списка подписок
- Добавить кэширование
- Реализовать уведомления об истекающих подписках
- Добавить метрики и мониторинг
- Создать клиентские SDK
## Лицензия
MIT © 2025 Mikebang-star