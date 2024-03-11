# pg-clickhouse-nats-redis-service

Проект представляет собой HTTP сервис, который взаимодействует с базами данных PostgreSQL, Redis, Clickhouse, реализуя CRUD API.

## Конфигурация

Конфигурация сервиса и баз данных выполняется в файле **.env**:

- ***APP_PORT*** - порт сервера
- ***POSTGRES_DB*** - название базы данных
- ***POSTGRES_USER*** - имя пользователя postgres
- ***POSTGRES_PASSWORD*** - пароль пользователя postgres
- ***POSTGRES_PORT*** - порт базы данных postgres
- ***POSTGRES_HOST*** - имя хоста postgres
- ***POSTGRES_MIGRATIONS_PATH*** - директория с миграциями базы данных postgres
- ***CLICKHOUSE_HOST*** - имя хоста clickhouse
- ***CLICKHOUSE_PORT*** - порт базы данных clickhouse
- ***CLICKHOUSE_USER*** - имя пользователя clickhouse
- ***CLICKHOUSE_PASSWORD*** - пароль пользователя clickhouse
- ***CLICKHOUSE_MIGRATIONS_PATH*** - директория с миграциями базы данных clickhouse
- ***NATS_HOST*** - имя хоста nats
- ***NATS_PORT*** - порт nats
- ***REDIS_HOST*** - имя хоста redis
- ***REDIS_PORT*** - порт базы данных redis
- ***REDIS_USER*** - имя пользователя redis
- ***REDIS_PASSWORD*** - пароль пользователя redis

## Запуск

Сервис и базы данных запускаются командой:
```bash
make all
```

## Использование сервиса

### Добавление новой записи:
```bash
curl -X POST -H "Content-Type: application/json" -d '{"name": "First record"}' 'http://localhost:8080/good/create?projectId=1'
```

### Получение списка всех записей:
```bash
curl -X GET 'http://localhost:8080/goods/list?limit=10&offset=0'
```

### Обновление записи:
```bash
curl -X PATCH -H "Content-Type: application/json" -d '{"name": "First record new name", "description": "First record new description"}' 'http://localhost:8080/good/update?id=1&projectId=1'
```

### Изменение приоритета записи:
```bash
curl -X PATCH -H "Content-Type: application/json" -d '{"newPriority": 5}' 'http://localhost:8080/good/reprioritize?id=1&projectId=1'
```

### Удаление записи:
```bash
curl -X DELETE 'http://localhost:8080/good/remove?id=1&projectId=1'
```

## P.S.
1. Так как в тз указано, что данные нужно кэшировать в redis при вызове GetList, то ключ для кэширования составляется из параметров limit и offset, а значением является список объектов goods.
2. При редактрировании какого-либо объекта, происходит инвалидация всего кэша redis, так как по отдельности это сделать слишком трудозатратно.
