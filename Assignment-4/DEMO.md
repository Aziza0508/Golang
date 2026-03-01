# Demo Video Scenario — Assignment 4

Follow these steps exactly for the 1–2 minute demo video.

## 1. Show No Containers Running

```bash
docker ps -a
```

Verify: no containers for this project are currently running or existing.

## 2. Build and Start

```bash
docker compose up --build
```

Wait for the build to finish and all services to start. The successful result is shown in the terminal.

## 3. Show Logs — DB Healthcheck Before Server Start

In the logs, point out:
- The **migrate** container runs and exits after applying migrations.
- The **web-app** prints `Waiting for database to be ready...` and then `Starting the Server...` AFTER the DB healthcheck passes.

## 4. Demonstrate Database Is Running

```bash
docker exec -it movies-db psql -U postgres -d moviesdb -c "\dt"
```

Shows the `movies` table exists.

## 5. Postman / curl — CRUD Operations (through Nginx on port 80)

### GET all movies (empty at first)
```bash
curl http://localhost/movies
```

### POST — Create movies
```bash
curl -X POST http://localhost/movies \
  -H "Content-Type: application/json" \
  -d '{"title":"SAW","genre":"horror","budget":500000,"hero":"JONNY DEPP","heroine":"Scarlet"}'

curl -X POST http://localhost/movies \
  -H "Content-Type: application/json" \
  -d '{"title":"TEST","genre":"Romance","budget":1000000,"hero":"BALE","heroine":"ARMAS"}'
```

### GET all movies (now has data)
```bash
curl http://localhost/movies
```

### GET single movie
```bash
curl http://localhost/movies/1
```

### PUT — Update movie
```bash
curl -X PUT http://localhost/movies/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"SAW 2","budget":750000}'
```

### DELETE movie
```bash
curl -X DELETE http://localhost/movies/2
```

### GET all movies (verify delete)
```bash
curl http://localhost/movies
```

## 6. Persistence Proof

### Add some data if not already added, then:
```bash
docker compose down
docker compose up -d
```

### Verify data is still there:
```bash
curl http://localhost/movies
```

The previously inserted movies are still present — named volume preserves data.

## 7. Docker Images — Show Size

```bash
docker images
```

**Point out:** The final Go image (`assignment-4-web-app`) is very small (~5–15MB) because of the multi-stage build with `scratch`. Compare to the base `golang:1.21-alpine` image (~250MB+).

## Architecture Summary

| Service       | Image               | Purpose                          |
|---------------|---------------------|----------------------------------|
| db            | postgres:15-alpine  | PostgreSQL database              |
| migrate       | migrate/migrate     | Runs schema migrations, then exits |
| web-app       | scratch (custom)    | Go API server on :8080           |
| nginx         | nginx:alpine        | Reverse proxy on :80             |

## Bonus Features Implemented

1. **Custom Network** — `movie-net` bridge network
2. **Makefile + run.sh** — `make up`, `make down`, `make scale`, `./run.sh`
3. **Graceful Shutdown** — SIGTERM/SIGINT handler with "Shutting down gracefully..."
4. **Scratch Image** — Final image uses `scratch` (smallest possible)
5. **Postgres Migrations** — golang-migrate runs as init container
6. **Nginx Reverse Proxy** — All traffic goes through nginx:80
7. **Service Scaling** — `docker compose up --scale web-app=3 -d` with nginx load balancing
