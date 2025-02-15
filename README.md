# Ozon-Internship
Тестовое задание для Go стажера-разработчика 

Задача: **Реализовать сервис, предоставляющий API по созданию сокращённых ссылок.**

Ссылка должна быть:

— Уникальной; на один оригинальный URL должна ссылаться только одна сокращенная ссылка;

— Длиной 10 символов;

— Из символов латинского алфавита в нижнем и верхнем регистре, цифр и символа _ (подчеркивание).

Сервис должен быть написан на Go и принимать следующие запросы по http:
1. Метод Post, который будет сохранять оригинальный URL в базе и возвращать сокращённый.
2. Метод Get, который будет принимать сокращённый URL и возвращать оригинальный.

Условие со звёздочкой (будет большим плюсом):

3. Сделать работу сервиса через GRPC, то есть составить proto и реализовать сервис с двумя соответствующими эндпойнтами

---

Решение должно соответствовать условиям:

— Сервис распространён в виде Docker-образа; 

— В качестве хранилища ожидаем in-memory решение и PostgreSQL. Какое хранилище использовать, указывается параметром при запуске сервиса; 

— Реализованный функционал покрыт Unit-тестами.

Результат предоставить в виде публичного репозитория на github.com

---

Что будем оценивать: 
    
    1) Как генерируются ссылки, почему предложенный алгоритм будет работать; насколько он соответствует заданию и прост в понимании.
    
    2) Как раскиданы типы по файлам, файлики по пакетам, пакеты по приложению: структуру проекта.
    
    3) Как обрабатываются ошибки в разных сценариях использования.
    
    4) Насколько удобен и логичен сервис в использовании.
    
    5) Как сервис будет себя вести, если им будут пользоваться одновременно сотни людей (как, например, YouTube, ya.cc).
    
    6) Что будет, если сервис оставить работать на очень долгое время.
    
    7) Общую чистоту кода.


--- 

## Документация по использованию

В первую очередь, клонируйте репозиторий проекта:
```bash
   git clone https://github.com/Bl00mGuy/Ozon-Internship.git
   cd urlshortener
```

Убедитесь, что у вас установлен Go, и выполните команду для установки зависимостей:
```
go mod tidy
```


#### Запуск с использованием in-memory хранилища:

Для того чтобы использовать in-memory хранилище, выполните команду:
```go
go run main.go -repo=inmemory -grpc-port=50051 -http-port=8080
```

#### Запуск с использованием PostgreSQL:

Для запуска с PostgreSQL в качестве хранилища передайте строку подключения через параметр -pgconn. Пример команды:
```go
go run main.go -repo=postgres -pgconn="postgres://postgres:password@localhost:5432/urlshortener?sslmode=disable" -grpc-port=50051 -http-port=8080
```

___

Пример Post запроса на сервер:
```
curl -X POST http://localhost:8080/urls/shorten \
-H "Content-Type: application/json" \
-d '{"original_url": "https://www.example.com/some-long-url"}'
```

Пример ответа от сервера:
```
{"shortenedUrl":"Hqveo2iu7D"}⏎
```