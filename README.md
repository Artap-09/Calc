# Запуск через docker-compose
1. Запускаем `docker-compose.yml`

для этого вводим в командную строку находясь в директории с `docker-compose.yml`:

```
docker-compose up -d
```

<details>
<summary>Эта команда создаст 4 контейнера</summary>

```go
Creating MONGO  ... done

Creating CALC   ... done 

Creating PSQL   ... done

Creating Reader ... done
```
</details>

2. Создаем БД и Таблицу в SQL

для этого надо ввести в `sh` контейнера:

```sql
create database mysum;
\c mysum;
create table calculator(
    A   int,
    B   int,
    Sum int
);
```

# Кратко, что делает приложение
## Caclculator

Это сервер который слушает порт `4969` и обрабатывает `post` запросы.

<details>
<summary>Пример:</summary>

```http
POST http://localhost:4969/
Content-Type: application/json

{
  "a": 12,
  "b": 54
}
```
</details>

`calculator` суммирует `a` и `b` в `sum`, записывает результат в соответствующие колонки `postgres` и в формате `JSON` в `mongo` 

## ReaderDB

Это сервер который слушает порт `4979` и обрабатывает `post` запросы.

<details>
<summary>Пример:</summary>

```http
POST http://localhost:4979/
Content-Type: application/json

{
  "where": "mongo",
  "first": 1,
  "last": 10
}
```
</details>

`readerdb` читает из таблицы указаной в `where` (либо `mongo`, либо `postgres`) и пишет результат в `logs` и `Response`

<details>
<summary>Пример:</summary>

```
1 + 3 = 4
42 + 21 = 63
12 + 12 = 24
```
</details>
