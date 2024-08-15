# Auth Service

### Команды

- Прежде всего надо:

1) Удостовериться, что установлена утилита `make`.  
2) Удостовериться, что установлен `docker`
2) Запустить `make bin-deps`

Получить информацию по командам `make help`

- Как запустить приложение?

1) `make compose-up` - запустить контейнеры
2) `make run-app` - запустить приложение (локально)

- Как остановить контейнеры?

1) `make compose-down`

- Как посмотреть БД?

1) `docker ps` - находим Container ID у образа image
2) `docker exec -it 2089ce2ba4be bash` - вставляем нужный Container ID
3) `psql -p 5432 user -d postgres` - открывает консоль Postgres
4) `exit` - чтобы выйти из контейнера

- Работа с миграция:

1) `make run-app` - автоматически запускает миграции
2) `make migrate-create name="migration_name"` - создать файлы для миграций
3) `make migrate-up` - запустить миграцию