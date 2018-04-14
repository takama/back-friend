# Back a friend

[![Build Status](https://travis-ci.org/takama/back-friend.svg?branch=master)](https://travis-ci.org/takama/back-friend)
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/takama/back-friend/issues)
[![Go Report Card](https://goreportcard.com/badge/github.com/takama/back-friend)](https://goreportcard.com/report/github.com/takama/back-friend)
[![codecov](https://codecov.io/gh/takama/back-friend/branch/master/graph/badge.svg)](https://codecov.io/gh/takama/back-friend)

Back a friend service: used [backer](https://github.com/takama/backer) package/library

## Known restrictions

1. A several verifications of input attributes skipped due to complicating the task
2. The project is used simplest `JSON` response without additional parameters like request attributes, totals and durations
3. The players and the tournaments could not be deleted

## Endpoints description

### Players endpoints

Funds (add to balance) player with 300 points. If no player exist should create new player with given
amount of points

```json
PUT /api/v1alpha/players/:id/fund
{
    "points": 300
}
```

Takes 300 points from player account

```json
PUT /api/v1alpha/players/:id/take
{
    "points": 300
}
```

Player balance

```json
GET /api/v1alpha/players/:id
{
    "player": "p1",
    "balance": 456.00
}
```

Players balances

```json
GET /api/v1alpha/players
[
    {
        "player": "p1",
        "balance": 456.0
    },
    {
        "player": "p2",
        "balance": 300.00
    }
]
```

### Tournaments endpoints

Announce tournament specifying the entry deposit

```json
POST /api/v1alpha/tournaments/:id/announce
{
    "deposit": 1000
}
```

Join player into a tournament and is he backed by a set of backers

```json
PUT /api/v1alpha/tournaments/:id/join
{
    "player": "p1",
    "backers": ["p2", "p3"]
}
```

Result tournament winners and prizes

```json
POST /api/v1alpha/tournaments/:id
{
    "winners": [
        {
            "player": "p1",
            "prize": 2000.00
        }
    ]
}
```

Get results of tournament winners and prizes

```json
GET /api/v1alpha/tournaments/:id
{
    "tournament": 1,
    "winners": [
        {
            "player": "p1",
            "prize": 2000.00
        }
    ]
}
```

Get results of tournaments

```json
GET /api/v1alpha/tournaments
[
    {
        "tournament": 1,
        "winners": [
            {
                "player": "p1",
                "prize": 2000.00
            }
        ]
    },
    {
        "tournament": 2,
        "winners": [
            {
                "player": "p1",
                "prize": 1000.00
            },
            {
                "player": "p2",
                "prize": 1000.00
            }
        ]
    }
]
```

### Engine endpoints

Reset the engine (database)

```sh
PUT /api/v1alpha/engine/reset
```

### Service endpoints

Service info and current statuses

```json
GET /info
{
    "version": "v0.1.0",
    "uptime": "1d 3h"
}
```

Service health

```sh
GET /healthz
```

Service readiness

```sh
GET /readyz
```

## Versioning automation

Using a special script to increase the release version

```sh
./bumper.sh
Current version v0.1.0.
Please enter new version [v0.1.1]:
```

## Contributing to the project

See the [contribution guidelines](docs/CONTRIBUTING.md) for information on how to
participate in the project by submitting pull requests or issues.

## Versioned changes

All changes in the project located in [changelog](docs/CHANGELOG.md)

## License

[MIT Public License](https://github.com/takama/back-friend/blob/master/LICENSE)
