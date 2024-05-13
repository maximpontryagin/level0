### WB Tech: level # 0 (Golang)
**Тестовое задание для стажировки wildberries:**
Необходимо разработать демонстрационный сервис с простейшим интерфейсом,
отображающий данные о заказе. Модель данных в формате JSON прилагается к
заданию.
Что нужно сделать:
1. Развернуть локально PostgreSQL:  
    1.1. Создать свою БД  
    1.2. Настроить своего пользователя  
    1.3. Создать таблицы для хранения полученных данных  
2. Разработать сервис:  
    2.1. Реализовать подключение и подписку на канал в nats-streaming  
    2.2. Полученные данные записывать в БД  
    2.3. Реализовать кэширование полученных данных в сервисе (сохранять in
    memory)  
    2.4. В случае падения сервиса необходимо восстанавливать кэш из БД  
    2.5 Запустить http-сервер и выдавать данные по id из кэша  
3. Разработать простейший интерфейс отображения полученных данных по id
заказа

Советы:
1. Данные статичны, исходя из этого подумайте насчет модели хранения в кэше и
в PostgreSQL. Модель в файле model.json
2. Подумайте как избежать проблем, связанных с тем, что в канал могут закинуть
что-угодно
3. Чтобы проверить работает ли подписка онлайн, сделайте себе отдельный
скрипт, для публикации данных в канал
4. Подумайте как не терять данные в случае ошибок или проблем с сервисом
5. Nats-streaming разверните локально (не путать с Nats)

### Стек используемых технологий:
* Golang
* PostgreSQL
* Nats-streaming
* Docker
* Viper

### Как развернуть проект:
1. Клонировать репозиторий
```
git clone https://github.com/maximpontryagin/level0
```
2. Запустить Nats-streaming. Docker образ:
```
docker run --name=NatsServer -p 4222:4222 -p 8222:8222 -d nats-streaming -p 4222 -m 8222
```
3. Запустить PostgreSQL. Docker образ:
```
docker run --name=NatsServer -p 4222:4222 -p 8222:8222 -d nats-streaming -p 4222 -m 8222
```
4. Запустить сервис:
```
go run .\cmd\consumer\main.go
```
5. Запустить отправку тестовых данных (для проверки работы сервиса):
```
go run .\cmd\publisher\main.go
```
5. Открыть файл `index.html` в браузере

### Пример запроса
```
http://127.0.0.1:8000/order/{order_id}
```
Через web интерфейс index.hml просто ввести order_id в поле id заказа

### Пример ответа:
{
    "order_uid": "b563feb7b2b84b6test6",
    "track_number": "WBILMTESTTRACK",
    "entry": "WBIL",
    "delivery": {
      "name": "Test Testov",
      "phone": "+9720000000",
      "zip": "2639809",
      "city": "Kiryat Mozkin",
      "address": "Ploshad Mira 15",
      "region": "Kraiot",
      "email": "test@gmail.com"
    },
    "payment": {
      "transaction": "b563feb7b2b84b6test",
      "request_id": "",
      "currency": "USD",
      "provider": "wbpay",
      "amount": 1817,
      "payment_dt": 1637907727,
      "bank": "alpha",
      "delivery_cost": 1500,
      "goods_total": 317,
      "custom_fee": 0
    },
    "items": [
      {
        "chrt_id": 9934930,
        "track_number": "WBILMTESTTRACK",
        "price": 453,
        "rid": "ab4219087a764ae0btest",
        "name": "Mascaras",
        "sale": 30,
        "size": "0",
        "total_price": 317,
        "nm_id": 2389212,
        "brand": "Vivienne Sabo",
        "status": 202
      }
    ],
    "locale": "en",
    "internal_signature": "",
    "customer_id": "test",
    "delivery_service": "meest",
    "shardkey": "9",
    "sm_id": 99,
    "date_created": "2021-11-26T06:22:19Z",
    "oof_shard": "1"
}
