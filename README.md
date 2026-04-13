# Chirpy
Chirpy is an http server made with REST API and was made following boot.dev's guided project.

## Project Summary
This guided project was made to learn about how http server's handles and directs requests and how the architecture of web servers work.
It was made in go, and interacting with a postgresql database. We were using goose schema to migrate our versioning up or down.

## Setup

### 1. Clone repository
```bash
git clone https://github.com/VokalTuna/chirpy.git
cd chirpy
```

### 2. Install prerequsites
Make sure you have the following installed:
- Go 1.26
- PostgreSQL 15.17
- sqlc 1.30.0
- Goose 3.26

### 3. Configure environment variables
Make a `.env` file in the root of the project. It should look like like this:
```
DB_URL="postgres://username:@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
JWT_SECRET=<your-generated-secret>
POLKA_KEY=<your-polka-key>
```
- `DB_URL`: PostgreSQL connection string
- `PLATFORM`: set to `dev` for local development
- `JWT_SECRET`: secret used to sign JWTs
- `POLKA_KEY`: API key for polka integration.

Generate a JWT secret with
```bash
openssl rand -hex 32
```

### 4. Run migrations
```bash
goose -dir sql/schema postgres "postgres://vtuna:@localhost:5432/chirpy" up
```

### 5. Generate SQL code
```
sqlc generate
```
### 6. Start the server
```bash
go build -o out && ./out
```
## Notes
- The database specified in the `DB_URL` should exists before running migrations.
- If any changes has been made in `sql/queries` so should you rerun `sqlc generate`
- If any changes has been made in `sql/schema` so rerun migrations as needed.
