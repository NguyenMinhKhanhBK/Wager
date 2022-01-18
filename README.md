## Context
This repo is a simple wager management backend system, where users can place/ buy wager as well as get wager list. Due to personal stuffs, I have just covered some basic aspects

## Dependencies
- Docker
- Docker-compose

## How to run
- To start service:
```
docker-compose up
```
- To run unit tests:
```
docker-compose ...
```

## How to test
### Place wager
- Valid request
```
curl --location --request POST 'http://localhost:8080/wagers' \
--header 'Content-Type: application/json' \
--data-raw '{
"total_wager_value": 100,
"odds": 120,
"selling_percentage": 1,
"selling_price": 200
}'
```
Response
```
{
  "id": 1,
  "total_wager_value": 100,
  "odds": 120,
  "selling_percentage": 1,
  "selling_price": 200,
  "current_selling_price": 200,
  "percentage_sold": null,
  "amount_sold": null,
  "place_at": 1642484487
}
```

- `total_wager_value` is not larger than 0
```
curl --location --request POST 'http://localhost:8080/wagers' \
--header 'Content-Type: application/json' \
--data-raw '{
"total_wager_value": 0,
"odds": 120,
"selling_percentage": 1,
"selling_price": 200
}'
```

Response 
```
{
  "error": [
    "TotalWagerValue must be larger than 0"
  ]
}
```
- `odds` is not large than 0
```
curl --location --request POST 'http://localhost:8080/wagers' \
--header 'Content-Type: application/json' \
--data-raw '{
"total_wager_value": 100,
"odds": 0,
"selling_percentage": 1,
"selling_price": 200
}'
```
Response
```
{
  "error": [
    "Odds must be larger than 0"
  ]
}
```
- `selling_percentage` is not between 1 and 100
```
curl --location --request POST 'http://localhost:8080/wagers' \
--header 'Content-Type: application/json' \
--data-raw '{
"total_wager_value": 100,
"odds": 100,
"selling_percentage": 0,
"selling_price": 200
}'
```
Response 
```
{
  "error": [
    "SellingPercentage must be larger than or equal 1"
  ]
}
```
Another request
```
curl --location --request POST 'http://localhost:8080/wagers' \
--header 'Content-Type: application/json' \
--data-raw '{
"total_wager_value": 100,
"odds": 100,
"selling_percentage": 101,
"selling_price": 200
}'

```
Response
```
{
  "error": [
    "SellingPercentage must be less than or equal 100"
  ]
}
```
- `selling_price` is less than `total_wager_value` * `selling_percentage` / 100
```
curl --location --request POST 'http://localhost:8080/wagers' \
--header 'Content-Type: application/json' \
--data-raw '{
"total_wager_value": 100,
"odds": 100,
"selling_percentage": 100,
"selling_price": 10
}'
```
Response
```
{
  "error": [
    "SellingPrice must be larger than TotalWagerValue * SellingPercentage"
  ]
}
```
- Invalid `total_wager_value` and `odds`
```
curl --location --request POST 'http://localhost:8080/wagers' \
--header 'Content-Type: application/json' \
--data-raw '{
"total_wager_value": 0,
"odds": 0,
"selling_percentage": 100,
"selling_price": 200
}'
```
Response
```
{
  "error": [
    "TotalWagerValue must be larger than 0",
    "Odds must be larger than 0"
  ]
}
```
### Get wager list
- Default filter (page = 1, limit = 10)
```
curl http://127.0.0.1:8080/wagers
```
Response
```
[
  {
    "id": 1,
    "total_wager_value": 100,
    "odds": 120,
    "selling_percentage": 1,
    "selling_price": 200,
    "current_selling_price": 200,
    "percentage_sold": null,
    "amount_sold": null,
    "place_at": 1642484487
  },
  {
    "id": 2,
    "total_wager_value": 100,
    "odds": 100,
    "selling_percentage": 100,
    "selling_price": 200,
    "current_selling_price": 200,
    "percentage_sold": null,
    "amount_sold": null,
    "place_at": 1642485725
  },
  
  ...
  
    {
    "id": 10,
    "total_wager_value": 100,
    "odds": 100,
    "selling_percentage": 100,
    "selling_price": 200,
    "current_selling_price": 200,
    "percentage_sold": null,
    "amount_sold": null,
    "place_at": 1642485730
  }
]

```
- Explicit page and limit
```
curl http://127.0.0.1:8080/wagers\?page\=1\&limit\=2
```
Response
```
[
  {
    "id": 1,
    "total_wager_value": 100,
    "odds": 120,
    "selling_percentage": 1,
    "selling_price": 200,
    "current_selling_price": 200,
    "percentage_sold": null,
    "amount_sold": null,
    "place_at": 1642484487
  },
  {
    "id": 2,
    "total_wager_value": 100,
    "odds": 100,
    "selling_percentage": 100,
    "selling_price": 200,
    "current_selling_price": 200,
    "percentage_sold": null,
    "amount_sold": null,
    "place_at": 1642485725
  }
]
```
- Invalid filters
```
curl http://127.0.0.1:8080/wagers\?page\=0\&limit\=0
```
Response
```
{
  "error": [
    "Page must be larger than 0",
    "Limit must be larger than 0"
  ]
}

```
### Buy wager
- `buying_price` is 0
```
curl --location --request POST 'http://localhost:8080/buy/1' \
--header 'Content-Type: application/json' \
--data-raw '{
"buying_price":0
}'
```
Response
```
{
  "error": [
    "BuyingPrice must be larger than 0"
  ]
}
```
- `buying_price` is larger than `current_selling_price`
```
curl --location --request POST 'http://localhost:8080/buy/1' \
--header 'Content-Type: application/json' \
--data-raw '{
"buying_price":1000
}'
```
Response
```
{
  "error": "buying price must be equal or smaller than current selling price"
}
```
- Invalid wager ID
```
curl --location --request POST 'http://localhost:8080/buy/100' \
--header 'Content-Type: application/json' \
--data-raw '{
"buying_price":10
}'
```
Response
```
{
  "error": "id not found"
}
```
- Success
```
curl --location --request POST 'http://localhost:8080/buy/1' \
--header 'Content-Type: application/json' \
--data-raw '{
"buying_price":50
}'
```
Response
```
{
  "id": 1,
  "wager_id": 1,
  "buying_price": 50,
  "bought_at": 1642486839
}
```
Get wager info
```
curl http://127.0.0.1:8080/wagers\?page\=1\&limit\=1
```
Response
```
[
  {
    "id": 1,
    "total_wager_value": 100,
    "odds": 120,
    "selling_percentage": 1,
    "selling_price": 200,
    "current_selling_price": 150,
    "percentage_sold": 25,
    "amount_sold": 50,
    "place_at": 1642484487
  }
]
```
## TODO
- Database migration
- CI/CD
