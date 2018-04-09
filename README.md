# Back a friend

Back a friend service: used [backer](https://github.com/takama/backer) package/library

## Known restrictions

1. A several verifications of input attributes skipped due to complicating the task
2. The project is used simplest `JSON` response without additional parameters like request attributes, totals and durations
3. The players and the tournaments could not be deleted

## Players endpoints: /api/v1alpha/players

Funds (add to balance) player with 300 points. If no player exist should create new player with given
amount of points

```sh
PUT /api/v1alpha/players/:id/fund
{"points": 300}
response http codes: 201, 404, 400, 500
```

Takes 300 points from player account

```sh
PUT /api/v1alpha/players/:id/take
{"points": 300}
response http codes: 201, 404, 400, 500
```

Player balance

```sh
GET /api/v1alpha/players/:id
{"player": "p1", "balance": 456.00}
response http codes: 200, 404, 500
```

Players balances

```sh
GET /api/v1alpha/players
[{"player": "p1", "balance": 456.00}, {"player": "p2", "balance": 300.00}]
response http codes: 200, 500
```

## Tournaments endpoints: /api/v1alpha/tournaments

Announce tournament specifying the entry deposit

```sh
POST /api/v1alpha/tournaments/:id/announce
{"deposit": 1000}
response http codes: 201, 404, 400, 500
```

Join player into a tournament and is he backed by a set of backers

```sh
PUT /api/v1alpha/tournaments/:id/join
{"player": "p1", "backers": ["p2", "p3"]}
response http codes: 201, 404, 400, 500
```

Result tournament winners and prizes

```sh
POST /api/v1alpha/tournaments/:id
{winners": [{"player": "p1", "prize": 2000.00}]}
response http codes: 201, 404, 400, 500
```

Get results of tournament winners and prizes

```sh
GET /api/v1alpha/tournaments/:id
{"tournament": 1, "winners": [{"player": "p1", "prize": 2000.00}]}
reponse http codes: 200, 404, 400, 500
```

Get results of tournaments

```sh
GET /api/v1alpha/tournaments
[{"tournament": 1, "winners": [{"player": "p1", "prize": 2000.00}]}, {"tournament": 2, "winners": [{"player": "p1", "prize": 1000.00}, {"player": "p2", "prize": 1000.00}]}]
response http codes: 200, 500
```

## Engine endpoints

Reset the engine (database)

```sh
PUT /api/v1alpha/engine/reset
response http codes: 201, 500
```

## Service endpoints

Service info and current statuses

```sh
GET /info
{"version": "v0.1.0", "uptime": "1d 3h"}
response http codes: 200, 500
```

Service health

```sh
GET /healthz
response http codes: 200, 500
```

Service readiness

```sh
GET /readyz
response http codes: 200, 500
```
