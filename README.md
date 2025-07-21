# 🛒 VK-internship - marketplace API

![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-8+-DC382D?logo=redis)
![Docker](https://img.shields.io/badge/Dockerlogo=docker)
![Swagger](https://img.shields.io/badge/Swagger?logo=swagger)

## 📋 Оглавление
- [🌟 Обзор проекта](#-обзор-проекта)
- [🚀 Основные возможности](#-основные-возможности)
- [⚙️ Технологический стек](#️-технологический-стек)
- [🏗️ Структура проекта](#️-структура-проекта)
- [🚀 Быстрый старт](#-быстрый-старт)
- [🔧 Использование API](#-использование-api)


## 🌟 Обзор проекта

Маркетплейс API - это высокопроизводительное RESTful приложение, реализующее функционал условного маркетплейса.

Проект решает следующие задачи:
- 🔐 Безопасная аутентификация и регистрация пользователей
- 📢 Управление объявлениями (создание, редактирование, удаление)
- 📋 Отображение ленты объявлений с фильтрацией и сортировкой
- 📚 Полная документация API через Swagger UI

## 🚀 Основные возможности

### 🔐 Аутентификация пользователей
- Регистрация новых пользователей с валидацией данных
- Аутентификация по логину/паролю с выдачей JWT токена
- Ролевая модель доступа (публичный/авторизованный доступ)
- Защита эндпоинтов middleware авторизации

### 📢 Управление объявлениями
- Создание объявлений с заголовком, описанием, изображением и ценой
- Редактирование и удаление объявлений (только для автора)
- Валидация данных объявления (длина текста, формат цены и URL)
- Автоматическое обновление кэша при изменении объявлений

### 📋 Лента объявлений
- Постраничный вывод объявлений с пагинацией
- Сортировка по дате создания и цене (возрастание/убывание)
- Фильтрация по диапазону цен
- Определение принадлежности объявления текущему пользователю
- Кэширование популярных запросов для ускорения ответа

### ⚙️ Дополнительные функции
- Подробное логирование всех операций
- Конфигурирование через переменные окружения
- Автоматическое применение миграций БД
- Полная документация API через Swagger UI
- Готовые Docker-образы для быстрого развертывания

## ⚙️ Технологический стек

- **Язык программирования**: Go 1.24+
- **База данных**: PostgreSQL 15
- **Кэширование**: Redis 8
- **Веб-фреймворк**: Chi Router
- **Аутентификация**: JWT (JSON Web Tokens)
- **Логирование**: Zerolog
- **Валидация**: Custom validator
- **Документация**: Swagger/OpenAPI 3.0
- **Контейнеризация**: Docker, Docker Compose
- **Миграции**: SQL-миграции
- **Конфигурация**: Environment variables

## 🏗️ Структура проекта

```bash
.
├── cmd                # Точка входа в приложение
├── docs               # Swagger документация
├── internal           # Внутренние пакеты
│   ├── app            # Инициализация приложения
│   ├── cache          # Кэширование
│   │   └── redis      # Реализация кэширования с помощью Redis
│   ├── config         # Конфигурация приложения
│   ├── database       # Работа с базой данных
│   │   ├── model      # Сущности БД
│   │   └── postgres   # Работа с PostgreSQL
│   ├── logger         # Логгер
│   │   └── zerolog    # Реализация логгера с zerolog
│   ├── server         # HTTP сервер и роутинг
│   │   ├── handler    # Обработчики эндпоинтов
│   │   └── middleware # Промежуточный слой
│   └── utils          # Вспомогательные утилиты
└── migrations         # Миграции БД
    └── postgres       # Миграции для PostgreSQL
```

## 🚀 Быстрый старт

### Предварительные требования
- Установленный Docker и Docker Compose
- Go 1.24+ (опционально, для локальной разработки)

### Запуск с Docker Compose
```bash
# Клонировать репозиторий
git clone https://github.com/rnymphaea/vk-internship.git
cd vk-internship

# Запустить сервисы
docker-compose up -d

# изменить переменные окружения (см. пример .env ниже)

# Приложение будет доступно по адресу:
# API: http://localhost:8080
# Swagger UI: http://localhost:8080/swagger/index.html
```
#### Пример .env
```
PORT=8080

JWT_SECRET=mysecret
JWT_TTL=24h
JWT_ISSUER=issuer

DB_TYPE=postgres
CACHE_TYPE=redis

LOGGER_TYPE=zerolog
LOGGER_LEVEL=debug
LOGGER_PRETTY=false

DB_URL=postgres://postgres:password@postgres:5432/marketplace?sslmode=disable

POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_HOST_PORT=5433
POSTGRES_DB_NAME=marketplace
POSTGRES_SSL_MODE=disable

REDIS_ADDR=redis:6379
REDIS_PORT=6379
REDIS_USER=appuser
REDIS_PASSWORD=your_secure_password_here
REDIS_DB=0
REDIS_MAX_RETRIES=3
REDIS_DIAL_TIMEOUT=10s
REDIS_TIMEOUT=5s
REDIS_TTL=24h
REDIS_MAX_FEED_ITEMS=10
```
## 🔧 Использование API
### Получение JWT токена
1. Зарегистрируйте нового пользователя:
```bash
POST /register
{
  "username": "newuser",
  "password": "SecurePass123!"
}
```
2. Если пользователь существует, авторизуйтесь
```bash
POST /login
{
  "username": "newuser",
  "password": "SecurePass123!"
}
```
3. Используйте полученный токен в заголовке запросов:
```bash
Authorization: Bearer <ваш_jwt_токен>
```

### Работа с объявлениями
- Создать объявление (доступно только с JWT токеном):
```bash
POST /ads
{
  "caption": "Продам ноутбук",
  "description": "Игровой ноутбук, 2023 года, идеальное состояние",
  "image_url": "https://example.com/laptop.jpg",
  "price": 75000.50
}
```

- Получить ленту объявлений:
```bash
GET /ads?page=1&page_size=10&sort_by=price&order=ASC&min_price=1000&max_price=100000
```

- Обновить объявление (доступно только с JWT токеном):
```bash
PUT /ads/{id}
{
  "caption": "Продам ноутбук (снижена цена)",
  "price": 70000.00
}
```

- Удалить объявление (доступно только с JWT токеном):
```bash
DELETE /ads/{id}
```

