# Golang Prayers

## _Islamic Prayer times calculations_

## Features

- Docker setup for containerization, enabling easy deployment, and hot reloading.
- Utilizing Redis as a caching layer to improve performance, optimize data retrieval, and reduce the number of requests made to third-party APIs.
- Using PostgreSQL as the primary database.
- Using the Gin web framework.
- Implementing a data seeding process to populate the database with initial data from a JSON file.
- Integrating Google APIs to fetch user location and retrieve the corresponding timezone based on latitude and longitude coordinates.
- Providing Islamic prayer times for the user based on their provided coordinates.

## Installation

This project requires [GoLang](https://go.dev/doc/install) to run.

clone the project and ``.

```sh
cd  golang-prayers
docker-compose up -d --build
docker logs -f backend #to monitor the hot reload result
```

You can test it by

```curl
curl --location 'http://localhost/api/v1/azan/' \
--header 'Content-Type: application/json' \
--data '{
    "lat": 30.012000,
    "lng": 31.194375
}'
```
