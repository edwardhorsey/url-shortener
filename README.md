# Url Shortener 

A simple URL shortener app written in Go.

URLs are stored in a PostgreSQL database then the database IDs are Base32 enocoded and used as the "short urls".

Features
 - Submit a URL and received a shortened URL
 - Lookup the destination address of a shortened URL (http://localhost:8080/dwlakj/show)

Technologies used:
 - [Golang](https://go.dev/)
 - [HTMX](https://htmx.org/)

## Running the app

Install [Go](https://go.dev/dl/)  (I use version 1.21 but this will very likely run on earlier versions).

Start a PostgreSQL database (via Docker Compose) and copy the contents of .env.example to .env.

```
docker-compose up -d
cp .env.example .env
```

Start the app
```
go run main.go

# and you should see
successfully connected to database, starting HTTP server on ":8080"
```

Navigate to [localhost:8080](http:localhost:8080) 🤏