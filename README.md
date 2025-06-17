# HSEPaymentSystem

## Введение
HSEPaymentSystem — это учебный проект, реализующий микросервисную архитектуру для управления платежами и заказами пользователей. Система состоит из двух сервисов — Payments Service и Orders Service, которые взаимодействуют между собой через брокер сообщений Kafka и используют паттерн Transactional Outbox для надёжной доставки событий. Проект демонстрирует современные подходы к построению отказоустойчивых и масштабируемых backend-систем на языке Go с использованием PostgreSQL и REST API.

---

## Важно (для ассистента)
Если честно, последние четыре часа я пытался закачать весь проект в докер, но в 3 ночи у меня уже нервы не выдерживают, поэтому я просто откатил все до момента, когда у меня все работало на локалке. `:(`

## Полезные ссылки
- [Cсылка на Postman-коллекцию](https://mikaeloganesian.postman.co/workspace/Mikael-Oganesian's-Workspace~b70df3ff-c2f3-465d-8fb5-aab76b425428/collection/44067443-34506887-5d08-4859-ac92-4f815b08dc4a?action=share&creator=44067443)
- Автор: [Mikael Oganesan](mailto:maoganesian@edu.hse.ru)
- Вопросы и предложения: issues/pull requests приветствуются!


# Payments Service

Микросервис на Go для управления пользовательскими счетами — создание, пополнение и просмотр баланса. Использует PostgreSQL и паттерн **Transactional Outbox** для надёжной доставки событий.

---

## Возможности

- Создание счёта по `user_id` (один счёт на пользователя)
- Пополнение счёта
- Получение баланса
- Transactional Outbox: события пишутся в outbox таблицу, затем надёжно обрабатываются
- Transactional Inbox: события получаются из inbox таблицы, затем надежно обрабатываются
- REST API с использованием **Gin**

---

## Стек технологий

- Go 1.21+
- [Gin](https://github.com/gin-gonic/gin) — HTTP-фреймворк
- PostgreSQL 14+
- `lib/pq` — PostgreSQL драйвер
- Ручная реализация Outbox-паттерна
- Kafka broker

---

## Структура микросервиса

```
- payments-service
    ├── cmd       # точка входа
    ├── db        # конфигурация и инициализация БД
    ├── handlers  # обработчики эндпоинтов
    ├── models    # типизированные модели
    └── services  # publisher и outbox
```

---

## Документация API

### POST `/accounts`
Создать новый аккаунт (один на пользователя).
http://localhost:8080/accounts
**Request:**
```json
{
  "user_id": "530fa038-f488-4a41-9ab7-7269f80af79d"
}
```

**Response:**
```json
{
    "message": "Account created"
}
```

---

### POST `/accounts/{id}/deposit`
Пополнить счет пользователя.
http://localhost:8080/accounts/530fa038-f488-4a41-9ab7-7269f80af79d/deposit
**Request:**
```json
{
  "user_id": "530fa038-f488-4a41-9ab7-7268f80af79d",
  "amount": 10000
}
```

**Response:**
```json
{
    "message": "Balance refilled"
}
```

---

### GET `/accounts/balance?user_id=`
Получить баланс счета по `user_id`.
http://localhost:8080/accounts/balance?user_id=530fa038-f488-4a41-9ab7-7268f80af79d
**Response:**
```json
{
    "balance": 103900
}
```

---


# Orders Service

Микросервис на Go для управления заказами пользователей — создание заказа, просмотр статуса и истории заказов. Использует PostgreSQL и паттерн **Transactional Outbox** для надёжной доставки событий.

---

## Возможности

- Создание заказа по `user_id`
- Получение информации о заказе
- Получение списка заказов пользователя
- Transactional Outbox: события пишутся в таблицу, затем надёжно обрабатываются
- REST API с использованием **Gin**

---

## Стек технологий

- Go 1.21+
- [Gin](https://github.com/gin-gonic/gin) — HTTP-фреймворк
- PostgreSQL 14+
- `lib/pq` — PostgreSQL драйвер
- Ручная реализация Outbox-паттерна
- Kafka broker

---

## Структура микросервиса

```
- orders-service
    ├── cmd       # точка входа
    ├── docs      # swagger-документация      
    ├── db        # конфигурация и инициализация БД
    ├── handlers  # обработчики эндпоинтов
    ├── models    # типизированные модели
    └── services  # outbox паттерн и kafka-приемник
```

---

## Документация API

### POST `/orders`
Создать новый заказ.
http://localhost:8081/orders
**Request:**
```json
{
    "user_id": "530fa038-f488-4a41-9ab7-7268f80af79d",
    "amount": 100,
    "description": "Покупка Iphone 17"
}
```

**Response:**
```json
{
    "message": "Order accepted"
}
```

---

### GET `/order/{order_id}`
Получить статус заказа по `id`.
http://localhost:8081/order/26d3ed81-e479-467c-ad73-5b4565e47c23
**Response:**
```json
{
    "id": "26d3ed81-e479-467c-ad73-5b4565e47c23",
    "status": "success"
}
```

---

### GET `/orders`
Получить список всех заказов.
http://localhost:8081/orders
**Response:**
```json
[
    {
        "ID": "26d3ed81-e479-467c-ad73-5b4565e47c23",
        "UserID": "530fa038-f488-4a41-9ab7-7268f80af79d",
        "Amount": 100,
        "Description": "Покупка Iphone 17",
        "Status": "success",
        "CreatedAt": "2025-06-17T20:55:35.806222Z",
        "UpdatedAt": "2025-06-17T20:55:53.109224Z"
    },
    {
        "ID": "24f00ead-931a-44f0-828d-ad8a360d806b",
        "UserID": "530fa038-f488-4a41-9ab7-7268f80af79d",
        "Amount": 100,
        "Description": "Покупка Iphone 16",
        "Status": "success",
        "CreatedAt": "2025-06-17T20:55:12.060518Z",
        "UpdatedAt": "2025-06-17T20:55:26.027711Z"
    },
    {
        "ID": "e8e48760-1b00-4d38-bc89-23c532065e1b",
        "UserID": "530fa038-f488-4a41-9ab7-7268f80af79d",
        "Amount": 100,
        "Description": "Покупка Iphone 16",
        "Status": "success",
        "CreatedAt": "2025-06-17T20:55:04.625613Z",
        "UpdatedAt": "2025-06-17T20:55:16.994625Z"
    },
    {
        "ID": "91b4ed8d-5711-4977-948e-3bc64858b954",
        "UserID": "530fa038-f488-4a41-9ab7-7269f80af79d",
        "Amount": 1000000,
        "Description": "Золото",
        "Status": "created",
        "CreatedAt": "2025-06-17T20:44:12.742694Z",
        "UpdatedAt": "2025-06-17T20:44:12.742694Z"
    },
    {
        "ID": "10411a26-ffb8-487c-ad73-d1b8ec065f74",
        "UserID": "530fa038-f488-4a41-9ab7-7268f80af79d",
        "Amount": 1000000,
        "Description": "Золото",
        "Status": "failed",
        "CreatedAt": "2025-06-17T18:41:57.370278Z",
        "UpdatedAt": "2025-06-17T18:42:15.255297Z"
    },
    {
        "ID": "d30b6cd9-ccdc-4a54-a479-44b24a101459",
        "UserID": "530fa038-f488-4a41-9ab7-7268f80af79d",
        "Amount": 100,
        "Description": "Сахар",
        "Status": "success",
        "CreatedAt": "2025-06-17T18:41:06.284596Z",
        "UpdatedAt": "2025-06-17T18:41:20.935591Z"
    }
]
```
