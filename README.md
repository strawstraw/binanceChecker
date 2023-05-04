# Тестовое
Создать Go приложение, состоящее из клиентской и серверной частей. 
Приложение должно быть реализовано в виде CLI приложения на основе https://github.com/spf13/cobra

## Серверная часть
Должнаа предоставлять данные, предоставляемые публичным API биржи Binanсе, доступным по адресу https://api.binance.com/api/v3/ticker/price

Нужно реализовать два метода:

GET
```
$ curl http://localhost:3001/api/v1/rates?pairs=BTC-USDT,ETH-USDT
{ "ETH-USDT": 1780.123, "BTC-USDT": 46956.45 }
```
POST
```
$ curl -X POST --data '{ "pairs": ["BTC-USDT", "ETH-USDT"] }' http://localhost:3001/api/v1/rates
{ "ETH-USDT": 1780.123, "BTC-USDT": 46956.45 }
```

Запуск серверной части
```
$ go run . server
Listening at port 3001...
```

## Клиентская часть
Выполняет GET запрос к серверу и выводит ответ в консоль
```
$ go run . rate ---pair=ETH-USDT
1780.123
```
