# Subscription Service

[![Go](https://img.shields.io/badge/Go-1.24-blue)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-%234169E1)](https://www.postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-✔-2496ED)](https://docker.com)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

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

### **Backend**
- **Go 1.24** - основной язык 
- **Chi Router** - высокопроизводительный HTTP роутер
- **Swagger** - автоматическая документация API

### **Data Layer** 
- **PostgreSQL 16** - реляционная база данных
- **pgx/pgxpool** - драйвер PostgreSQL с пулом соединений

### **Infrastructure**
- **Docker & Docker Compose** - контейнеризация и оркестрация
- **Database Migrations** 
- **Environment Configuration** - управление настройками через переменные окружения

##  API Endpoints

### Подписки

- `POST /subscriptions` - Создать новую подписку
- `GET /subscriptions` - Получить все подписки
- `GET /subscriptions/{id}` - Получить подписку по ID
- `PUT /subscriptions/{id}` - Обновить подписку
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
│   └── main.go                 # Точка входа
├── docs           
│   ├── swagger.json                
│   └── swagger.yaml           
├── internal/
│   ├── api/
│   │   ├── handlers/           # HTTP обработчики
│   │   ├── dto/                # Data Transfer Objects  
│   │   └── server/             # HTTP сервер
│   ├── database/               # Подключение к БД
│   └── model/                  # Модели данных (сущности БД)
├── migrations/                 # Миграции БД
├── docker-compose.yml          # Docker Compose
├── Dockerfile                  # Docker образ
└── .env                        # Переменные окружения
```
##  Быстрый старт

1. Клонировать репозиторий
```bash
git clone https://github.com/MikebangSfilya/subscription.git
cd subscription
```

2. Запусти через Docker
```bash
docker-compose up -d
```
3. API будет доступно на http://localhost:9091


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
    "user_id": "a37a0327-99af-4e62-8b33-55dc3863cdc6",
    "start_date": "01-2024"
  }'
```

### Расчет суммы подписок

```bash
curl "http://localhost:9091/subscriptions/sum?from=01-2024&to=12-2024"
```





