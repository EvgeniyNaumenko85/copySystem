# copySys

CopySys - файловый менеджер(аналог "Google Диск")

## Стек технологий
- Golang: Gin, Viper, GoDotEnv, lumberjack.v2, "regexp"
- JWT
- Postgres

## Логирование
* В системы реализованы 4 уровня логов:

[//]: # (1. Debug)
2. Errors
3. Info

[//]: # (4. Warning)

## Возможности
- сохранение файлов в специально отведенную директорию
- загрузка файлов из специально отведенной директории
- управление доступом пользователей к файлам
- поиск по названию файлов
- ограничение по размеру загружаемых файлов
- ограничение по размеру отведенного хранилища

[//]: # (- фильтр по типу фйлов)
[//]: # (- предоставление статистических данных о хранимых файлах &#40;тип хранимых файлов в % от хранилища&#41;)


## База данных
База данных находится в СУБД PostgreSQL
Файлы пользователей сохраняются локально в дирикторию "storage" 
(с распределением по папкам - именам пользователей)


## Запуск приложения
1. Скачайте необходимые зависимости: `go mod tidy`
2. Запустите приложение: `go run main.go`
* также в репозитории есть postman_collection всех запросов приложения 
