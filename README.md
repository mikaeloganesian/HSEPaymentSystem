# Payments Service

Микросервис на Go для управления пользовательскими счетами — создание, пополнение и просмотр баланса. Использует PostgreSQL и паттерн **Transactional Outbox** для надёжной доставки событий.

---

## Возможности

- Создание счёта по `user_id` (один счёт на пользователя)
- Пополнение счёта
- Получение баланса
- Transactional Outbox: события пишутся в таблицу, затем надёжно обрабатываются
- REST API с использованием **Gin**

---

## Стек технологий

- Go 1.21+
- [Gin](https://github.com/gin-gonic/gin) — HTTP-фреймворк
- PostgreSQL 14+
- `lib/pq` — PostgreSQL драйвер
- Ручная реализация Outbox-паттерна

---

## Структура микросервиса

```
- payments-service
--> db - конфигурация и инициализация БД
--> handlers - обработчики эндпоинтов
--> models - типизированные модели
--> outbox - publisher для вставки из outbox
```

---

## Эндпоинты

### POST /accounts
Создает новый аккаунт (в единственном экземпляре)
```json
{
  "user_id": "uuid"
}
```
```json
{
    "id": "int",
    "user_id": "uuid",
    "balance": 0,
    "created_at": "datetime"
}
```

### POST /accounts/{id}/deposit
Делает депозит на соответсвующий счет
```json
{
    "amount": "int"
}
```
```json
{
    "status": "success"
}
```

### GET /accounts/balance?user_id=
Получить баланс счета
```json
{
    "balance": "int"
}
```

---

## Список тестов и запуск тестов

### Cписок тестов
1. TestCreateAccount - создание счета и получение ответа
2. TestDeposit - депозит суммы на счет
3. TestGetBalance - получение суммы на счете

### Запуск
Из корневой директории проекта:
```bash
    cd payments-service
    go test ./tests/...
```
